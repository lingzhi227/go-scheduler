                                      The Thirty-Fifth AAAI Conference on Artificial Intelligence (AAAI-21)




Positions, Channels, and Layers: Fully Generalized Non-Local Network for Singer
                                  Identification
                                            I-Yuan Kuo, Wen-Li Wei, Jen-Chun Lin
                                     Institute of Information Science, Academia Sinica, Taiwan
                               iyuan.i.kuo@gmail.com, lilijinjin@gmail.com, jenchunlin@gmail.com




                            Abstract                                         focus, but also improves the feature representation by cap-
                                                                             turing long-range spatial (or spatio-temporal) dependencies
  Recently, a non-local (NL) operation has been designed as                  (Wang et al. 2018; Vaswani et al. 2017; Bahdanau, Cho, and
  the central building block for deep-net models to capture
  long-range dependencies (Wang et al. 2018). Despite its ex-
                                                                             Bengio 2015; Ramachandran et al. 2019). Among a mass
  cellent performance, it does not consider the interaction be-              of attention mechanisms, a non-local (NL) operation that
  tween positions across channels and layers, which is crucial               belongs to the self-attention mechanism has recently been
  in fine-grained classification tasks. To address the limitation,           proposed, and has achieved great success in various vision
  we target at singer identification (SID) task and present a fully          and audio processing tasks (Wang et al. 2018; Hsieh et al.
  generalized non-local (FGNL) module to help identify fine-                 2019; Jung et al. 2020; Zhang et al. 2019; Li et al. 2019).
  grained vocals. Specifically, we first propose a FGNL opera-               As illustrated in Figure 1 (a), the NL operation computes
  tion, which extends the NL operation to explore the correla-               the response at a position in an image (or audio frame) by
  tions between positions across channels and layers. Secondly,              attending to all positions and taking their weighted aver-
  we further apply a depth-wise convolution with Gaussian ker-               age in an embedding space to achieve the goal of capturing
  nel in the FGNL operation to smooth feature maps for better
  generalization. More, we modify the squeeze-and-excitation
                                                                             long-range dependencies. Despite its excellent performance,
  (SE) scheme into the FGNL module to adaptively emphasize                   the original NL module only considers the global spatial
  correlated feature channels to help uncover relevant feature               (or spatio-temporal) correlation by merging channels, which
  responses and eventually the target singer. Evaluating results             would lose important cues across channels and layers.
  on the benchmark artist20 dataset shows that the FGNL mod-                    To mimic the functionalities of the human auditory sys-
  ule significantly improves the accuracy of the deep-net mod-               tem and improve the effectiveness in singer identification
  els in SID. Codes are available at https://github.com/ian-k-               (SID) task, this study proposes a fully generalized non-local
  1217/Fully-Generalized-Non-Local-Network.                                  (FGNL) module that extends the NL module by learning
                                                                             explicit correlations among all of the elements (positions)
                        Introduction                                         across channels and layers, as shown in Figure 1 (b). Specif-
                                                                             ically, FGNL module contributes in three key aspects. First,
The ability of humans to identify singers under limited guid-                we propose the FGNL operation, which scales up the rep-
ance is remarkable. Take, for example, humans can quickly                    resentation power of NL operation to attend the interac-
learn to identify singers by listening to only a few clips of                tion among feature maps across channels and layers and re-
music from those singers. Even without prior knowledge                       veal the mutual similarity of the corresponding parts. Sec-
about singers, the human auditory system has evolved to                      ond, we suppress the noise in each feature map by integrat-
be able to handle such a task by performing different func-                  ing the Gaussian smoothing filter into the FGNL operation.
tionalities that include exhibiting attention for specific fre-              Third, we modify the squeeze-and-excitation (SE) scheme
quency bands, capturing long-range dependencies of audio                     (Hu, Shen, and Sun 2018) into the end of the FGNL mod-
features as a whole, and extracting distinctive cues for com-                ule to adaptively recalibrate channel-wise feature responses
parison. All these can be done under the influence of back-                  by explicitly modeling inter-dependencies among channels.
ground instrumental music and subtle sound variations from                   Figure 2 illustrates the details of the FGNL module.
different singers (fine-grained vocals).
   The functionalities provided by the human auditory sys-                      To the best of our knowledge, our work is the first to in-
tem were a perfect match to a particular class of deep learn-                troduce the attention mechanism for solving the SID task.
ing algorithms called attention mechanism. Attention mech-                   Extensive experimental results show that: 1) Compared with
anism is an attempt to mimic human brain action, that is,                    the NL module, our FGNL module can capture richer fea-
to selectively concentrate on a few relevant things, while ig-               ture representations and distinctive cues for prediction, and
noring others in deep-net models. It not only tells where to                 achieve the state-of-the-art results on the SID task; 2) The
                                                                             proposed FGNL module is flexible in the sense that it can
Copyright © 2021, Association for the Advancement of Artificial              be integrated into different deep-net architectures and be
Intelligence (www.aaai.org). All rights reserved.                            trained in an end-to-end fashion.


                                                                      8217
Figure 1: Compared with the original non-local (NL) operation computes the response at each position by attending to all other
positions in a single channel, the proposed fully generalized non-local (FGNL) operation further considers correlations among
all of the positions across channels and layers.


                      Related Work                                       ture representations for SID. Zhang et al. (Zhang et al. 2020)
