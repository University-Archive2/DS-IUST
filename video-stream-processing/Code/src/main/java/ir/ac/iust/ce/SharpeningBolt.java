package ir.ac.iust.ce;

import org.apache.storm.topology.BasicOutputCollector;
import org.apache.storm.topology.OutputFieldsDeclarer;
import org.apache.storm.topology.base.BaseBasicBolt;
import org.apache.storm.tuple.Fields;
import org.apache.storm.tuple.Tuple;
import org.apache.storm.tuple.Values;
import org.opencv.core.CvType;
import org.opencv.core.Mat;
import org.opencv.imgproc.Imgproc;


public class SharpeningBolt extends BaseBasicBolt {
    @Override
    public void execute(Tuple tuple, BasicOutputCollector collector) {
        Mat frame = (Mat) tuple.getValueByField("processedFrame");
        Object frameId = tuple.getValueByField("frameId");

        // Apply sharpening filter
        Mat sharpenedFrame = new Mat();
        Mat kernel = new Mat(3, 3, CvType.CV_32F) {
            {
                put(0, 0, -1);
                put(0, 1, -1);
                put(0, 2, -1);

                put(1, 0, -1);
                put(1, 1, 9);
                put(1, 2, -1);

                put(2, 0, -1);
                put(2, 1, -1);
                put(2, 2, -1);
            }
        };
        Imgproc.filter2D(frame, sharpenedFrame, frame.depth(), kernel);

        collector.emit(new Values(frameId, sharpenedFrame));
    }

    @Override
    public void declareOutputFields(OutputFieldsDeclarer declarer) {
        declarer.declare(new Fields("frameId", "sharpenedFrame"));
    }
}
