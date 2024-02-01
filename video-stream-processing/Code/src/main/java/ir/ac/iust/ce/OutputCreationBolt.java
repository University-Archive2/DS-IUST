package ir.ac.iust.ce;

import org.apache.storm.task.TopologyContext;
import org.apache.storm.topology.BasicOutputCollector;
import org.apache.storm.topology.OutputFieldsDeclarer;
import org.apache.storm.topology.base.BaseBasicBolt;
import org.apache.storm.tuple.Tuple;
import org.opencv.core.Mat;
import org.opencv.core.Size;
import org.opencv.videoio.VideoWriter;

import java.io.BufferedWriter;
import java.io.FileWriter;
import java.io.IOException;
import java.util.Map;

import static ir.ac.iust.ce.Config.*;

public class OutputCreationBolt extends BaseBasicBolt {
    private int frameCount = 0;
    private VideoWriter videoWriter;

    private final String videoOutputPath;

    private final String textOutputPath;
    private BufferedWriter writer;

    public OutputCreationBolt(String videoOutputPath, String textOutputPath) {
        this.videoOutputPath = videoOutputPath;
        this.textOutputPath = textOutputPath;

    }

    @Override
    public void prepare(Map<String, Object> topoConf, TopologyContext context) {
        super.prepare(topoConf, context);
        int code = VideoWriter.fourcc('M', 'J', 'P', 'G');
        this.videoWriter = new VideoWriter(videoOutputPath, code, fps, new Size(width, height), false);
        try {
            this.writer = new BufferedWriter(new FileWriter(textOutputPath));
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    @Override
    public void execute(Tuple input, BasicOutputCollector collector) {
        Mat frame = (Mat) input.getValueByField("aggregatedFrame");
        frameCount++;

        if (videoWriter.isOpened()) {
            videoWriter.write(frame);
        }
    }

    @Override
    public void cleanup() {
        if (videoWriter.isOpened()) {
            videoWriter.release();
        }

        try {
            writer.write("Number of Frames: " + frameCount);
            writer.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    @Override
    public void declareOutputFields(OutputFieldsDeclarer declarer) {
        // No fields to declare as this bolt does not emit any tuples
    }
}