As our goal is to develop attention mechanisms for captur-               use the WaveNet to learn feature representations directly
ing richer feature representations and distinctive cues so that          from the raw audio waveform in the time domain to iden-
they could be used to facilitate the SID task. We discuss rel-           tify singers. Despite the recent success of CNN and RNN,
evant literature and recent progress on both topics.                     both convolutional and recurrent operations can only pro-
                                                                         cess a local neighborhood (Wang et al. 2018), making it
                                                                         difficult to learn non-local context relations between audio
Singer Identification                                                    feature representations, which is essential for distinguishing
SID is a classic task in the field of music information re-              fine-grained vocals.
trieval (MIR) (Nasrullah and Zhao 2019; Hsieh et al. 2020;                  Regarding the second challenge, the key is to separate the
Zhang et al. 2020; Van, Quang, and Thanh 2019). It aims to               vocal parts of the given audio clips to minimize the influ-
automatically identify the performing singers in given au-               ence of instruments on the learning of the SID model. For
dio clips to facilitate the management of music libraries.               example, Van et al. (Van, Quang, and Thanh 2019) combine
There are two main challenges in the SID task. First, due                a gated recurrent unit (GRU) on U-Net to separate the vocal
to subtle differences in vocal organs, singers may have sim-             parts from the song that with mixed background accompani-
ilar singing voices (fine-grained vocals), resulting in small            ment. Sharma et al. (Sharma, Das, and Li 2019) introduce an
inter-class variations (Hsieh et al. 2020; Sundberg 1989). As            end-to-end trainable Wave-U-Net to learn the separation of
the number of singers to be considered increases, this issue             singing voices, thereby eliminating the interference of back-
becomes crucial. Second, since the songs in each singer’s al-            ground accompaniment on singer identity cues. Hsieh et al.
bums usually contain instrumental accompaniment, it is dif-              (Hsieh et al. 2020) use an open source tool called Open-
ficult for the SID model to extract vocal-only features from             Unmix (Stöter et al. 2019), which combines a three-layer
such recordings, which will reduce the generalization abil-              bidirectional LSTM and multiplicative skip connection to
ity of the SID model (Hsieh et al. 2020; Van, Quang, and                 separate the vocal and instrumental tracks of music, and has
Thanh 2019; Sharma, Das, and Li 2019; Rafii et al. 2018;                 made great progress. As source separation technology has
Sturm 2014).                                                             become more mature and has been successfully used to im-
   With the success of deep learning, deep-net models such               prove the performance of SID task, in this study, we inte-
as convolutional neural network (CNN) and recurrent neural               grate the source separation model (Stöter et al. 2019) into
network (RNN) have been widely used to address both chal-                our system and focus on solving the first challenge by intro-
lenges. For the first challenge, the core behind these methods           ducing attention mechanism.
is to learn discriminative feature representations for singers
to be identified. For example, Nasrullah and Zhao (Nasrul-               Attention Mechanism
lah and Zhao 2019) introduce an end-to-end trainable convo-              Attention mechanism has enjoyed widespread adoption as
lutional recurrent neural network (CRNN) to learn the dis-               a computational module for modeling sequences because
criminative feature representations and their temporal de-               of its ability to capture long-range dependencies and se-
pendency to achieve the SID task. Hsieh et al. (Hsieh et al.             lectively concentrate on the relevant subset of the input
2020) further add a branch in CRNN to incorporate melody                 (Vaswani et al. 2017; Bahdanau, Cho, and Bengio 2015; De-
features for better performance. Van et al. (Van, Quang, and             vlin et al. 2019; Yu et al. 2018). For example, Bahdanau et
Thanh 2019) use the bidirectional long short-term memory                 al. (Bahdanau, Cho, and Bengio 2015) present for the first
(LSTM) network to learn the temporal dependency of fea-                  time an attention mechanism and combine it with the RNN


                                                                  8218
encoder-decoder in a neural machine translation model to                 local context relations (i.e., long-range dependencies) across
allow selective attention to relevant information from a vari-           the whole feature map by weighting sum of the features at
able length source sentence. Vaswani et al. (Vaswani et al.              all positions,
2017) further propose a Transformer architecture to draw
global dependencies between input and output. This archi-                                 Y = f (θ(X), φ(X))g(X),                   (1)
tecture entirely replaces recurrence with self-attention, and            where T denotes the number of input video frames (when
greatly improves the performance of machine translation.                 the input is a single image, T can be ignored), H and W
Such a self-attention mechanism has also been extended to                denote the height and width of the feature map, C is the
other language representation models such as BERT (Devlin                number of channels, m is a reduction ratio, which refers to
et al. 2019) and achieved the state-of-the-art results.                  the bottleneck design used for reducing the computational
   Creating attention mechanisms to compensate for the                   complexity (Wang et al. 2018), f (·, ·) represents the pair-
weakness of convolution has also become an emerging                      wise function, which calculates the affinity between all po-
theme in vision tasks (Hu, Shen, and Sun 2018; Wang                      sitions, and θ(·), φ(·), and g(·) are learnable transformations
et al. 2018; Ramachandran et al. 2019; Woo et al. 2018;                  recommended to be implemented by using 1 × 1 or 1 × 1 ×
Roy, Navab, and Wachinger 2019; Bello et al. 2019). For                  1 convolution (Wang et al. 2018). Thus, the transformations
example, Hu et al. (Hu, Shen, and Sun 2018) present a                    can be written as
channel-wise attention mechanism to explicitly model the                                                           C
inter-dependencies between the channels of its spatial fea-                               θ(X) = XWθ ∈ RN × m ,                     (2)
tures. It is intended to select the useful feature maps and                                                       C
                                                                                                               N× m
