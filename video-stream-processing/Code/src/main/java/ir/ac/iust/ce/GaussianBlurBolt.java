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

public class GaussianBlurBolt extends BaseBasicBolt {


    @Override
    public void execute(Tuple tuple, BasicOutputCollector collector) {
        Mat frame = (Mat) tuple.getValueByField("processedFrame");
        Object frameId = tuple.getValueByField("frameId");

        // Apply Gaussian blur
        Mat blurredFrame = new Mat();
        int kernelSize = 31;
        Size kSize = new Size(kernelSize, kernelSize);
        Imgproc.GaussianBlur(frame, blurredFrame, kSize, 0);

        collector.emit(new Values(frameId, blurredFrame));
    }

    @Override
    public void declareOutputFields(OutputFieldsDeclarer declarer) {
        declarer.declare(new Fields("frameId", "blurredFrame"));
    }
}
