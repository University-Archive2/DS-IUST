package ir.ac.iust.ce;

import org.apache.storm.task.TopologyContext;
import org.apache.storm.topology.BasicOutputCollector;
import org.apache.storm.topology.OutputFieldsDeclarer;
import org.apache.storm.topology.base.BaseBasicBolt;
import org.apache.storm.tuple.Fields;
import org.apache.storm.tuple.Tuple;
import org.apache.storm.tuple.Values;
import org.opencv.core.Core;
import org.opencv.core.Mat;
import org.opencv.imgproc.Imgproc;

import java.io.BufferedWriter;
import java.io.FileWriter;
import java.io.IOException;
import java.util.Map;


public class FrameAnalysisBolt extends BaseBasicBolt {
    private final String textFilePath;
    private BufferedWriter writer;

    public FrameAnalysisBolt(String textFilePath) {
        this.textFilePath = textFilePath;
    }

    @Override
    public void prepare(Map<String, Object> topoConf, TopologyContext context) {
        super.prepare(topoConf, context);
        try {
            this.writer = new BufferedWriter(new FileWriter(textFilePath));
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    @Override
    public void execute(Tuple tuple, BasicOutputCollector collector) {
        Mat frame = (Mat) tuple.getValueByField("frame");
        Object frameId = tuple.getValueByField("frameId");

        double averageBrightness = calculateAverageBrightness(frame);

        try {
            writer.write("Average Brightness of Frame " + frameId.toString() + ": " + averageBrightness + "\n");
        } catch (IOException e) {
            e.printStackTrace();
        }

        collector.emit(new Values(frameId, frame));
    }

    private double calculateAverageBrightness(Mat frame) {
        Mat gray = new Mat();
        Imgproc.cvtColor(frame, gray, org.opencv.imgproc.Imgproc.COLOR_BGR2GRAY);
        return Core.mean(gray).val[0];
    }

    @Override
    public void declareOutputFields(OutputFieldsDeclarer declarer) {
        declarer.declare(new Fields("frameId", "frame"));
    }

    @Override
    public void cleanup() {
        try {
            writer.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