suppress the others by considering the global information                                 φ(X) = XWφ ∈ R               ,            (3)
of each channel. Woo et al. (Woo et al. 2018) and Roy et al.             and                                       C
(Roy, Navab, and Wachinger 2019) further explore both spa-                                g(X) = XWg ∈ RN × m ,                     (4)
tial and channel-wise attentions, and verify that using both is
                                                                         parameterized by the weight matrices Wθ , Wφ , and Wg ∈
superior to using only the channel-wise attention. Recently,                   C
Wang et al. (Wang et al. 2018) show that self-attention is               RC× m , respectively. Here N denotes the collapsing of all
an instantiation of non-local mean (Buades, Coll, and Morel              the spatial or spatio-temporal positions in one dimension,
2005), and present a NL operation for the convolution-                   i.e., N = HW or N = HW T . In the implementation, the
based deep-net models to capture long-range dependencies.                original NL operation provides multiple options for f . For
Specifically, the NL operation computes the correlation ma-              simplicity, we choose the dot product as an example, i.e.,
trix between each spatial point in the feature maps to gener-
                                                                                       f (θ(X), φ(X)) = θ(X)φ(X)T ,                 (5)
ate an attention map, and then perform the attention-guided
dense context information aggregation. Such a NL opera-                  where the size of the resulting pairwise function f (·, ·) de-
tion has become the core component for various deep-net                                 C       C
                                                                         notes as RN × m × R m ×N → RN ×N . Thus, by substituting
architectures to capture non-local context relations, and has            equations (2) to (5) into (1), the response Y can be obtained
been successfully applied in various fields, including vi-               as
sion, audio, etc. (Wang et al. 2018; Li et al. 2019). Despite                             Y = XWθ WφT XT XWg ,                      (6)
its excellent performance, the original NL operation only
considers the global spatial (or spatio-temporal) correlation            where the pairwise matrix XWθ WφT XT ∈ RN ×N encodes
by merging channels, which may miss subtle but important                 the mutual similarity between any positions of the input fea-
cues across channels and layers in fine-grained classification           ture. The effect of NL operation can be understood as the
tasks. In this work, we propose the FGNL operation, which                self-attention mechanism (Vaswani et al. 2017) in the sense
extends the NL operation to further explore the explicit cor-            that each position (row) in the resulting Y is a linear combi-
relations among all of the elements (positions) across chan-             nation of all the positions (columns) of XWg weighted by
nels and layers to obtain richer feature representations and             the corresponding row of the pairwise matrix.
distinctive cues.
                                                                         Our FGNL Module
                        Approach                                         The original NL operation aims to capture the long-range
In this section, we elaborate the proposed FGNL module.                  dependencies between any two positions in one convolu-
We first revisit the original NL operation (Wang et al. 2018).           tional layer. However, it only calculates the dependencies of
Then, we will introduce three extensions of FGNL module                  any two positions in each channel separately, and aggregates
in detail, including FGNL operation, Gaussian smoothing                  all channel information in one convolutional layer together
filter, and modified squeeze-and-excitation (MoSE) scheme.               through a joint location-wise matrix f (θ(X), φ(X)). Thus,
                                                                         it will lose the interaction between positions across chan-
Revisiting NL Operation                                                  nels and layers. To this end, we generalize the original NL
                                                                         operation so that the long-range dependencies between any
The original NL operation (Wang et al. 2018) is revisited                positions of any channels and layers can be modeled.
in matrix form shortly. Given the input feature map X ∈                     We formulate the proposed FGNL module as follows.
RT ×H×W ×C for the NL operation, the goal is to obtain a                 Given a set of input feature maps F = {X1 , X2 , . . . , XL }
                         C
response Y ∈ RT ×H×W × m , which aims to capture the non-                for the FGNL module, the goal of the FGNL operation


                                                                  8219
Figure 2: A spatio-temporal FGNL module. The feature maps are shown as the shape of their tensors, e.g., T ×Hi ×Wi ×C for C
channels (proper reshaping is performed when noted). θ, φ, and g denote 1×1×1 convolutions, ⊗ denotes matrix multiplication,
  denotes the element-wise product, and ⊕ denotes element-wise sum. The computation of softmax is performed on each row.


is to obtain a set of non-local context responses R =                      bilities. [·, ·] is the operation of matrix concatenation, and
{Y1 , Y2 , . . . , YM }, where L represents the number of lay-             r(·) represents a rolling function, which rolls the elements
ers and M is the number of responses. For the sake of clar-                of the matrix along the channel axis. Thus, by concate-
ity, in the following, we use two-layer input feature maps                 nating the matrices from different layers and subsequently
Xi ∈ RT ×Hi ×Wi ×C and Xj ∈ RT ×Hj ×Wj ×C , namely                         rolling the matrix along the channel axis, the long-range de-
F = {Xi , Xj } as an example to explain the FGNL oper-                     pendencies between the positions across channels and lay-
ation. (See Figure 2.) To this end, each response Yk in R                  ers can be obtained through the operation of the pairwise
can be calculated by weighting sum of the features at all po-              function f (·, ·). Similar to (5), we choose the dot prod-
sitions,                                                                   uct as the operation for f , and then normalize it by us-
                                                                           ing the softmax computation. The size of the resulting pair-
     Yk = f (G(θ(Xi )), r([G(φi (Xi )), G(φj (Xj ))]))                                                              C      C
                                                                           wise function f (·, ·) denotes as RNi × m × R m ×(Ni +Nj ) →
          × [G(gi (Xi )), G(gj (Xj ))].                (7)
                                                                           RNi ×(Ni +Nj ) . Here Ni and Nj denote the collapsing of all
