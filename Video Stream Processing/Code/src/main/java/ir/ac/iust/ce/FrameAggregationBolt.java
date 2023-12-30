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

import java.util.HashMap;
import java.util.Map;

public class FrameAggregationBolt extends BaseBasicBolt {
    private Map<Long, Tuple> gaussianBlurData;
    private Map<Long, Tuple> sharpeningData;

    @Override
    public void prepare(Map<String, Object> topoConf, TopologyContext context) {
        super.prepare(topoConf, context);
        this.gaussianBlurData = new HashMap<>();
        this.sharpeningData = new HashMap<>();
    }

    @Override
    public void execute(Tuple tuple, BasicOutputCollector collector) {
        Long frameId = tuple.getLongByField("frameId");

        String sourceComponent = tuple.getSourceComponent();

        if (sourceComponent.equals("gaussianBlurBolt")) {
            gaussianBlurData.put(frameId, tuple);
        } else if (sourceComponent.equals("sharpeningBolt")) {
            sharpeningData.put(frameId, tuple);
        }

        // Check if we have data from both bolts for this frame
        if (gaussianBlurData.containsKey(frameId) && sharpeningData.containsKey(frameId)) {
            // Get data from both maps
            Tuple gaussianTuple = gaussianBlurData.get(frameId);
            Tuple sharpeningTuple = sharpeningData.get(frameId);

            Mat blurredFrame = (Mat) gaussianTuple.getValueByField("blurredFrame");
            Mat sharpenedFrame = (Mat) sharpeningTuple.getValueByField("sharpenedFrame");


            Mat aggregatedFrame = new Mat();
            Core.addWeighted(blurredFrame, 0.5, sharpenedFrame, 0.5, 0, aggregatedFrame);

            gaussianBlurData.remove(frameId);
            sharpeningData.remove(frameId);

            collector.emit(new Values(aggregatedFrame));
        }
    }

    @Override
    public void declareOutputFields(OutputFieldsDeclarer declarer) {
        declarer.declare(new Fields("aggregatedFrame"));
    }
}
