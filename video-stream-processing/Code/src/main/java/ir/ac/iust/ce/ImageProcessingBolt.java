package ir.ac.iust.ce;

import org.apache.storm.topology.BasicOutputCollector;
import org.apache.storm.topology.OutputFieldsDeclarer;
import org.apache.storm.topology.base.BaseBasicBolt;
import org.apache.storm.tuple.Fields;
import org.apache.storm.tuple.Tuple;
import org.apache.storm.tuple.Values;
import org.opencv.core.Mat;
import org.opencv.core.Size;
import org.opencv.imgproc.Imgproc;

public class ImageProcessingBolt extends BaseBasicBolt {

    private int resizeWidth;
    private int resizeHeight;

    public ImageProcessingBolt(int resizeWidth, int resizeHeight) {
        this.resizeWidth = resizeWidth;
        this.resizeHeight = resizeHeight;
    }

    @Override
    public void execute(Tuple tuple, BasicOutputCollector collector) {
        Mat frame = (Mat) tuple.getValueByField("frame");
        Object frameId = tuple.getValueByField("frameId");

        // Convert to grayscale
        Mat grayFrame = new Mat();
        Imgproc.cvtColor(frame, grayFrame, Imgproc.COLOR_BGR2GRAY);

        // Resize the frame
        Mat resizedFrame = new Mat();
        Size size = new Size(resizeWidth, resizeHeight);
        Imgproc.resize(grayFrame, resizedFrame, size);

        collector.emit(new Values(frameId, resizedFrame));
    }

    @Override
    public void declareOutputFields(OutputFieldsDeclarer declarer) {
        declarer.declare(new Fields("frameId", "processedFrame"));
    }
}