Similar to the NL operation, θ(·), φi (·), φj (·), gi (·), and             the spatial or spatio-temporal positions for layer i and j,
                                                                                                                                   C
gj (·) are learnable transformations. In the implementation,               respectively. To this end, the response Yk ∈ RNi × m can
we set the number of channels represented by the weight                    be calculated by the linear combination between the two
matrices Wθ , Wφi , Wφj , Wgi , and Wgj in the above                       matrices resulted from f (·, ·) and [G(gi (Xi )), G(gj (Xj ))].
                         1
transformations to m        of the number of channels in Xi and            As a result, R = {Y1 , Y2 , . . . , YM = C } can be ob-
                                                                                                                         m
Xj . Here m is a reduction ratio, and is set to 32 in our ex-              tained by repeating the operation of (7), which rolls matrix
periments, and G(·) represents the Gaussian smoothing fil-                 r([G(φi (Xi )), G(φj (Xj ))] m  C
                                                                                                              − 1 times along the channel
ter, which suppresses noise by performing the depth-wise                   axis.
convolution between the feature map and the Gaussian ker-                     Besides capturing the long-range dependencies with
nel. A two-dimensional Gaussian kernel, i.e., G(p, q) =                    FGNL operation, we further explore the relatedness between
  1     −(p2 +q 2 )/2σ 2
2πσ 2 e                   is adopted, where p and q represent the          each response Yk in R through the modified squeeze-and-
spatial coordinates of the feature map (i.e., resulting from               excitation (MoSE) scheme, which adaptively recalibrates
θ(Xi ) or φi (Xi ) or φj (Xj ) or gi (Xi ) or gj (Xj )), and σ             each response Yk by considering inter-dependencies over
is the standard deviation. As σ grows, the feature map be-                 the M responses. Specifically, the squeeze step spatially
comes smoother, providing more noise suppression capa-                     summarizes each response with global average pooling,


                                                                    8220
                                                                         Zhao 2019) is employed, which ensures that the songs from
                                                                         the same album are split either in the training, validation, or
                                                                         the test set, to eliminate additional clues provided by the al-
                                                                         bum. All evaluated deep-net models are trained with audio
                                                                         clips of length {3s, 5s, 10s}. Among them, 90% audio clips
                                                                         (so-called frames) are used for training and the rest are used
                                                                         for testing. The data in the validation set is split from 10%
                                                                         of the training data.

                                                                         Evaluation Protocols
                                                                         We integrate the proposed FGNL module into two state-of-
                                                                         the-art SID models, the convolutional recurrent neural net-
                                                                         work (CRNN) (Nasrullah and Zhao 2019) and the convo-
                                                                         lutional recurrent neural network with melody (CRNNM)
                                                                         (Hsieh et al. 2020) in order to compare performance. For
                                                                         both CRNN and CRNNM, we follow their original archi-
                                                                         tecture settings as benchmarks. Briefly, the CRNN archi-
Figure 3: The proposed CRNN FGNL and CRNNM FGNL                          tecture is defined as a stack of four convolutional layers,
architectures. © means concatenation.                                    two GRU layers, and one fully connected (FC) layer. The
                                                                         CRNNM architecture is basically the same as CRNN, ex-
                                                                         cept that CRNNM also includes a branch related to melody.
while the excitation function emphasizes those responses                 Such a melody branch consists of a stack of four convo-
that play crucial roles in identifying the target. In between            lutional layers, and its output will be concatenated to the
the squeeze layer and the excitation layer, we use two con-              main branch of the CRNN for subsequent processing. For
volutional layers that connect ReLU and softmax activations              the proposed FGNL networks, we insert the FGNL module
respectively to modify the original SE scheme (Hu, Shen,                 into the network architecture of CRNN and CRNNM, re-
and Sun 2018). We depict the MoSE operation as follows:                  spectively named CRNN FGNL and CRNNM FGNL. For
                                                                         CRNN FGNL, as shown in Figure 3 (a), we insert a FGNL
              MoSE(R) = w,        R̃ = w     R,            (8)           module after the fourth convolutional layer to model the
                                                                         non-local context relations between the feature maps of the
where R̃ is a re-weighted set of non-local context responses,            fourth and third convolutional layers. For CRNNM FGNL,
w ∈ RM is the excitation vector, and denotes the element-                as shown in Figure 3 (b), we insert the FGNL module after
wise product.                                                            the fourth convolutional layer of the mel-spectrogram and
   Finally, as in the design of the NL module (Wang et al.               melody branches to model the non-local context relations
2018), we use residual connection to generate the output                 between the feature maps of the fourth and third convolu-
feature representation (map) Z ∈ RT ×Hi ×Wi ×C (refer to                 tional layers. For the training of the above deep-net models,
Figure 2) of the FGNL module as follows:                                 we apply random initialization for the weights, a constant
                                                                         learning rate of 10−4 , the dropout and batch normalization
                     Z = R̃Wz + Xi ,                       (9)           to avoid over-fitting, and the Adam solver (Kingma and Ba
where Wz is a learnable weight matrix, which can be im-                  2015) for optimization. Each model is trained by using back-
plemented by using 1 × 1 or 1 × 1 × 1 convolution (i.e.,                 propagation algorithm (including back-propagation through
depends on frame-wise (spatial) classification or sequence-              time algorithm) with the objective of softmax cross entropy
wise (spatio-temporal) classification task), and the number              under the supervision of the ground truth artist (singer) la-
of channels in Wz is scaled up to match the number of chan-              bel. The meta-parameters of each model are set based on the
nels in Xi . “+Xi ” denotes a residual connection (He et al.             validation error.
2016). Such a residual connection allows us to insert a new                 To evaluate whether the background accompaniment will
FGNL module into any pre-trained model, without breaking                 affect the generalization ability of the above deep-net mod-
its initial behavior (e.g., if Wz is initialized as zero). As a          els, two evaluation settings are considered, including the
result, by further considering the re-weighted non-local con-            original audio file and the vocal-only. The difference be-
                                                                         tween them is that the vocal-only setting further employs
