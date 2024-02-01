package ir.ac.iust.ce;

import org.apache.storm.spout.SpoutOutputCollector;
import org.apache.storm.task.TopologyContext;
import org.apache.storm.topology.OutputFieldsDeclarer;
import org.apache.storm.topology.base.BaseRichSpout;
import org.apache.storm.tuple.Fields;
import org.apache.storm.tuple.Values;
import org.opencv.core.Mat;
import org.opencv.videoio.VideoCapture;
import org.opencv.videoio.Videoio;

import java.util.Map;

public class VideoFileSpout extends BaseRichSpout {
    private SpoutOutputCollector collector;
    private VideoCapture videoCapture;
    private final String videoPath;
    private long frameId = 0;

    public VideoFileSpout(String videoPath) {
        this.videoPath = videoPath;
    }

    @Override
    public void open(Map<String, Object> conf, TopologyContext context, SpoutOutputCollector collector) {
        this.collector = collector;
        this.videoCapture = new VideoCapture();
        this.videoCapture.open(videoPath);

        double fps = videoCapture.get(Videoio.CAP_PROP_FPS);
        System.out.println("FPS: " + fps);
    }

    @Override
    public void nextTuple() {
        Mat frame = new Mat();
        if (videoCapture.isOpened() && videoCapture.read(frame)) {
            collector.emit(new Values(frameId++, frame));
            System.out.println("FRAMEID: " + frameId);
        }
    }

    @Override
    public void declareOutputFields(OutputFieldsDeclarer outputFieldsDeclarer) {
        outputFieldsDeclarer.declare(new Fields("frameId", "frame"));
    }

    @Override
    public void close() {
        if (videoCapture.isOpened()) {
            videoCapture.release();
        }
    }
}
