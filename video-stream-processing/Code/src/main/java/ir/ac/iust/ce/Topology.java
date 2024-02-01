package ir.ac.iust.ce;

import org.apache.storm.Config;
import org.apache.storm.LocalCluster;
import org.apache.storm.topology.TopologyBuilder;
import org.apache.storm.tuple.Fields;
import org.opencv.core.Core;

import static ir.ac.iust.ce.Config.*;


public class Topology {

    static {
        System.loadLibrary(Core.NATIVE_LIBRARY_NAME);
    }

    public static void main(String[] args) throws Exception {
        TopologyBuilder builder = new TopologyBuilder();

        builder.setSpout("videoFileSpout", new VideoFileSpout(videoPath));
        builder.setBolt("frameAnalysisBolt", new FrameAnalysisBolt(averageBrightnessPath)).shuffleGrouping("videoFileSpout");
        builder.setBolt("imageProcessingBolt", new ImageProcessingBolt(width, height)).shuffleGrouping("frameAnalysisBolt");
        builder.setBolt("gaussianBlurBolt", new GaussianBlurBolt()).shuffleGrouping("imageProcessingBolt");
        builder.setBolt("sharpeningBolt", new SharpeningBolt()).shuffleGrouping("imageProcessingBolt");
        builder.setBolt("frameAggregationBolt", new FrameAggregationBolt()).fieldsGrouping("gaussianBlurBolt", new Fields("frameId")).fieldsGrouping("sharpeningBolt", new Fields("frameId"));
        builder.setBolt("outputCreationBolt", new OutputCreationBolt(outputVideoPath, outputFrameCountPath)).shuffleGrouping("frameAggregationBolt");

        Config conf = new Config();
        conf.setDebug(true);

        System.out.println("Submitting topology to local cluster.");
        LocalCluster cluster = new LocalCluster();
        cluster.submitTopology("videoProcessingTopology", conf, builder.createTopology());
        Thread.sleep(10000);
        cluster.shutdown();
    }
}