text responses R̃, the information in Z is richer so Z can be
                                                                         the Open-Unmix toolkit (Stöter et al. 2019) to separate the
regarded as enhanced Xi .
                                                                         vocal parts from each audio file in training and test. In the
                                                                         experiments, we report the evaluation results of each deep-
                      Experiments                                        net model at the frame level and the song level. Specifically,
To demonstrate the effectiveness of the proposed FGNL                    at the frame level, each t-length (3s, 5s, or 10s) audio spec-
module, we conduct SID experiments on the benchmark                      trogram is treated as an independent sample, and the perfor-
artist20 dataset (Ellis 2007), which includes a total of 1,413           mance is measured by taking the F1 score across all samples
complete songs collected from 20 artists (singers). In the               in the test set. For evaluation at the song level, majority vot-
experiments, album-split (Hsieh et al. 2020; Nasrullah and               ing will be applied to select the most frequent frame level


                                                                  8221
                                                 Original Audio File                         Vocal-Only
                                            Frame Level       Song Level             Frame Level      Song Level
  Model                     Type         3s    5s   10s    3s    5s   10s         3s    5s   10s   3s    5s   10s #Parameters
  CRNN                      Average     0.44 0.45 0.48 0.57 0.55 0.58            0.42 0.46 0.51 0.72 0.74 0.74
                                                                                                                    394,516
  (Nasrullah and Zhao 2019) Best        0.46 0.47 0.53 0.62 0.59 0.60            0.44 0.48 0.53 0.76 0.79 0.77
                            Average     0.52 0.54 0.55 0.72 0.73 0.73            0.44 0.47 0.51 0.79 0.80 0.79
  CRNN FGNL (Ours)                                                                                                  584,141
                            Best        0.54 0.57 0.58 0.76 0.79 0.78            0.44 0.48 0.53 0.81 0.82 0.83
  CRNNM                     Average     0.47 0.47 0.51 0.62 0.61 0.65            0.42 0.46 0.49 0.73 0.75 0.73
                                                                                                                    778,772
  (Hsieh et al. 2020)       Best        0.48 0.50 0.53 0.67 0.68 0.69            0.43 0.47 0.50 0.75 0.79 0.75
                            Average     0.54 0.55 0.58 0.74 0.74 0.73            0.42 0.47 0.52 0.77 0.83 0.81
  CRNNM FGNL (Ours)                                                                                                1,175,381
                            Best        0.55 0.57 0.63 0.82 0.81 0.83            0.44 0.47 0.53 0.83 0.84 0.86
Table 1: The average and best F1 scores of the frame level and the song level in various length settings. Each t-length (3s, 5s, or
10s) experiment repeats three independent runs. Bold is the comparison winner of the same series (CRNN or CRNNM) model.

                                                                Original Audio File                         Vocal-Only
                                                           Frame Level       Song Level             Frame Level      Song Level
   Model                                     Type       3s    5s   10s    3s    5s   10s         3s    5s   10s   3s    5s   10s
                                             Average   0.44 0.45 0.48 0.57 0.55 0.58            0.42 0.46 0.51 0.72 0.74 0.74
   CRNN (Nasrullah and Zhao 2019)
                                             Best      0.46 0.47 0.53 0.62 0.59 0.60            0.44 0.48 0.53 0.76 0.79 0.77
   CRNN NL                                   Average   0.51 0.52 0.54 0.71 0.69 0.69            0.42 0.46 0.50 0.77 0.78 0.76
   (w/o the cues across channels and layers) Best      0.53 0.53 0.55 0.76 0.74 0.74            0.43 0.46 0.51 0.81 0.81 0.79
   CRNN FGNL LIGHT                           Average   0.51 0.54 0.54 0.70 0.73 0.69            0.43 0.47 0.51 0.77 0.77 0.77
   (w/o the cues across layers)              Best      0.54 0.55 0.55 0.78 0.77 0.78            0.44 0.48 0.53 0.80 0.80 0.82
                                             Average   0.52 0.54 0.55 0.72 0.73 0.73            0.44 0.47 0.51 0.79 0.80 0.79
   CRNN FGNL (Ours)
                                             Best      0.54 0.57 0.58 0.76 0.79 0.78            0.44 0.48 0.53 0.81 0.82 0.83
Table 2: Ablation experiments of CRNN with three attention modules, including NL (Wang et al. 2018), FGNL LIGHT, and
FGNL. Each t-length experiment repeats three independent runs. Bold indicates the comparison winner of the model.


artist prediction as the final prediction for each song. Note            ment of the FGNL module is not just because it adds the
that in the implementation, if the confidence (softmax out-              number of parameters to the baseline model. To see this, we
put) of the test frame is less than 0.5, it will be removed and          note that in Table 1, CRNN FGNL has better performance
will not participate in voting (Nasrullah and Zhao 2019). The            than CRNNM but has fewer parameters.
F1 score is then reported by song to quantify performance.                  For comparing the original audio file setting with the
                                                                         vocal-only setting at the frame level and the song level, as
Results and Comparisons                                                  shown in Table 1, we first notice that the vocal-only setting
For all the above competition methods, Table 1 summarizes                at the frame level performs worse than the original audio file
the average and best test F1 scores of the frame level and               setting. Such results indicate that a model trained with the
the song level resulted from three independent runs. For                 original audio files may benefit from the additional infor-
the comparison between CRNN and CRNNM, similar to                        mation in the accompaniment. This is supported by another
the results in (Hsieh et al. 2020), the results first show that          observation. It is observed that the model can identify the
CRNNM is better than CRNN in most settings. Such re-                     singer, even if some segments (e.g., intro, inter, or outro)
sults indicate that further consideration of melody-related              in the song do not contain the vocals. However, it is inter-
features is positive for SID. However, although CRNNM                    esting that the vocal-only setting at the song level performs
outperformed CRNN, the performance is still limited. One                 better than the original audio file setting. This is because in
explanation for this may be that both convolutional and re-              song level prediction, lower confidence frames will be re-
current operations in CRNNM only consider a local neigh-                 moved and will not contribute to the voting. Although the
borhood (Wang et al. 2018), so it is difficult to capture the            accompaniment in the original audio file setting will pro-
non-local context relations (i.e., long-range dependencies)              vide extra information, the confidence is usually low. This is
between audio features to distinguish singer. To tackle the is-          because the accompaniment in some vocal segments of the
sue, we introduce the attention mechanism and develop the                song could confuse the identification. In this case, the source
FGNL module to explicitly model the correlations among                   separation technique, which is used to separate the human
all of the positions in the feature map across channels and              voice from the original audio, would increase the identifica-
layers. By further integrating the FGNL module, the results              tion confidence of the SID model. Thus, the results indicate
support that CRNN FGNL and CRNNM FGNL can learn                          that source separation plays a positive role when considering
richer feature representations and distinctive cues to com-              the identification confidence in the song level. All in all, the
plete SID. That is, compared with the original CRNN and                  proposed FGNL module makes significant improvements to
CRNNM, CRNN FGNL and CRNNM FGNL achieve great                            CRNN and CRNNM in both the original audio file and the
improvements. In addition, it is noteworthy that the improve-            vocal-only settings.


                                                                  8222
Figure 4: Visualization of the embeddings (projected into 2-D space by t-SNE) under the original audio file setting of the 5-sec
frame level test samples. From left to right are CRNN, CRNN NL, CRNN FGNL LIGHT, and CRNN FGNL.

                                                            Original Audio File                         Vocal-Only
                                                       Frame Level       Song Level             Frame Level      Song Level
     Model                               Type       3s    5s   10s    3s    5s   10s         3s    5s   10s   3s    5s   10s
     CRNN FGNL                           Average   0.52 0.54 0.54 0.71 0.73 0.72            0.44 0.46 0.50 0.78 0.79 0.76
     without Gaussian smoothing          Best      0.53 0.55 0.55 0.77 0.80 0.74            0.45 0.47 0.52 0.81 0.83 0.81
     CRNN FGNL                           Average   0.52 0.53 0.54 0.71 0.72 0.71            0.43 0.45 0.51 0.78 0.78 0.78
     without MoSE                        Best      0.54 0.54 0.55 0.77 0.77 0.77            0.44 0.46 0.53 0.81 0.82 0.83
     CRNN FGNL                           Average   0.53 0.54 0.54 0.72 0.71 0.70            0.44 0.46 0.51 0.80 0.76 0.77
     with Gaussian smoothing and SE      Best      0.54 0.55 0.57 0.78 0.79 0.78            0.44 0.48 0.52 0.83 0.78 0.83
     CRNN FGNL                           Average   0.52 0.54 0.55 0.72 0.73 0.73            0.44 0.47 0.51 0.79 0.80 0.79
     with Gaussian smoothing and MoSE    Best      0.54 0.57 0.58 0.76 0.79 0.78            0.44 0.48 0.53 0.81 0.82 0.83
Table 3: Ablation experiments of CRNN FGNL with and without Gaussian smoothing, MoSE, and SE (Hu, Shen, and Sun
2018) mechanisms. Bold indicates the comparison winner of the model.


   To verify whether the cues across channels and layers in            embedding vectors to a 2-D space for visualization. Briefly,
the proposed FGNL module are useful, we conducted ab-                  for each of the above models, we regard the output of the last
lation experiments under the CRNN architecture. Specif-                layer of GRU in the CRNN architecture as embedding and
ically, the proposed CRNN FGNL is compared with the                    visualize it through t-SNE. For space limit, we visualize the
CRNN FGNL LIGHT (i.e., without the cues across layers),                four competing models under the setting of the original au-
the CRNN NL (Wang et al. 2018) (i.e., without the cues                 dio file at the 5-sec frame level. The audio samples of testing
across channels and layers), and the original CRNN. All at-            set are drawn and colored according to the ground truth artist
tention modules (i.e., NL, FGNL LIGHT, and FGNL) are                   (singer) labels in Figure 4. It can be seen from the result of
inserted after the fourth convolutional layer of the CRNN              CRNN FGNL LIGHT that samples from different singers
architecture (Nasrullah and Zhao 2019). For performance                are fairly well-separated in the embedding space. The re-
comparison, it is obvious from Table 2 that CRNN NL is                 sult of CRNN NL looks chaotic and less separated, suggest-
superior to CRNN in almost all settings. The results con-              ing again that a model taking the cues across channels (i.e.,
firm that by further introducing NL module to model the                CRNN FGNL LIGHT) may achieve SID better. Finally, as
non-local context relations of audio features, SID perfor-             shown in Table 2 and Figure 4, by simultaneously explor-
mance can indeed be improved. Despite its excellent perfor-            ing the correlations between positions across channels and
mance, it can only compute the response at each position by            layers, the CRNN FGNL can indeed capture richer feature
attending to all other positions in each channel separately,           representations and distinctive cues to facilitate the identifi-
which will lose important information between positions                cation of singers. Overall, CRNN FGNL achieves the best
across channels and layers. Compared with CRNN NL, the                 performance among the above competing models.
CRNN FGNL LIGHT demonstrates that by further consid-                      Finally, to evaluate whether integrating the Gaussian
ering the cues across channels, the performance can indeed             smoothing filter and the modified squeeze-and-excitation
be improved. The results can be further verified by visualiz-          (MoSE) scheme into the FGNL module (refer to Figure 2)
ing the feature embedding in each competing model. To this             can improve the generalization ability of the model, abla-
end, we employ t-distributed stochastic neighbor embedding             tion experiments are further conducted in CRNN FGNL (re-
(t-SNE) (Maaten and Hinton 2008) to project the computed               fer to Figure 3(a)). Four settings are considered, including


                                                                8223
CRNN FGNL without Gaussian smoothing, CRNN FGNL                          IEEE/CVF International Conference on Computer Vision
without MoSE, CRNN FGNL with Gaussian smoothing and                      (ICCV), 3286–3295.
squeeze-and-excitation (SE) (Hu, Shen, and Sun 2018), and                Buades, A.; Coll, B.; and Morel, J.-M. 2005. A non-local al-
CRNN FGNL with Gaussian smoothing and MoSE (our full                     gorithm for image denoising. In Proceedings of the IEEE
version). Table 3 summarizes performance comparison. It                  Conference on Computer Vision and Pattern Recognition
is observed that the CRNN FGNL with Gaussian smooth-                     (CVPR), 60–65.
ing and MoSE outperforms CRNN FGNL without Gaus-
sian smoothing and CRNN FGNL without MoSE, demon-                        Devlin, J.; Chang, M.-W.; Lee, K.; and Toutanova, K. 2019.
strating that generalizing the ability of the model could re-            BERT: Pre-training of Deep Bidirectional Transformers for
sult from both using the Gaussian smoothing filter to sup-               Language Understanding. In Proceedings of NAACL-HLT,
press the noise of the feature map, and using the MoSE                   4171–4186.
scheme to recalibrate the channel-wise feature responses.                Ellis, D. P. W. 2007. Classifying Music Audio with
Besides, comparing the original SE and the proposed MoSE                 Timbral and Chroma Features. In Proceedings Inter-
scheme, the results show that CRNN FGNL with Gaus-                       national Society for Music Information Retrieval Con-
sian smoothing and MoSE is better than CRNN FGNL with                    ference (ISMIR), 339–340.           [Online] https://labrosa.
Gaussian smoothing and SE. Such results indicate that us-                ee.columbia.edu/projects/artistid/.
ing the convolutional layer and softmax operation instead of
fully connected layer and sigmoid operation in SE module                 He, K.; Zhang, X.; Ren, S.; and Sun, J. 2016. Deep Resid-
can indeed increase the generalization ability of the model.             ual Learning for Image Recognition. In Proceedings of the
All in all, the ablation experiments show that the Gaus-                 IEEE Conference on Computer Vision and Pattern Recogni-
sian smoothing filter and the MoSE scheme improve the                    tion (CVPR), 770–778.
SID performance for the deep-net model. More experiments                 Hsieh, T.-H.; Cheng, K.-H.; Fan, Z.-C.; Yang, Y.-C.; and
and a demo video can be found at https://github.com/ian-k-               Yang, Y.-H. 2020. Addressing The Confounds Of Accom-
1217/Fully-Generalized-Non-Local-Network.                                paniments In Singer Identification. In IEEE International
                                                                         Conference on Acoustics, Speech and Signal Processing
                       Conclusions                                       (ICASSP), 1–5.
We have introduced a new attention mechanism called the                  Hsieh, T.-I.; Lo, Y.-C.; Chen, H.-T.; and Liu, T.-L. 2019.
fully generalized non-local (FGNL) module, which can bet-                One-Shot Object Detection with Co-Attention and Co-
ter capture the non-local context relations (i.e., long-range            Excitation. In Advances in Neural Information Processing
dependencies) of audio features to help identify fine-grained            Systems (NeurIPS), 2725–2734.
vocals. The results have demonstrated that the FGNL mod-
                                                                         Hu, J.; Shen, L.; and Sun, G. 2018. Squeeze-and-Excitation
ule significantly improves the accuracy of the deep-net mod-
                                                                         Networks. In Proceedings of the IEEE Conference on Com-
els in singer identification (SID) task and achieves the state-
                                                                         puter Vision and Pattern Recognition (CVPR), 7132–7141.
of-the-art level. Moreover, it is shown that the proposed
FGNL module is superior to the popular non-local (NL)                    Jung, Y.; Kim, D.; Woo, S.; Kim, K.; Kim, S.; and Kweon,
module (Wang et al. 2018) by explicitly modeling the rich                I. S. 2020. Hide-and-Tell: Learning to Bridge Photo Streams
inter-dependencies between any positions across channels                 for Visual Storytelling. In The Thirty-Fourth AAAI Confer-
and layers in the feature space, while the NL module only                ence on Artificial Intelligence (AAAI), 11213–11220.
considers the correlations between positions along the spe-              Kingma, D. P.; and Ba, J. L. 2015. Adam: A Method for
cific channel. Based on the promising outcomes, our future               Stochastic Optimization. In International Conference on
work will focus on developing more effective loss functions              Learning Representations (ICLR).
to improve the fineness of the learned feature representation.
We also plan to expand the scale of the experiments to other             Li, X.; Li, Y.; Li, M.; Xu, S.; Dong, Y.; Sun, X.; and Xiong,
tasks in the future, such as vision tasks.                               S. 2019. A Convolutional Neural Network with Non-Local
                                                                         Module for Speech Enhancement. In Proc. Interspeech,
                   Acknowledgments                                       1796–1800.
This work was supported in part by MOST grant 107-2218-                  Maaten, L. v. d.; and Hinton, G. 2008. Visualizing Data us-
E-001-010-MY2. The third author would like to thank his                  ing t-SNE. Journal of Machine Learning Research 9: 2579–
best friend, Chao-Wei Chen, for his suggestion and encour-               2605.
agement.                                                                 Nasrullah, Z.; and Zhao, Y. 2019. Music Artist Classifica-
                                                                         tion with Convolutional Recurrent Neural Networks. In In-
                        References                                       ternational Joint Conference on Neural Networks (IJCNN),
Bahdanau, D.; Cho, K.; and Bengio, Y. 2015. Neural Ma-                   1–8.
chine Translation by Jointly Learning to Align and Trans-                Rafii, Z.; Liutkus, A.; Stöter, F.-R.; Mimilakis, S. I.; FitzGer-
late. In International Conference on Learning Representa-                ald, D.; and Pardo, B. 2018. An Overview of Lead and
tions (ICLR), 1–15.                                                      Accompaniment Separation in Music. IEEE/ACM Trans-
Bello, I.; Zoph, B.; Vaswani, A.; Shlens, J.; and Le, Q. V.              actions on Audio, Speech and Language Processing 26(8):
2019. Attention Augmented Convolutional Networks. In                     1307–1335.


                                                                  8224
Ramachandran, P.; Parmar, N.; Vaswani, A.; Bello, I.; Lev-
skaya, A.; and Shlens, J. 2019. Stand-Alone Self-Attention
in Vision Models. In Advances in Neural Information Pro-
cessing Systems (NeurIPS), 68–80.
Roy, A. G.; Navab, N.; and Wachinger, C. 2019. Recalibrat-
ing Fully Convolutional Networks With Spatial and Channel
‘Squeeze & Excitation’ Blocks. IEEE Transactions on Med-
ical Imaging 38(2): 540–549.
Sharma, B.; Das, R. K.; and Li, H. 2019. On the Impor-
tance of Audio-Source Separation for Singer Identification
in Polyphonic Music. In Proc. Interspeech, 2020–2024.
Stöter, F.-R.; Uhlich, S.; Liutkus, A.; and Mitsufuji, Y.
2019. Open-Unmix - A Reference Implementation for Mu-
sic Source Separation. Journal of Open Source Software
[Online] https://sigsep.github.io/open-unmix/.
Sturm, B. L. 2014. A Simple Method to Determine if a Mu-
sic Information Retrieval System is a “Horse”. IEEE Trans-
actions on Multimedia 16(6): 1636–1644.
Sundberg, J. 1989. The Science of the Singing Voice. North-
ern Illinois University Press.
Van, T. P.; Quang, N. T. N.; and Thanh, T. M. 2019. Deep
Learning Approach for Singer Voice Classification of Viet-
namese Popular Music. In Proceedings of the Tenth Interna-
tional Symposium on Information and Communication Tech-
nology (SoICT), 255–260.
Vaswani, A.; Shazeer, N.; Parmar, N.; Uszkoreit, J.; Jones,
L.; Gomez, A. N.; Kaiser, Ł.; and Polosukhin, I. 2017. At-
tention is All you Need. In Advances in Neural Information
Processing Systems (NeurIPS), 5998–6008.
Wang, X.; Girshick, R.; Gupta, A.; and He, K. 2018. Non-
Local Neural Networks. In Proceedings of the IEEE Confer-
ence on Computer Vision and Pattern Recognition (CVPR),
7794–7803.
Woo, S.; Park, J.; Lee, J.-Y.; and Kweon, I. S. 2018. CBAM:
Convolutional Block Attention Module. In Proceedings of
the European Conference on Computer Vision (ECCV), 3–
19.
Yu, A. W.; Dohan, D.; Luong, M.-T.; Zhao, R.; Chen, K.;
Norouzi, M.; and Le, Q. V. 2018. QANet: Combining Local
Convolution with Global Self-Attention for Reading Com-
prehension. In International Conference on Learning Rep-
resentations (ICLR).
Zhang, H.; Goodfellow, I.; Metaxas, D.; and Odena, A.
2019. Self-Attention Generative Adversarial Networks. In
Proceedings of the 36th International Conference on Ma-
chine Learning (ICML), 7354–7363.
Zhang, X.; Gao, Y.; Yu, Y.; and Li, W. 2020. Music Artist
Classification with WaveNet Classifier for Raw Waveform
Audio Data. arXiv preprint arXiv:2004.04371v1 .




                                                              8225
