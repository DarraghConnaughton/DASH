# Dynamic Adaptive Streaming over HTTP (DASH)
Dynamic Adaptive Streaming over HTTP (DASH) is a streaming protocol which has been widely adopted by modern streaming service providers due to effectiveness at improving the end users Quality of Experience (QoE),
which is achieves through the optimizing for reduction of video stalls during a viewing session. DASH achieves improved QoE by varying the quality of the video across time in response to changing network conditions.


### How DASH works

- A video is encoded into various distinct representations of varying bitrate levels. Each representation is then segemented across time, the details of which are captured in a manifest file.
- DASH client requests the manifest for the video they wish to stream. The manifest will be used as reference material as they sequentially request video segments.
- Sequential retrieval of video segments occurs in two phases:
  - Initial playback buffer hydration; the client will request n segments to fill the playback buffer of t seconds. 
  - Steady state mode; the client will request segments from the server once the buffer dips below t seconds. A full buffer is a healthy buffer, signifying the client’s ability to hydrate the buffer in excess of the video players data ingestion requirements.

#### Unhealthy playback buffer

A depleting buffer indicates a lag in segment retrieval and forewarns of a stall if not remediated.
Likely cause of delay/loss of retrieval is poor network conditions. Client adapts by reducing the quality of the proceeding segments until a healthy hydration level has been achieved.
Tradeoff here being between video quality for video continuity.

## Vulnerability

The sequential delivery of segments exposes DASH to side channel attacks. Specifically, the bytes per second of the segments traversing the network, even in encrypted communication,
leaks an identifiable pattern for a subset of videos.
- To be explored; the Advanced Video Coding (AVC), H.264 for example, seeks to reduce the size of the video based on visual elements within the video, such as temporal and spatial redundancy.
  - **Temporal Redundancy (Inter-frame Compression)**
    - Given a set of consecutive frames with little variance, rather than sending frames for each, we send an initial frame which represents the initial scene and the delta describing how the scene changes across time.
  - **Spatial Redundancy (Intra-frame Compression)**
    - The frame is divided into blocks and only the difference between these blocks is encoded. 

Visual content within the video, which makes it unique, will result in 


![alt text](./images/networkdiagram.svg)
