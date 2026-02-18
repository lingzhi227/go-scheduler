::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_main
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_content
# Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs {#efficient-pre-training-of-llms-via-topology-aware-communication-alignment-on-more-than-9600-gpus .ltx_title .ltx_title_document}

::: ltx_authors
[ [Guoliang He    ^1^, Youhe Jiang[^1^[[^1^[footnotemark: ]{.ltx_note_type}[1]{.ltx_tag .ltx_tag_note}]{.ltx_note_content}]{.ltx_note_outer}]{#footnotex1 .ltx_note .ltx_role_footnotemark}    ^1^, Wencong Xiao^2^, Kaihua Jiang^2^, Shuguang Wang^2^,\
[Jun Wang^2^, Zixian Du^2^, Zhuo Jiang^2^, Xinlei Zhang^2^, Binhang Yuan^3^, Eiko Yoneki^1^\
^[1]{.ltx_text .ltx_font_medium}^]{.ltx_text .ltx_font_bold}University of Cambridge, ^2^ByteDance Seed, ^3^HKUST\
[{gh512, yj367}@cam.ac.uk, biyuan@ust.hk, eiko.yoneki@cl.cam.ac.uk,\
{hanli.hl, jiangkaihua, wangshuguang, wangjun.289}@bytedance.com,\
{duzixian, jiangzhuo.cs, zhangxinlei.123}@bytedance.com]{.ltx_text .ltx_font_typewriter} ]{.ltx_personname}[Equal contribution.Work done during internship at ByteDance Seed.Correspondence to: Youhe Jiang \<yj367@cam.ac.uk\>.]{.ltx_author_notes}]{.ltx_creator .ltx_role_author}
:::

::: ltx_abstract
###### Abstract {#abstract .ltx_title .ltx_title_abstract}

The scaling law for large language models (LLMs) depicts that the path towards machine intelligence necessitates training at large scale. Thus, companies continuously build large-scale GPU clusters, and launch training jobs that span over thousands of computing nodes. However, LLM pre-training presents unique challenges due to its complex communication patterns, where GPUs exchange data in sparse yet high-volume bursts within specific groups. Inefficient resource scheduling exacerbates bandwidth contention, leading to suboptimal training performance. This paper presents Arnold, a scheduling system summarizing our experience to effectively align LLM communication patterns with data center topology at scale. An in-depth characteristic study is performed to identify the impact of physical network topology to LLM pre-training jobs. Based on the insights, we develop a scheduling algorithm to effectively align communication patterns with the physical network topology in modern data centers. Through simulation experiments, we show the effectiveness of our algorithm in reducing the maximum spread of communication groups by up to $1.67$x. In production training, our scheduling system improves the end-to-end performance by $10.6\%$ when training with more than $9600$ GPUs, a significant improvement for our training pipeline.
:::

::::::::::::: {#S1 .section .ltx_section}
## [1 ]{.ltx_tag .ltx_tag_section}Introduction {#introduction .ltx_title .ltx_title_section}

::: {#S1.p1 .ltx_para}
Pre-training large language models (LLMs) at scale is a highly resource-intensive process that requires vast computational infrastructure. The performance of LLM training is fundamentally dependent on three factors: dataset size, computational power, and model parameters [scaling_law](#bib.bib19){.ltx_ref} . To meet these demands, companies continually enhance their computing infrastructure by incorporating cutting-edge GPUs and redesigning network architectures [imbue_infrastructure](#bib.bib14){.ltx_ref} ; [ali_hpn](#bib.bib33){.ltx_ref} ; [rail_only](#bib.bib41){.ltx_ref} . However, LLM pre-training presents unique challenges that distinguish it from conventional deep learning tasks --- in this paper, we explore [how to develop an efficient resource scheduling mechanism to support the LLM training workflow to accommodate the resource-intensive and complex communication patterns in modern data centers]{.ltx_text .ltx_font_italic}.
:::

::: {#S1.p2 .ltx_para}
LLM pre-training is an exceptionally resource-intensive process. Given the pressing need to commercialize LLMs swiftly, accelerating the training process is paramount. However, training these models often spans weeks, requiring the deployment of thousands of GPU nodes per run. The ability to efficiently schedule and allocate resources is critical for both performance optimization and cost management. Furthermore, the unique data transmission patterns in LLM training --- wherein GPUs communicate sparsely but at a high volume within specific groups --- pose an additional challenge in leveraging modern multi-tier network topologies effectively.
:::

::: {#S1.p3 .ltx_para}
Existing cluster schedulers (e.g., [mlaas](#bib.bib43){.ltx_ref} ; [fgd](#bib.bib44){.ltx_ref} ; [gandiva](#bib.bib45){.ltx_ref} ; [antman](#bib.bib46){.ltx_ref} ; [mast](#bib.bib7){.ltx_ref} ; [crux](#bib.bib5){.ltx_ref} ) fail to integrate network topology-aware scheduling specific to LLM workloads. The primary limitation is their lack of awareness of the high-volume, yet sparsely active distributed communication patterns inherent in LLM training. For example, Figure [[1(a)]{.ltx_text .ltx_ref_tag}](#S1.F1.sf1 "In Figure 1 ‣ 1 Introduction ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} indicates that 30% - 50% of the time is spent on communication during production LLM training, but studies [rail_only](#bib.bib41){.ltx_ref} show that more than $99\%$ of the GPU pairs do not exhibit direct traffic, with data exchange occurring exclusively within specific communication groups, as shown in Figure [[1(b)]{.ltx_text .ltx_ref_tag}](#S1.F1.sf2 "In Figure 1 ‣ 1 Introduction ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}. Meanwhile, modern GPU clusters use multi-tier, fat-tree network topologies [flat_tree](#bib.bib1){.ltx_ref} (Figure [[2(b)]{.ltx_text .ltx_ref_tag}](#S2.F2.sf2 "In Figure 2 ‣ 2 Background ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}), and inefficient job placement leads to significant bandwidth loss and communication overhead. Current schedulers are not designed to optimize network-aware placement at the scale required by LLM pre-training jobs (LPJs).
:::

<figure id="S1.F1" class="ltx_figure ltx_align_floatright">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S1.F1.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/images/intro1.png" id="S1.F1.sf1.g1" class="ltx_graphics ltx_img_square" width="598" height="539" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(a)</span> </span><span class="ltx_text" style="font-size:90%;">Time break-down.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S1.F1.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/images/intro2.png" id="S1.F1.sf2.g1" class="ltx_graphics ltx_img_landscape" width="598" height="436" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(b)</span> </span><span class="ltx_text" style="font-size:90%;">Traffic heatmap.</span></figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 1</span>: </span><span class="ltx_text" style="font-size:90%;">Communication characteristics of LLMs training in production.</span></figcaption>
</figure>

::: {#S1.p4 .ltx_para}
To enable effective scheduling of LPJs in data centers, we identify two key challenges that limit the effectiveness of existing cluster schedulers.
:::

::: {#S1.p5 .ltx_para}
- [[•]{.ltx_tag .ltx_tag_item}]{#S1.I1.i1}

  ::: {#S1.I1.i1.p1 .ltx_para}
  [Misalignment of communication patterns with data center topology.]{.ltx_text .ltx_font_bold} Schedulers optimize GPU locality through bin-packing but lack awareness of LPJ communication structures. Consequently, jobs may experience inefficient cross-switch communication despite being tightly packed.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S1.I1.i2}

  ::: {#S1.I1.i2.p1 .ltx_para}
  [Unaddressed trade-offs across communication dimensions.]{.ltx_text .ltx_font_bold} There is a fundamental trade-off between aligning data parallel (DP) and pipeline parallel (PP) communication patterns. Since GPUs participate in both groups, perfect alignment for both is unachievable, and schedulers must carefully balance the two during placement.
  :::
:::

::: {#S1.p6 .ltx_para}
To address the challenges, we present Arnold, a system that co-designs training frameworks and cluster scheduling, effectively aligning LPJs with modern data center network topology. To optimize training performance, we performed an in-depth characterization study to investigate the impact of physical network topology on LLM training. Based on the observation, we devise a scheduling algorithm to reduce the maximum weighted spread of communication groups for LPJs. We also develop a resource management policy that manages job queues to reserve nodes for imminent LPJs.
:::

::: {#S1.p7 .ltx_para}
Through trace-based experiments, we show the effectiveness of our scheduling algorithm by benchmarking against other SOTA algorithms. We also perform a production run with $9600$+ GPUs and show our proposed system improves the end-to-end training performance by $10.6\%$. In summary, our contributions include the following:
:::

::: {#S1.p8 .ltx_para}
[Contribution 1.]{.ltx_text .ltx_font_bold .ltx_framed .ltx_framed_underline} We identify the challenge of aligning LLM communication patterns with modern data center topology for large-scale pre-training. We characterize the impact of physical network topology on individual communication operations and end-to-end training performance in modern data centers.
:::

::: {#S1.p9 .ltx_para}
[Contribution 2.]{.ltx_text .ltx_font_bold .ltx_framed .ltx_framed_underline} We design a scheduling algorithm to effectively align communication patterns to the topology of the data center for LPJs, and a resource management policy to reserve nodes for the placement.
:::

::: {#S1.p10 .ltx_para}
[Contribution 3.]{.ltx_text .ltx_font_bold .ltx_framed .ltx_framed_underline} Through comprehensive simulation experiments, we evaluate the effectiveness of the scheduling algorithm in reducing the maximum spread of communication groups by up to $1.67$x. In a production run, we verify the proposed scheduler can improve a $9600$+ GPUs LPJ by $10.6\%$.
:::
:::::::::::::

:::::::: {#S2 .section .ltx_section}
## [2 ]{.ltx_tag .ltx_tag_section}Background {#background .ltx_title .ltx_title_section}

::: {#S2.p1 .ltx_para}
[Distributed training.]{.ltx_text .ltx_font_bold} LLMs are billion-parameter transformer-based models that must be trained with multi-GPU systems[attention](#bib.bib40){.ltx_ref} ; [moe](#bib.bib36){.ltx_ref} . Common training frameworks [megatron1](#bib.bib37){.ltx_ref} ; [deepspeed](#bib.bib35){.ltx_ref} employ hybrid parallelization strategies to parallelize and accelerate the training process, including:
:::

::: {#S2.p2 .ltx_para}
- [[•]{.ltx_tag .ltx_tag_item}]{#S2.I1.i1}

  ::: {#S2.I1.i1.p1 .ltx_para}
  Data parallelism (DP). The Zero Redundancy Optimizer (ZeRO) [zero](#bib.bib34){.ltx_ref} shards model weights and gradients across data parallel processes and performs synchronization at the end of a training step by all-gather and reduce-scatter communication operations.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S2.I1.i2}

  ::: {#S2.I1.i2.p1 .ltx_para}
  Pipeline parallelism (PP). The layers of models are divided into several stages (PP size), and each stage interleave communication to adjacent stages as well as the computation within the stages. Inter-stage communication is performed by P2P communication operation like send-recv.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S2.I1.i3}

  ::: {#S2.I1.i3.p1 .ltx_para}
  Tensor parallelism (TP). Model weights within an PP stage are further sharded across multi-GPUs to alleviate the memory pressure. All-gather and reduce-scatter are necessary to synchronize the intermediate activation during forward pass and backward pass.
  :::
:::

::: {#S2.p3 .ltx_para}
The combination of parallelism, i.e. hybrid parallelism, forms diverse communication patterns for GPUs, and training frameworks use communication group to manage the complexity. Each GPU is assigned to a DP, TP, and PP communication group at initialization. The illustrations of different parallelisms and communication groups are demonstrated in Figure [[2(a)]{.ltx_text .ltx_ref_tag}](#S2.F2.sf1 "In Figure 2 ‣ 2 Background ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}.
:::

::: {#S2.p4 .ltx_para}
Previous works [megatron](#bib.bib24){.ltx_ref} ; [megatron1](#bib.bib37){.ltx_ref} have identified that TP communication groups should be prioritized to GPUs located within the same node to utilize the high-bandwidth NVLink interconnection due to stringent data dependencies. Thus, the inter-node communication only takes place within the DP and PP groups. As only inter-node communication is sensitive to physical network topology, the communication patterns of DP group and PP groups are the focus of this paper.
:::

<figure id="S2.F2" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S2.F2.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/x1.png" id="S2.F2.sf1.g1" class="ltx_graphics ltx_img_square" width="830" height="679" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(a)</span> </span><span class="ltx_text" style="font-size:90%;">Hybrid parallelism.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S2.F2.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/x2.png" id="S2.F2.sf2.g1" class="ltx_graphics ltx_img_landscape" width="830" height="617" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(b)</span> </span><span class="ltx_text" style="font-size:90%;">Data center topology.</span></figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 2</span>: </span><span class="ltx_text" style="font-size:90%;">LLMs parallelism and data center topology.</span></figcaption>
</figure>

::: {#S2.p5 .ltx_para}
[Data center topology.]{.ltx_text .ltx_font_bold} Figure [[2(b)]{.ltx_text .ltx_ref_tag}](#S2.F2.sf2 "In Figure 2 ‣ 2 Background ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} gives an overview of our HPC cluster, which is similar to other modern data centers. More than $2000$ nodes are interconnected by three layers of switches, forming a CLOS-like topology [clos](#bib.bib8){.ltx_ref} . The leaf switch is denoted as s0, which interconnects nodes within the same rack. Then, several s0 switches link to a spine switch (s1), forming a minipod of nodes. Finally, s1 switches link to core switches, enabling communication between minipods. The switches in each layer have $32$ ports both for uplinks and downlinks. The greatest number of hops occurs when the nodes of two different minipods communicate with each other. The compute nodes are equipped with $8$ H800 GPUs, each of which is connected to an InfiniBand [InfiniBand](#bib.bib27){.ltx_ref} NIC. GPUs within a node are connected by high-bandwidth links such as NVLink [NVLink](#bib.bib30){.ltx_ref} with a bandwidth of 400Gbps, while inter-node communication is achieved via the InfiniBand network.
:::
::::::::

::::::: {#S3 .section .ltx_section}
## [3 ]{.ltx_tag .ltx_tag_section}Observation and Challenges {#observation-and-challenges .ltx_title .ltx_title_section}

::: {#S3.p1 .ltx_para}
Given a user-specified number of GPUs and degree of hybrid parallelism of an LPJ, job scheduling systems enqueue the job and perform resource scheduling to find the best placement in our GPU cluster. However, we observe existing scheduling systems fail to align LLM communication patterns with data center topology in practice.
:::

<figure id="S3.F3" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S3.F3.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/images/align-non.png" id="S3.F3.sf1.g1" class="ltx_graphics ltx_img_landscape" width="476" height="289" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(a)</span> </span><span class="ltx_text" style="font-size:90%;">Misaligned.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S3.F3.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/images/align-pp.png" id="S3.F3.sf2.g1" class="ltx_graphics ltx_img_landscape" width="476" height="311" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(b)</span> </span><span class="ltx_text" style="font-size:90%;">PP-aligned.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S3.F3.sf3" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/images/align-dp.png" id="S3.F3.sf3.g1" class="ltx_graphics ltx_img_landscape" width="476" height="296" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(c)</span> </span><span class="ltx_text" style="font-size:90%;">DP-aligned.</span></figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 3</span>: </span><span class="ltx_text" style="font-size:90%;">Alignment of communication patterns. One DP group and PP group is highlighted.</span></figcaption>
</figure>

::: {#S3.p2 .ltx_para}
[Observation 1: Misalignment of job placement results in increased cross-switch communication.]{.ltx_text .ltx_font_bold .ltx_font_italic} SOTA cluster schedulers [mlaas](#bib.bib43){.ltx_ref} ; [topo_aware](#bib.bib2){.ltx_ref} ; [gandiva](#bib.bib45){.ltx_ref} apply a bin-packing strategy to enhance GPU locality of LPJs. However, as shown in Figure [[3(a)]{.ltx_text .ltx_ref_tag}](#S3.F3.sf1 "In Figure 3 ‣ 3 Observation and Challenges ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}, even if the scheduler packs an 4-node (32 GPUs) LPJ inside a minipod, the communication groups may still not be aligned, because both DP and PP groups engage cross-spine-switch communication that has a longer distance. This misalignment stems from the scheduler's lack of awareness of the LPJ's communication structure at scheduling time, limiting its ability to allocate GPU resources according to the job's communication patterns.
:::

::: {#S3.p3 .ltx_para}
[Observation 2: Unresolved trade-offs between DP and PP communication priorities.]{.ltx_text .ltx_font_bold .ltx_font_italic} Figure [[3(b)]{.ltx_text .ltx_ref_tag}](#S3.F3.sf2 "In Figure 3 ‣ 3 Observation and Challenges ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} and [[3(c)]{.ltx_text .ltx_ref_tag}](#S3.F3.sf3 "In Figure 3 ‣ 3 Observation and Challenges ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} show two potential alignments of the LPJ, with one prioritizing DP communication and the other prioritizing PP. This presents a fundamental trade-off between the two, because DP and PP are orthogonal parallelism strategies widely used in LLM training. Each GPU participates in both a DP and an PP group, making it impossible to perfectly align both communication patterns simultaneously. A well-designed scheduler must consider this trade-off and balance the alignment needs of both group types during job placement.
:::

::: {#S3.p4 .ltx_para}
[Challenge: Communication and topology-aligned scheduling for LPJs.]{.ltx_text .ltx_font_bold .ltx_font_italic} To effectively schedule LPJs, the scheduler must be aware of the diverse communication patterns and minimize their spread in data centers. Furthermore, effective balance between the spread of DP and PP groups is critical, which requires in-depth characterization of communication patterns in modern data centers.
:::
:::::::

:::::::::::: {#S4 .section .ltx_section}
## [4 ]{.ltx_tag .ltx_tag_section}Characterization of Communication Patterns for LPJs {#characterization-of-communication-patterns-for-lpjs .ltx_title .ltx_title_section}

::: {#S4.p1 .ltx_para}
Although prior studies [gandiva](#bib.bib45){.ltx_ref} ; [topo_aware](#bib.bib2){.ltx_ref} have explored locality and topology, their scope is constrained by (1) a focus on data-parallel (all-reduce) workloads and (2) limited consideration of inter-node topology. To address these gaps, we conduct NCCL tests, a benchmarking suite designed to measure the latency and bandwidth of communication operations used by NCCL [nccl](#bib.bib28){.ltx_ref} ; [nccl_test](#bib.bib29){.ltx_ref} , to study the impact of inter-node topology. We focus on inter-node topology across minipods, as the scale of LPJs typically necessitates allocating computing nodes across multiple minipods, where the slowest communication path often dictates the overall communication overhead.
:::

::: {#S4.p2 .ltx_para}
[Communication operation performance.]{.ltx_text .ltx_font_bold} Figure [[4]{.ltx_text .ltx_ref_tag}](#S4.F4 "Figure 4 ‣ 4 Characterization of Communication Patterns for LPJs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} studies the performance of communication operations. We use the bus bandwidth (BusBw) as a performance measurement, which reflects the peak hardware bandwidth by accounting for the number of ranks for collective communication.
:::

<figure id="S4.F4" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S4.F4.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/x3.png" id="S4.F4.sf1.g1" class="ltx_graphics ltx_img_landscape" width="660" height="376" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(a)</span> </span><span class="ltx_text" style="font-size:90%;">BusBw over message sizes.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S4.F4.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/x4.png" id="S4.F4.sf2.g1" class="ltx_graphics ltx_img_landscape" width="661" height="402" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(b)</span> </span><span class="ltx_text" style="font-size:90%;">BusBw of collective operations.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S4.F4.sf3" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/x5.png" id="S4.F4.sf3.g1" class="ltx_graphics ltx_img_landscape" width="661" height="413" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(c)</span> </span><span class="ltx_text" style="font-size:90%;">BusBw of send-recv.</span></figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 4</span>: </span><span class="ltx_text" style="font-size:90%;">Performance of communication operation. (AR: all-reduce, AG: all-gather, RS: reduce-scatter, SR: send-recv)</span></figcaption>
</figure>

::: {#S4.p3 .ltx_para}
Figure [[4(a)]{.ltx_text .ltx_ref_tag}](#S4.F4.sf1 "In Figure 4 ‣ 4 Characterization of Communication Patterns for LPJs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} illustrates the performance of different communication patterns in message sizes. For collective communication, the message size must be larger than $2^{8}$ ($\approx 256$) megabytes to fully utilize the bandwidth, while for P2P communication like send-recv, a small message size ($\approx 2$ megabytes) can saturate the bandwidth. Using a widely adopted analytical model detailed in Appendix §[[C]{.ltx_text .ltx_ref_tag}](#A3 "Appendix C Analytical Estimation for Communication Volume ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} and substituting the parameters with a 7B GPT-based model, we can obtain that the data volumes of the DP and PP groups are $2$ GB and $30$ MB respectively, indicating that the bandwidth is fully utilized.
:::

::: {#S4.p4 .ltx_para}
The degradation in BusBw by expanding communication groups across minipods is illustrated in Figures [[4(b)]{.ltx_text .ltx_ref_tag}](#S4.F4.sf2 "In Figure 4 ‣ 4 Characterization of Communication Patterns for LPJs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} and [[4(c)]{.ltx_text .ltx_ref_tag}](#S4.F4.sf3 "In Figure 4 ‣ 4 Characterization of Communication Patterns for LPJs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}. Performance decreases by up to 17% for collective operations and 70% for the P2P operation as communication extends over additional minipods, highlighting the critical importance of GPU locality and alignment. Additionally, our findings suggest that co-located jobs may experience reduced bandwidth contention in multi-tenant cluster environments, as evidenced by the inter-job interference patterns documented in Appendix [[D]{.ltx_text .ltx_ref_tag}](#A4 "Appendix D Sensitivity to Shared Load ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}.
:::

::: {#S4.p5 .ltx_para}
[End-to-end training performance.]{.ltx_text .ltx_font_bold} Based on the characterization of individual communication operations, we proportionally down-scaled a production model and ran the workload with $96$ GPUs, spanning $2$ minipods, to further understand the impact of network topology on LLM training.
:::

<figure id="S4.F5" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S4.F5.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/x6.png" id="S4.F5.sf1.g1" class="ltx_graphics ltx_img_landscape" width="830" height="301" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(a)</span> </span><span class="ltx_text" style="font-size:90%;">Comparison of three different placement strategies for two types of LLMs.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S4.F5.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/x7.png" id="S4.F5.sf2.g1" class="ltx_graphics ltx_img_landscape" width="664" height="332" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(b)</span> </span><span class="ltx_text" style="font-size:90%;">Performance improvement by scaling model sizes under the optimized alignment.</span></figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 5</span>: </span><span class="ltx_text" style="font-size:90%;">End-to-end training performance.</span></figcaption>
</figure>

::: {#S4.p6 .ltx_para}
Figure [[5(a)]{.ltx_text .ltx_ref_tag}](#S4.F5.sf1 "In Figure 5 ‣ 4 Characterization of Communication Patterns for LPJs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} shows the throughput of the LPJ under the three different placement strategies. The dp-aligned placement is illustrated in Figure [[12]{.ltx_text .ltx_ref_tag}](#A6.F12 "Figure 12 ‣ Appendix F Communication Matrix ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}. The throughput becomes stable after $200$ steps except for a slight fluctuation around the $550$th step due to garbage collection. PP-aligned placement consistently outperforms the other two, demonstrating that prioritizing PP group communication leads to improved performance. The average improvement for the dense model and the MoE model is $2.3\%$ and $1.8\%$ respectively. For the dense model, we also observe that the PP communication dominates the communication overhead, since prioritizing the placement of DP groups leads to no speedup. For the MoE model, reducing the spread of both the DP and PP groups contributes to performance gains, with the optimizations of the PP group providing more improvements.
:::

::: {#S4.p7 .ltx_para}
Figure [[5(b)]{.ltx_text .ltx_ref_tag}](#S4.F5.sf2 "In Figure 5 ‣ 4 Characterization of Communication Patterns for LPJs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} shows that if we scale the model size by adding more layers, the performance improvement continues to increase. We attribute this to communication being the primary performance bottleneck, with larger models further amplifying communication volume. Thus, more pronounced benefits are obtained as the model size increases.
:::

::: {#S4.p8 .ltx_para}
We further examine the sensitivity of training performance to intra-minipod network topology by varying node placement within a single minipod. For a dense 24B parameter model, the maximum observed performance variation is 0.3%, and the impact is negligible for other models. Since the communication overhead of a group is typically dominated by the slowest link, and LPJ communication groups frequently span multiple minipods, we conclude that training performance is largely insensitive to intra-minipod topology.
:::

::: {#S4.p9 .ltx_para}
We repeated the characterization experiment in another GPU cluster detailed in Appendix [[E]{.ltx_text .ltx_ref_tag}](#A5 "Appendix E Ada Lovelace GPUs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}, and found that the best placement can be subjected to model sizes and GPU types. Since LPJs are typically scheduled in advance and deployed for a long duration, it is essential to perform a characterization beforehand to identify communication bottlenecks within the group. This enables the optimization of placement strategies accordingly and balances the trade-off, as detailed in §[[5]{.ltx_text .ltx_ref_tag}](#S5 "5 Scheduling Algorithm ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}.
:::
::::::::::::

::::::::::::::::::::::: {#S5 .section .ltx_section}
## [5 ]{.ltx_tag .ltx_tag_section}Scheduling Algorithm {#scheduling-algorithm .ltx_title .ltx_title_section}

::: {#S5.p1 .ltx_para}
The core component of Arnold is its scheduling module, which is designed to effectively allocate GPUs to LPJs based on the user-specified number of GPUs and the degree of hybrid parallelism. In this section, we formally define the scheduling problem and present our proposed solution.
:::

:::::: {#S5.SS1 .section .ltx_subsection}
### [5.1 ]{.ltx_tag .ltx_tag_subsection}Workload Representation {#workload-representation .ltx_title .ltx_title_subsection}

::: {#S5.SS1.p1 .ltx_para}
Arnold represents an LPJ by a communication matrix, where a row represents a PP group and a column represents a DP group. Formally, given a job specification including the total number of GPUs, the degree of PP, TP, and then Arnold computes the size of the communication matrix using Equation [[1]{.ltx_text .ltx_ref_tag}](#S5.E1 "In 5.1 Workload Representation ‣ 5 Scheduling Algorithm ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}.
:::

::: {#S5.SS1.p2 .ltx_para}
  -- ----------------------------------- -- ----------------------------------------------------
     $$\begin{split}&DP=\#GPUs/TP/PP\\      [(1)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
     &\#row=DP/(8/TP)\\                     
     &\#col=PP\end{split}$$                 
  -- ----------------------------------- -- ----------------------------------------------------
:::

::: {#S5.SS1.p3 .ltx_para}
An example of $96$ GPUs and $DP=6,PP=2$ is shown in Figure [[12]{.ltx_text .ltx_ref_tag}](#A6.F12 "Figure 12 ‣ Appendix F Communication Matrix ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}. For a node $v_{ij}$ in the communication matrix, it is attached with a vector $[v_{w},v_{d},v_{p}]$, representing the size of weight, the DP and PP communication volume per GPU, respectively. Those values are computed using the analytical model detailed in Appendix §[[C]{.ltx_text .ltx_ref_tag}](#A3 "Appendix C Analytical Estimation for Communication Volume ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}, and is used to balance the trade-off between DP and PP groups.
:::
::::::

:::::::::::::: {#S5.SS2 .section .ltx_subsection}
### [5.2 ]{.ltx_tag .ltx_tag_subsection}Objectives {#objectives .ltx_title .ltx_title_subsection}

::: {#S5.SS2.p1 .ltx_para}
The scheduling objective function aims to minimize the maximum spread for both DP and PP groups. For a node $v_{ij}$ in the communication matrix, the number of possible assignment to a minipod is $k$. Therefore, it has a one-hot decision vector $x_{ij}$ of length $k$, representing the decision to place to one of the $k$ minipods. Then the objective function can be written as Equation [[2]{.ltx_text .ltx_ref_tag}](#S5.E2 "In 5.2 Objectives ‣ 5 Scheduling Algorithm ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}.
:::

::: {#S5.SS2.p2 .ltx_para}
  -- --------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
     $$\text{{MIN}}\,\,\,[\alpha\max_{i}(D(x_{i1},x_{i2},...x_{in}))+\beta\max_{j}(D(x_{1j},x_{2j},...,x_{mj}))]$$      [(2)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- --------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
:::

::: {#S5.SS2.p3 .ltx_para}
Where, the distance $D$ between the $n$ vectors is defined as follows:
:::

::: {#S5.SS2.p4 .ltx_para}
  -- ---------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
     $$D(v_{1},v_{2},\dots,v_{n})=|\{i:\exists j,l\in\{1,2,\dots,n\},j\neq l,v_{j}[i]\neq v_{l}[i]\}|$$      [(3)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- ---------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
:::

::: {#S5.SS2.p5 .ltx_para}
Intuitively, the distance measures the spread of a communication group, i.e. if any two vectors differ in the $i$th position, the $i$th position adds one to the distance. The objective function aims to minimize the weighted sum of the maximum spread of DP and PP groups, and the maximum is taken because the slowest communication group is the straggler to slow down the whole training process. The weight parameters $\alpha$ and $\beta$ represent the affinity that controls the trade-off between DP and PP groups, and $\alpha+\beta=1$. Equation [[2]{.ltx_text .ltx_ref_tag}](#S5.E2 "In 5.2 Objectives ‣ 5 Scheduling Algorithm ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} cannot be solved by off-the-shelf solvers efficiently for online scheduling due to the discrete distance calculation, and we seek for simplifications.
:::

<figure id="S5.SS2.fig1" class="ltx_figure">
<table id="A9.EGx1" class="ltx_equationgroup ltx_eqn_align ltx_eqn_table">
<tbody id="S5.E4">
<tr class="ltx_equation ltx_eqn_row ltx_align_baseline">
<td class="ltx_eqn_cell ltx_eqn_center_padleft"></td>
<td class="ltx_td ltx_align_right ltx_eqn_cell"><span class="ltx_text ltx_markedasmath ltx_font_italic" style="font-size:90%;">MIN</span></td>
<td class="ltx_td ltx_align_left ltx_eqn_cell"><span class="math inline">[<em>α</em>∑<sub><em>j</em></sub>(<em>y</em><sub><em>j</em></sub>) + <em>β</em><em>T</em>]</span></td>
<td class="ltx_eqn_cell ltx_eqn_center_padright"></td>
<td class="ltx_eqn_cell ltx_eqn_eqno ltx_align_middle ltx_align_right"><span class="ltx_tag ltx_tag_equation ltx_align_right">(4)</span></td>
</tr>
</tbody>
<tbody id="S5.E5">
<tr class="ltx_equation ltx_eqn_row ltx_align_baseline">
<td class="ltx_eqn_cell ltx_eqn_center_padleft"></td>
<td class="ltx_td ltx_align_right ltx_eqn_cell"><span class="ltx_text ltx_markedasmath" style="font-size:90%;">s.t.</span></td>
<td class="ltx_td ltx_align_left ltx_eqn_cell"><span class="math inline">∀<em>i</em> : ∑<sub><em>j</em></sub><em>s</em><sub><em>i</em><em>j</em></sub> ≤ <em>T</em>   (Max Spread)</span></td>
<td class="ltx_eqn_cell ltx_eqn_center_padright"></td>
<td class="ltx_eqn_cell ltx_eqn_eqno ltx_align_middle ltx_align_right"><span class="ltx_tag ltx_tag_equation ltx_align_right">(5)</span></td>
</tr>
</tbody>
<tbody id="S5.E6">
<tr class="ltx_equation ltx_eqn_row ltx_align_baseline">
<td class="ltx_eqn_cell ltx_eqn_center_padleft"></td>
<td class="ltx_td ltx_align_right ltx_eqn_cell"><span class="math inline">∀<em>j</em>:</span></td>
<td class="ltx_td ltx_align_left ltx_eqn_cell"><span class="math inline">∑<sub><em>i</em></sub><em>p</em><sub><em>i</em><em>j</em></sub> ≤ <em>c</em><sub><em>j</em></sub><em>y</em><sub><em>j</em></sub>   (Capacity Const.)</span></td>
<td class="ltx_eqn_cell ltx_eqn_center_padright"></td>
<td class="ltx_eqn_cell ltx_eqn_eqno ltx_align_middle ltx_align_right"><span class="ltx_tag ltx_tag_equation ltx_align_right">(6)</span></td>
</tr>
</tbody>
<tbody id="S5.E7">
<tr class="ltx_equation ltx_eqn_row ltx_align_baseline">
<td class="ltx_eqn_cell ltx_eqn_center_padleft"></td>
<td class="ltx_td ltx_align_right ltx_eqn_cell"><span class="math inline">∀<em>i</em>:</span></td>
<td class="ltx_td ltx_align_left ltx_eqn_cell"><span class="math inline">∑<sub><em>j</em></sub><em>p</em><sub><em>i</em><em>j</em></sub> = 1   (Allocation Const.)</span></td>
<td class="ltx_eqn_cell ltx_eqn_center_padright"></td>
<td class="ltx_eqn_cell ltx_eqn_eqno ltx_align_middle ltx_align_right"><span class="ltx_tag ltx_tag_equation ltx_align_right">(7)</span></td>
</tr>
</tbody>
<tbody id="S5.E8">
<tr class="ltx_equation ltx_eqn_row ltx_align_baseline">
<td class="ltx_eqn_cell ltx_eqn_center_padleft"></td>
<td class="ltx_td ltx_align_right ltx_eqn_cell"><span class="math inline">∀<em>i</em>, <em>j</em>:</span></td>
<td class="ltx_td ltx_align_left ltx_eqn_cell"><span class="math inline"><em>p</em><sub><em>i</em><em>j</em></sub> ≤ <em>s</em><sub><em>i</em><em>j</em></sub>   (Minipod Selection)</span></td>
<td class="ltx_eqn_cell ltx_eqn_center_padright"></td>
<td class="ltx_eqn_cell ltx_eqn_eqno ltx_align_middle ltx_align_right"><span class="ltx_tag ltx_tag_equation ltx_align_right">(8)</span></td>
</tr>
</tbody>
<tbody id="S5.E9">
<tr class="ltx_equation ltx_eqn_row ltx_align_baseline">
<td class="ltx_eqn_cell ltx_eqn_center_padleft"></td>
<td class="ltx_td ltx_align_right ltx_eqn_cell"><span class="math inline">∀<em>j</em>:</span></td>
<td class="ltx_td ltx_align_left ltx_eqn_cell"><span class="math inline"><em>y</em><sub><em>j</em></sub> ∈ {0, 1}</span></td>
<td class="ltx_eqn_cell ltx_eqn_center_padright"></td>
<td class="ltx_eqn_cell ltx_eqn_eqno ltx_align_middle ltx_align_right"><span class="ltx_tag ltx_tag_equation ltx_align_right">(9)</span></td>
</tr>
</tbody>
<tbody id="S5.E10">
<tr class="ltx_equation ltx_eqn_row ltx_align_baseline">
<td class="ltx_eqn_cell ltx_eqn_center_padleft"></td>
<td class="ltx_td ltx_align_right ltx_eqn_cell"><span class="math inline">∀<em>i</em>, <em>j</em>:</span></td>
<td class="ltx_td ltx_align_left ltx_eqn_cell"><span class="math inline"><em>s</em><sub><em>i</em><em>j</em></sub> ∈ {0, 1}, <em>p</em><sub><em>i</em><em>j</em></sub> ∈ [0, 1]</span></td>
<td class="ltx_eqn_cell ltx_eqn_center_padright"></td>
<td class="ltx_eqn_cell ltx_eqn_eqno ltx_align_middle ltx_align_right"><span class="ltx_tag ltx_tag_equation ltx_align_right">(10)</span></td>
</tr>
</tbody>
</table>
</figure>

::: {#S5.SS2.p6 .ltx_para}
[Domain-specific simplification.]{.ltx_text .ltx_font_bold} We identify that communication groups are homogeneous and synchronous for LPJs because nodes are gang-scheduled and must synchronize their gradients at the end of a training step. As a result, each PP group always starts approximately at the same time and performs the same amount of computation and communication. Similarly, DP groups perform gradient synchronization at the same time. The characteristics allow us to simplify the scheduling objective function by coarsening a scheduling unit as a communication group. We therefore transform Equation [[2]{.ltx_text .ltx_ref_tag}](#S5.E2 "In 5.2 Objectives ‣ 5 Scheduling Algorithm ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} into a bin-packing-like formulation:
:::

::: {#S5.SS2.p7 .ltx_para}
Where $y_{j}$ indicates whether the $j$-th minipod is used and $c_{j}$ is the normalized capacity of the minipod, updated dynamically based on the number of available nodes. $s_{ij}$ denotes whether the $i$-th communication group is allocated to minipod $j$, $p_{ij}$ denotes the percentage of the $i$-th communication group allocated to minipod $j$. $T$ is an introduced auxiliary variable that allows us to minimize the maximum spread of communication groups. $\alpha$ and $\beta$ are the affinity parameters as before. Overall, minimizing $T$ effectively consolidates communication groups into the smallest possible number of minipods, while the objective term $\sum_{j}(y_{j})$ controls the spread of the other communication group.
:::

::: {#S5.SS2.p8 .ltx_para}
For example, we consider the PP group as a scheduling unit. By setting $\alpha=0$, the scheduler can place each PP group into a minipod, causing more cross-switch communication for DP groups, while by setting $\beta=0$, the scheduler minimizes the overall usage of minipods, although the placement may cause cross-switch communication for PP groups. In this formulation, all variables are either an integer or a fraction. Therefore, the objective function can be solved using off-the-shelf mixed-integer programming (MIP) solvers efficiently for online scheduling [scip](#bib.bib4){.ltx_ref} . After solving the MIP, we assign continuous rank indices to nodes belonging to the same minipod to reduce cross-switch communication within each communication group.
:::

::: {#S5.SS2.p9 .ltx_para}
[Balancing the trade-off.]{.ltx_text .ltx_font_bold} The affinity parameters in Equation [[4]{.ltx_text .ltx_ref_tag}](#S5.E4 "In 5.2 Objectives ‣ 5 Scheduling Algorithm ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} require balance of the trade-off between the DP and PP groups, which depends on the model configurations and GPU types (§[[4]{.ltx_text .ltx_ref_tag}](#S4 "4 Characterization of Communication Patterns for LPJs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}). To perform online scheduling, we store the characterization results in §[[4]{.ltx_text .ltx_ref_tag}](#S4 "4 Characterization of Communication Patterns for LPJs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} to a database, and we search for the best match to determine the values of the affinity parameters. The communication matrix computes the per-GPU parameters ($v_{w}$) and communication volumes ($v_{d},v_{p}$). We then compute the average ratio of computation-to-communication and DP-to-PP volume as $r_{1}=\frac{mb\times v_{w}}{v_{d}+v_{p}}$ and $r_{2}=\frac{v_{d}}{v_{p}}$, where $mb$ is the size of the microbatch. These ratios are used to find the best matching job in the database by comparing the Euclidean distance, i.e. $\text{{MIN}}_{r_{i},r_{j}}\,\,\,\sqrt{(r_{1}-r_{i})^{2}+(r_{2}-r_{j})^{2}}$, because GPUs exhibit comparable performance characteristics if they have similar computational load and communication volume.
:::

::: {#S5.SS2.p10 .ltx_para}
LPJs are associated with metadata $\langle\text{GPU}_{\textit{type}},j_{dp},j_{pp}\rangle$, where $j_{dp},j_{pp}$ corresponds to the improvement of DP-aligned, and PP-aligned placement strategies. The affinity parameters $\alpha$ and $\beta$ are then derived based on the relative performance improvement of $j_{dp}$ and $j_{pp}$, i.e. $\alpha=\frac{j_{dp}}{j_{dp}+j_{pp}}$ and $\beta=\frac{j_{pp}}{j_{dp}+j_{pp}}$.
:::

::: {#S5.SS2.p11 .ltx_para}
Due to the importance of LLM training and their unified architectures, LPJs are scheduled in advanced and pre-characterized, so the profiling data in the database can be looked up in online scheduling. For example, for a $24$b dense model in the H800 GPU cluster, the scheduling unit is set to the PP group and $\alpha$ is set to zero as PP groups clearly dominate the communication overhead. For the $24$b MoE model, $\alpha=0.3$ and $\beta=0.7$.
:::
::::::::::::::

::::: {#S5.SS3 .section .ltx_subsection}
### [5.3 ]{.ltx_tag .ltx_tag_subsection}Resource Management {#resource-management .ltx_title .ltx_title_subsection}

::: {#S5.SS3.p1 .ltx_para}
The scheduling algorithm computes a globally optimal placement for LPJs in the GPU cluster, which inevitably conflicts with other jobs. To address this, we develop a queuing policy to manage the job queue and reserve resources for the imminent LPJ.
:::

::: {#S5.SS3.p2 .ltx_para}
Algorithm [[1]{.ltx_text .ltx_ref_tag}](#alg1 "Algorithm 1 ‣ Appendix G Queue Management ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} illustrates our scheduling policy. Once the LPJ is planned, the scheduler solves the MIP equation and reserves the resources. Since then, incoming jobs are preferentially allocated outside the reserved zone. Otherwise, to improve resource utilization, if the predicted JCT of an incoming job precedes the arrival time of the LPJ, it may still be scheduled within the reserved zone. If neither of the conditions can be satisfied, the job is delayed in the scheduling interval. We also employ an ML-driven job completion time (JCT) predictor to balance the trade-off of queuing delay and resource utilization. The setup and evaluation are detailed in Appendix [[G]{.ltx_text .ltx_ref_tag}](#A7 "Appendix G Queue Management ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} and [[H]{.ltx_text .ltx_ref_tag}](#A8 "Appendix H Evaluation of Queue Management ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} respectively.
:::
:::::
:::::::::::::::::::::::

:::: {#S6 .section .ltx_section}
## [6 ]{.ltx_tag .ltx_tag_section}System Implementation {#system-implementation .ltx_title .ltx_title_section}

<figure id="S6.F6" class="ltx_figure">
<img src="/html/2509.15940/assets/x8.png" id="S6.F6.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="830" height="480" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 6</span>: </span><span class="ltx_text" style="font-size:90%;">Architecture overview of Arnold.</span></figcaption>
</figure>

::: {#S6.p1 .ltx_para}
We have implemented a prototype of Arnold with more than 3k lines of Python codes. The prototype consists of the scheduling module and a trace-driven simulator that can replay production traces. We also have a version of the deployment integrated with Kubernetes [kubernetes2023](#bib.bib21){.ltx_ref} . For training frameworks, we build on top of Megatron [megatron](#bib.bib24){.ltx_ref} and modify it to ensure that communication groups follow the placement provided by Arnold. Figure [[6]{.ltx_text .ltx_ref_tag}](#S6.F6 "Figure 6 ‣ 6 System Implementation ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} gives an overview of Arnold.
:::
::::

:::::::::::::::::: {#S7 .section .ltx_section}
## [7 ]{.ltx_tag .ltx_tag_section}Evaluation {#evaluation .ltx_title .ltx_title_section}

::: {#S7.p1 .ltx_para}
We evaluate Arnold using both simulation and real-cluster experiments. To benchmark scheduling algorithms cost-effectively, we develop a simulator, as direct evaluation on production clusters is prohibitively expensive. After identifying the highest-performing scheduling algorithm through simulation, we deploy it on our production cluster to validate its effectiveness under real workloads.
:::

:::::::::: {#S7.SS1 .section .ltx_subsection}
### [7.1 ]{.ltx_tag .ltx_tag_subsection}Simulation Experiments {#simulation-experiments .ltx_title .ltx_title_subsection}

::: {#S7.SS1.p1 .ltx_para}
[Baselines.]{.ltx_text .ltx_font_bold} We compare the scheduling algorithm with the following baselines.
:::

::: {#S7.SS1.p2 .ltx_para}
1.  [[1.]{.ltx_tag .ltx_tag_item}]{#S7.I1.i1}

    ::: {#S7.I1.i1.p1 .ltx_para}
    Best-fit [best_fit](#bib.bib32){.ltx_ref} assigns the nodes to the minipod with the least remaining resources.
    :::
2.  [[2.]{.ltx_tag .ltx_tag_item}]{#S7.I1.i2}

    ::: {#S7.I1.i2.p1 .ltx_para}
    Random-fit [fgd](#bib.bib44){.ltx_ref} assigns nodes to minipods randomly such that the assignment is balanced and fair.
    :::
3.  [[3.]{.ltx_tag .ltx_tag_item}]{#S7.I1.i3}

    ::: {#S7.I1.i3.p1 .ltx_para}
    GPU-packing [mlaas](#bib.bib43){.ltx_ref} ; [gandiva](#bib.bib45){.ltx_ref} is an effective placement strategy applied by state-of-the-art GPU cluster schedulers that pack multiple jobs to the same GPU. We modify the algorithm to pack multi-GPU jobs to a minipod to satisfy the network topology semantics.
    :::
4.  [[4.]{.ltx_tag .ltx_tag_item}]{#S7.I1.i4}

    ::: {#S7.I1.i4.p1 .ltx_para}
    Topo-aware [topo_aware](#bib.bib2){.ltx_ref} is a GPU topology-aware placement algorithm. It represents the workload as a job graph (similar to our communication matrix) and the topology as a physical graph. Then it recursively bi-partitions the physical graph and maps the job graph to the sub-graphs (Hierarchical Static Mapping Dual Recursive Bi-partitioning [static_map](#bib.bib10){.ltx_ref} ). The graph bi-partitioning is implemented by the Fiduccia Mattheyses algorithm [fm_cut](#bib.bib11){.ltx_ref} .
    :::
:::

::: {#S7.SS1.p3 .ltx_para}
[Metrics.]{.ltx_text .ltx_font_bold} The weighted sum of the maximum spread as in Equation [[2]{.ltx_text .ltx_ref_tag}](#S5.E2 "In 5.2 Objectives ‣ 5 Scheduling Algorithm ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} is used to evaluate scheduling algorithms. To evaluate scalability, we measure the scheduling latency.
:::

<figure id="S7.T1" class="ltx_table">
<table class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<thead class="ltx_thead">
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_r ltx_border_t"><span class="ltx_text ltx_font_bold">Settings</span></th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_r ltx_border_t"><span class="ltx_text ltx_font_bold">Network Topology</span></th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_t"><span class="ltx_text ltx_font_bold">Job Configs</span></th>
</tr>
</thead>
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">(i)</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">3, 18</td>
<td class="ltx_td ltx_align_center ltx_border_t">12, 4, 2</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">(ii)</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">5, 438</td>
<td class="ltx_td ltx_align_center ltx_border_t">24, 4, 8</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_b ltx_border_r ltx_border_t">(iii)</td>
<td class="ltx_td ltx_align_center ltx_border_b ltx_border_r ltx_border_t">11, 1019</td>
<td class="ltx_td ltx_align_center ltx_border_b ltx_border_t">46, 8, 8</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 1</span>: </span><span class="ltx_text" style="font-size:90%;">Benchmark setting. Network topology <span class="math inline">{<em>x</em>}, {<em>y</em>}</span> represent {x} minipod and {y} nodes in total, and the numbers in job configurations are the degree of DP, TP, PP. The scheduling unit is the PP group.</span></figcaption>
</figure>

::: {#S7.SS1.p4 .ltx_para}
[Setups.]{.ltx_text .ltx_font_bold} We use $3$ settings in the benchmark as listed in Table [[1]{.ltx_text .ltx_ref_tag}](#S7.T1 "Table 1 ‣ 7.1 Simulation Experiments ‣ 7 Evaluation ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}, where the network topology is taken from a subset of our GPU cluster, and the job configurations are representative for small, medium, large jobs respectively. We also vary the value of $\alpha$ to investigate different degree of affinity.
:::

<figure id="S7.F7" class="ltx_figure">
<img src="/html/2509.15940/assets/x9.png" id="S7.F7.g1" class="ltx_graphics ltx_img_landscape" width="830" height="330" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 7</span>: </span><span class="ltx_text" style="font-size:90%;">Weighted spread of communication groups under different scheduling algorithms.</span></figcaption>
</figure>

::: {#S7.SS1.p5 .ltx_para}
Figure [[7]{.ltx_text .ltx_ref_tag}](#S7.F7 "Figure 7 ‣ 7.1 Simulation Experiments ‣ 7 Evaluation ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} compares the performance of different algorithms. Our algorithm consistently outperforms other baselines and up to $1.67$x compared to the best baseline. On average, it leads to $1.2$x reduction of the weighted sum of the maximum spread for communication groups. In the simple topology (setting [(i)]{.ltx_text .ltx_font_italic}), our algorithm achieves the same score as best-fit, gpu-pack and topo-aware, because the network topology and job configurations are relatively simple, so there is no room to improve the placement. For medium and large jobs, our algorithm is better than the other baselines due to the large search space of possible placement.
:::

::: {#S7.SS1.p6 .ltx_para}
We also observe that as $\alpha$ increases, our scheduling algorithm is closer to other baselines. This is because $\alpha$ controls the affinity of the DP group, and if $\alpha$ is $1$, the objective function reduces to a bin-packing formulation and therefore has no difference from other bin-packing algorithms. In practice, we would not set $\alpha$ greater than $0.5$ as observed from our characterization results (§[[4]{.ltx_text .ltx_ref_tag}](#S4 "4 Characterization of Communication Patterns for LPJs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}). As a result, our algorithm usually scores higher than other baselines.
:::

<figure id="S7.F8" class="ltx_figure">
<img src="/html/2509.15940/assets/x10.png" id="S7.F8.g1" class="ltx_graphics ltx_img_landscape" width="830" height="335" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 8</span>: </span><span class="ltx_text" style="font-size:90%;">Scheduling latency under different configurations.</span></figcaption>
</figure>

::: {#S7.SS1.p7 .ltx_para}
[Scalability.]{.ltx_text .ltx_font_bold} To evaluate the scalability of our algorithm, we implement an enumeration approach and compare the scheduling latency. Figure [[8]{.ltx_text .ltx_ref_tag}](#S7.F8 "Figure 8 ‣ 7.1 Simulation Experiments ‣ 7 Evaluation ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} investigates the scalability. The enumeration approach is guaranteed to obtain the optimal placement strategy. However, it is not scalable, as it can incur a scheduling latency of $30$s in a simple topology (setting [(i)]{.ltx_text .ltx_font_italic}) when scheduling $14$ nodes. In a medium topology (setting [(ii)]{.ltx_text .ltx_font_italic}) it takes $100$s+ to schedule a job with $10$ nodes. In contrast, our algorithm has a low latency even if it is required to schedule a $512$ node job in a cluster with $1000+$ nodes.
:::
::::::::::

:::::::: {#S7.SS2 .section .ltx_subsection}
### [7.2 ]{.ltx_tag .ltx_tag_subsection}Cluster Experiment {#cluster-experiment .ltx_title .ltx_title_subsection}

::: {#S7.SS2.p1 .ltx_para}
To evaluate the effectiveness of Arnold in real-world environments, we run experiments in our GPU cluster. [The specific information such as the number of GPUs and the model size, is hidden due to business concerns.]{.ltx_text .ltx_font_italic} One of our LLMs is a MoE variant and was trained previously with more than $9600$ GPUs ($1200+$ nodes). We first run the experiment by scheduling the job with $208$ GPUs, and validate the speedup achieved by Arnold. We then run the pre-training at full scale. We compare Arnold with an SOTA production system for LLMs, MegaScale [megascale](#bib.bib18){.ltx_ref} , which takes a full-stack solution to optimize LLMs training and scale to $O(10,000)$ GPUs.
:::

<figure id="S7.F9" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S7.F9.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/x11.png" id="S7.F9.sf1.g1" class="ltx_graphics ltx_img_square" width="659" height="659" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(a)</span> </span><span class="ltx_text" style="font-size:90%;">208 GPUs.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S7.F9.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/x12.png" id="S7.F9.sf2.g1" class="ltx_graphics ltx_img_square" width="660" height="658" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(b)</span> </span><span class="ltx_text" style="font-size:90%;">9600+ GPUs.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S7.F9.sf3" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2509.15940/assets/x13.png" id="S7.F9.sf3.g1" class="ltx_graphics ltx_img_landscape" width="660" height="384" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(c)</span> </span><span class="ltx_text" style="font-size:90%;">Throughput over steps.</span></figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 9</span>: </span><span class="ltx_text" style="font-size:90%;">Cluster experiments.</span></figcaption>
</figure>

::: {#S7.SS2.p2 .ltx_para}
[End-to-end experiments.]{.ltx_text .ltx_font_bold} Figure [[9(a)]{.ltx_text .ltx_ref_tag}](#S7.F9.sf1 "In Figure 9 ‣ 7.2 Cluster Experiment ‣ 7 Evaluation ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} and [[9(b)]{.ltx_text .ltx_ref_tag}](#S7.F9.sf2 "In Figure 9 ‣ 7.2 Cluster Experiment ‣ 7 Evaluation ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} illustrate the average throughput of the two systems. Arnold achieves an average speedup of $5.7\%$ and $10.6\%$ respectively. We observe that Arnold reduces the maximum spread for the DP group and the PP group by $3$x and $2$x in the medium-scale experiment, while for the full-scale experiment, the reduction is $1.5$x and $1.3$x. This is because it is more likely to spread nodes across minipods in the cluster for medium-scale experiment if not planned carefully. However, for the full-scale experiment, the requested GPUs take up more than 50% of the total number of GPUs in the cluster, and therefore the space of scheduling is shrunk.
:::

::: {#S7.SS2.p3 .ltx_para}
Nevertheless, we observe that as the model size scales and more nodes are added, the speedup achieved by Arnold also increases. The finding is consistent with Figure [[5(b)]{.ltx_text .ltx_ref_tag}](#S4.F5.sf2 "In Figure 5 ‣ 4 Characterization of Communication Patterns for LPJs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} in §[[4]{.ltx_text .ltx_ref_tag}](#S4 "4 Characterization of Communication Patterns for LPJs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}, in which the effectiveness of optimized placement is more prominent as the model size scales up. Our production deployment encompasses significantly larger models exceeding 400 billion parameters distributed across a substantially higher number of nodes, resulting in intensive network communication demands.
:::

::: {#S7.SS2.p4 .ltx_para}
Figure [[9(c)]{.ltx_text .ltx_ref_tag}](#S7.F9.sf3 "In Figure 9 ‣ 7.2 Cluster Experiment ‣ 7 Evaluation ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} plots the performance of the full-scale experiment over the training steps. Despite the performance fluctuation at the $160$-th step due to torch profiling, we observe that the Arnold outperforms MegaScale consistently. The LPJ runs for more than one month, and thus the improvement is significant for downstream tasks as well as cost savings (GPU hours and human resources). Moreover, the proposed optimization is orthogonal to those reviewed in Appendix [[A]{.ltx_text .ltx_ref_tag}](#A1 "Appendix A Related Works ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}, allowing it to be applied in conjunction with existing methods. Importantly, the optimization is fully transparent to end users.
:::

<figure id="S7.F10" class="ltx_figure">
<img src="/html/2509.15940/assets/x14.png" id="S7.F10.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="829" height="308" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 10</span>: </span><span class="ltx_text" style="font-size:90%;">Breakdown analysis.</span></figcaption>
</figure>

::: {#S7.SS2.p5 .ltx_para}
[Breakdown analysis.]{.ltx_text .ltx_font_bold} Figure [[10]{.ltx_text .ltx_ref_tag}](#S7.F10 "Figure 10 ‣ 7.2 Cluster Experiment ‣ 7 Evaluation ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} shows kernel-level examination of both systems using the Torch profiler. It reveals that communication and topology-aligned placement strategies yield a nuanced impact: while they enhance the performance of a communication kernel as expected, they simultaneously introduce performance degradation in other kernels, including a computation kernel. Through systematic ablation studies in Appendix [[I]{.ltx_text .ltx_ref_tag}](#A9 "Appendix I Break-down analysis ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}, we identify resource contention between GPU streams as the fundamental mechanism underlying this phenomenon, which presents when communication and computation kernels execute concurrently across multiple streams. These observations broaden the scope of topology-aware scheduling by showing that its impact extends beyond communication efficiency, influencing the execution characteristics of computation kernels as well.
:::
::::::::
::::::::::::::::::

:::: {#S8 .section .ltx_section}
## [8 ]{.ltx_tag .ltx_tag_section}Conclusion {#conclusion .ltx_title .ltx_title_section}

::: {#S8.p1 .ltx_para}
In this work, we present Arnold, a scheduling system that summarizes our experience in effectively scheduling LPJs at scale. In-depth characterization is performed and a scheduling algorithm is developed to align LPJs with the topology of modern data centers. Through experiments, we show the effectiveness both in simulation-based and real-world GPU clusters.
:::
::::

::: {#bib .section .ltx_bibliography}
## References {#references .ltx_title .ltx_title_bibliography}

- [[\[1\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Mohammad Al-Fares, Alexander Loukissas, and Amin Vahdat. ]{.ltx_bibblock} [A scalable, commodity data center network architecture. ]{.ltx_bibblock} [In [Proceedings of the ACM SIGCOMM 2008 Conference on Data Communication]{.ltx_text .ltx_font_italic}, SIGCOMM '08, page 63--74, New York, NY, USA, 2008. Association for Computing Machinery. ]{.ltx_bibblock}]{#bib.bib1}
- [[\[2\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Marcelo Amaral, Jordà Polo, David Carrera, Seetharami Seelam, and Malgorzata Steinder. ]{.ltx_bibblock} [Topology-aware gpu scheduling for learning workloads in cloud environments. ]{.ltx_bibblock} [In [Proceedings of the International Conference for High Performance Computing, Networking, Storage and Analysis]{.ltx_text .ltx_font_italic}, SC '17, New York, NY, USA, 2017. Association for Computing Machinery. ]{.ltx_bibblock}]{#bib.bib2}
- [[\[3\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ George Amvrosiadis, Jun Woo Park, Gregory R. Ganger, Garth A. Gibson, Elisabeth Baseman, and Nathan DeBardeleben. ]{.ltx_bibblock} [On the diversity of cluster workloads and its impact on research results. ]{.ltx_bibblock} [In [2018 USENIX Annual Technical Conference (USENIX ATC 18)]{.ltx_text .ltx_font_italic}, pages 533--546, Boston, MA, July 2018. USENIX Association. ]{.ltx_bibblock}]{#bib.bib3}
- [[\[4\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Suresh Bolusani, Mathieu Besançon, Ksenia Bestuzheva, Antonia Chmiela, João Dionísio, Tim Donkiewicz, Jasper van Doornmalen, Leon Eifler, Mohammed Ghannam, Ambros Gleixner, Christoph Graczyk, Katrin Halbig, Ivo Hedtke, Alexander Hoen, Christopher Hojny, Rolf van der Hulst, Dominik Kamp, Thorsten Koch, Kevin Kofler, Jurgen Lentz, Julian Manns, Gioni Mexi, Erik Mühmer, Marc E. Pfetsch, Franziska Schlösser, Felipe Serrano, Yuji Shinano, Mark Turner, Stefan Vigerske, Dieter Weninger, and Liding Xu. ]{.ltx_bibblock} [The scip optimization suite 9.0, 2024. ]{.ltx_bibblock}]{#bib.bib4}
- [[\[5\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jiamin Cao, Yu Guan, Kun Qian, Jiaqi Gao, Wencong Xiao, Jianbo Dong, Binzhang Fu, Dennis Cai, and Ennan Zhai. ]{.ltx_bibblock} [Crux: Gpu-efficient communication scheduling for deep learning training. ]{.ltx_bibblock} [In [Proceedings of the ACM SIGCOMM 2024 Conference]{.ltx_text .ltx_font_italic}, ACM SIGCOMM '24, page 1--15, New York, NY, USA, 2024. Association for Computing Machinery. ]{.ltx_bibblock}]{#bib.bib5}
- [[\[6\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Li-Wen Chang, Wenlei Bao, Qi Hou, Chengquan Jiang, Ningxin Zheng, Yinmin Zhong, Xuanrun Zhang, Zuquan Song, Chengji Yao, Ziheng Jiang, Haibin Lin, Xin Jin, and Xin Liu. ]{.ltx_bibblock} [Flux: Fast software-based communication overlap on gpus through kernel fusion, 2024. ]{.ltx_bibblock}]{#bib.bib6}
- [[\[7\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Arnab Choudhury, Yang Wang, Tuomas Pelkonen, Kutta Srinivasan, Abha Jain, Shenghao Lin, Delia David, Siavash Soleimanifard, Michael Chen, Abhishek Yadav, Ritesh Tijoriwala, Denis Samoylov, and Chunqiang Tang. ]{.ltx_bibblock} [MAST: Global scheduling of ML training across Geo-Distributed datacenters at hyperscale. ]{.ltx_bibblock} [In [18th USENIX Symposium on Operating Systems Design and Implementation (OSDI 24)]{.ltx_text .ltx_font_italic}, pages 563--580, Santa Clara, CA, July 2024. USENIX Association. ]{.ltx_bibblock}]{#bib.bib7}
- [[\[8\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Charles Clos. ]{.ltx_bibblock} [A study of non-blocking switching networks. ]{.ltx_bibblock} [[The Bell System Technical Journal]{.ltx_text .ltx_font_italic}, 32(2):406--424, 1953. ]{.ltx_bibblock}]{#bib.bib8}
- [[\[9\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tri Dao, Daniel Y. Fu, Stefano Ermon, Atri Rudra, and Christopher Ré. ]{.ltx_bibblock} [Flashattention: Fast and memory-efficient exact attention with io-awareness, 2022. ]{.ltx_bibblock}]{#bib.bib9}
- [[\[10\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ F. Ercal, J. Ramanujam, and P. Sadayappan. ]{.ltx_bibblock} [Task allocation onto a hypercube by recursive mincut bipartitioning. ]{.ltx_bibblock} [In [Proceedings of the Third Conference on Hypercube Concurrent Computers and Applications: Architecture, Software, Computer Systems, and General Issues - Volume 1]{.ltx_text .ltx_font_italic}, C3P, page 210--221, New York, NY, USA, 1988. Association for Computing Machinery. ]{.ltx_bibblock}]{#bib.bib10}
- [[\[11\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ C.M. Fiduccia and R.M. Mattheyses. ]{.ltx_bibblock} [A linear-time heuristic for improving network partitions. ]{.ltx_bibblock} [In [19th Design Automation Conference]{.ltx_text .ltx_font_italic}, pages 175--181, 1982. ]{.ltx_bibblock}]{#bib.bib11}
- [[\[12\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Guoliang He and Eiko Yoneki. ]{.ltx_bibblock} [Cuasmrl: Optimizing gpu sass schedules via deep reinforcement learning, 2025. ]{.ltx_bibblock}]{#bib.bib12}
- [[\[13\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Qinghao Hu, Zhisheng Ye, Zerui Wang, Guoteng Wang, Meng Zhang, Qiaoling Chen, Peng Sun, Dahua Lin, Xiaolin Wang, Yingwei Luo, Yonggang Wen, and Tianwei Zhang. ]{.ltx_bibblock} [Characterization of large language model development in the datacenter. ]{.ltx_bibblock} [In [NSDI]{.ltx_text .ltx_font_italic}, pages 709--729, 2024. ]{.ltx_bibblock}]{#bib.bib13}
- [[\[14\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ imbue team. ]{.ltx_bibblock} [From bare metal to a 70b model: infrastructure set-up and scripts. ]{.ltx_bibblock} [[https://imbue.com/research/70b-infrastructure/](https://imbue.com/research/70b-infrastructure/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock} [Accessed: 2024-12-01. ]{.ltx_bibblock}]{#bib.bib14}
- [[\[15\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Youhe Jiang, Fangcheng Fu, Xupeng Miao, Xiaonan Nie, and Bin Cui. ]{.ltx_bibblock} [Osdp: Optimal sharded data parallel for distributed deep learning. ]{.ltx_bibblock} [In [Proceedings of the Thirty-Second International Joint Conference on Artificial Intelligence]{.ltx_text .ltx_font_italic}, pages 2142--2150, 2023. ]{.ltx_bibblock}]{#bib.bib15}
- [[\[16\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Youhe Jiang, Huaxi Gu, Yunfeng Lu, and Xiaoshan Yu. ]{.ltx_bibblock} [2d-hra: Two-dimensional hierarchical ring-based all-reduce algorithm in large-scale distributed machine learning. ]{.ltx_bibblock} [[IEEE Access]{.ltx_text .ltx_font_italic}, 8:183488--183494, 2020. ]{.ltx_bibblock}]{#bib.bib16}
- [[\[17\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zhifeng Jiang, Zhihua Jin, and Guoliang He. ]{.ltx_bibblock} [Safeguarding system prompts for llms. ]{.ltx_bibblock} [[arXiv preprint arXiv:2412.13426]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib17}
- [[\[18\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ziheng Jiang, Haibin Lin, Yinmin Zhong, Qi Huang, Yangrui Chen, Zhi Zhang, Yanghua Peng, Xiang Li, Cong Xie, Shibiao Nong, Yulu Jia, Sun He, Hongmin Chen, Zhihao Bai, Qi Hou, Shipeng Yan, Ding Zhou, Yiyao Sheng, Zhuo Jiang, Haohan Xu, Haoran Wei, Zhang Zhang, Pengfei Nie, Leqi Zou, Sida Zhao, Liang Xiang, Zherui Liu, Zhe Li, Xiaoying Jia, Jianxi Ye, Xin Jin, and Xin Liu. ]{.ltx_bibblock} [MegaScale: Scaling large language model training to more than 10,000 GPUs. ]{.ltx_bibblock} [In [21st USENIX Symposium on Networked Systems Design and Implementation (NSDI 24)]{.ltx_text .ltx_font_italic}, pages 745--760, Santa Clara, CA, April 2024. USENIX Association. ]{.ltx_bibblock}]{#bib.bib18}
- [[\[19\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jared Kaplan, Sam McCandlish, Tom Henighan, Tom B. Brown, Benjamin Chess, Rewon Child, Scott Gray, Alec Radford, Jeffrey Wu, and Dario Amodei. ]{.ltx_bibblock} [Scaling laws for neural language models, 2020. ]{.ltx_bibblock}]{#bib.bib19}
- [[\[20\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Guolin Ke, Qi Meng, Thomas Finley, Taifeng Wang, Wei Chen, Weidong Ma, Qiwei Ye, and Tie-Yan Liu. ]{.ltx_bibblock} [Lightgbm: A highly efficient gradient boosting decision tree. ]{.ltx_bibblock} [In I. Guyon, U. Von Luxburg, S. Bengio, H. Wallach, R. Fergus, S. Vishwanathan, and R. Garnett, editors, [Advances in Neural Information Processing Systems]{.ltx_text .ltx_font_italic}, volume 30. Curran Associates, Inc., 2017. ]{.ltx_bibblock}]{#bib.bib20}
- [[\[21\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Kubernetes Authors. ]{.ltx_bibblock} [Kubernetes: Production-grade container orchestration. ]{.ltx_bibblock} [[https://kubernetes.io](https://kubernetes.io){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}, 2023. ]{.ltx_bibblock} [Accessed: 2024-12-01. ]{.ltx_bibblock}]{#bib.bib21}
- [[\[22\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xupeng Miao, Yujie Wang, Youhe Jiang, Chunan Shi, Xiaonan Nie, Hailin Zhang, and Bin Cui. ]{.ltx_bibblock} [Galvatron: Efficient transformer training over multiple gpus using automatic parallelism. ]{.ltx_bibblock} [[Proceedings of the VLDB Endowment]{.ltx_text .ltx_font_italic}, 16(3):470--479, 2022. ]{.ltx_bibblock}]{#bib.bib22}
- [[\[23\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xupeng Miao, Yujie Wang, Youhe Jiang, Chunan Shi, Xiaonan Nie, Hailin Zhang, and Bin Cui. ]{.ltx_bibblock} [Galvatron: Efficient transformer training over multiple gpus using automatic parallelism. ]{.ltx_bibblock} [[Proc. VLDB Endow.]{.ltx_text .ltx_font_italic}, 16(3):470--479, November 2022. ]{.ltx_bibblock}]{#bib.bib23}
- [[\[24\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Deepak Narayanan, Mohammad Shoeybi, Jared Casper, Patrick LeGresley, Mostofa Patwary, Vijay Korthikanti, Dmitri Vainbrand, Prethvi Kashinkunti, Julie Bernauer, Bryan Catanzaro, Amar Phanishayee, and Matei Zaharia. ]{.ltx_bibblock} [Efficient large-scale language model training on gpu clusters using megatron-lm. ]{.ltx_bibblock} [In [SC21: International Conference for High Performance Computing, Networking, Storage and Analysis]{.ltx_text .ltx_font_italic}, pages 1--14, 2021. ]{.ltx_bibblock}]{#bib.bib24}
- [[\[25\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Andrew Newell, Dimitrios Skarlatos, Jingyuan Fan, Pavan Kumar, Maxim Khutornenko, Mayank Pundir, Yirui Zhang, Mingjun Zhang, Yuanlai Liu, Linh Le, Brendon Daugherty, Apurva Samudra, Prashasti Baid, James Kneeland, Igor Kabiljo, Dmitry Shchukin, Andre Rodrigues, Scott Michelson, Ben Christensen, Kaushik Veeraraghavan, and Chunqiang Tang. ]{.ltx_bibblock} [Ras: Continuously optimized region-wide datacenter resource allocation. ]{.ltx_bibblock} [In [Proceedings of the ACM SIGOPS 28th Symposium on Operating Systems Principles]{.ltx_text .ltx_font_italic}, SOSP '21, page 505--520, New York, NY, USA, 2021. Association for Computing Machinery. ]{.ltx_bibblock}]{#bib.bib25}
- [[\[26\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ NVIDIA. ]{.ltx_bibblock} [Gb200-nvl72. ]{.ltx_bibblock} [[https://www.nvidia.com/en-us/data-center/gb200-nvl72/](https://www.nvidia.com/en-us/data-center/gb200-nvl72/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock} [Accessed: 2024-12-01. ]{.ltx_bibblock}]{#bib.bib26}
- [[\[27\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ NVIDIA. ]{.ltx_bibblock} [Infiniband networking solutions. ]{.ltx_bibblock} [[https://www.nvidia.com/en-us/networking/products/infiniband/](https://www.nvidia.com/en-us/networking/products/infiniband/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock} [Accessed: 2024-12-01. ]{.ltx_bibblock}]{#bib.bib27}
- [[\[28\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ NVIDIA. ]{.ltx_bibblock} [Nccl. ]{.ltx_bibblock} [[https://developer.nvidia.com/nccl](https://developer.nvidia.com/nccl){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock} [Accessed: 2024-12-01. ]{.ltx_bibblock}]{#bib.bib28}
- [[\[29\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ NVIDIA. ]{.ltx_bibblock} [Nccl-test. ]{.ltx_bibblock} [[https://github.com/NVIDIA/nccl-tests](https://github.com/NVIDIA/nccl-tests){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock} [Accessed: 2024-12-01. ]{.ltx_bibblock}]{#bib.bib29}
- [[\[30\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ NVIDIA. ]{.ltx_bibblock} [Nvlink. ]{.ltx_bibblock} [[https://www.nvidia.com/en-us/data-center/nvlink/](https://www.nvidia.com/en-us/data-center/nvlink/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock} [Accessed: 2024-12-01. ]{.ltx_bibblock}]{#bib.bib30}
- [[\[31\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ OpenAI. ]{.ltx_bibblock} [Chatgpt. ]{.ltx_bibblock} [[https://openai.com/index/chatgpt/](https://openai.com/index/chatgpt/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock} [Accessed: 2024-12-01. ]{.ltx_bibblock}]{#bib.bib31}
- [[\[32\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Rina Panigrahy, Vijayan Prabhakaran, Kunal Talwar, Udi Wieder, and Rama Ramasubramanian. ]{.ltx_bibblock} [Validating heuristics for virtual machines consolidation. ]{.ltx_bibblock} [Technical Report MSR-TR-2011-9, January 2011. ]{.ltx_bibblock}]{#bib.bib32}
- [[\[33\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Kun Qian, Yongqing Xi, Jiamin Cao, Jiaqi Gao, Yichi Xu, Yu Guan, Binzhang Fu, Xuemei Shi, Fangbo Zhu, Rui Miao, Chao Wang, Peng Wang, Pengcheng Zhang, Xianlong Zeng, Eddie Ruan, Zhiping Yao, Ennan Zhai, and Dennis Cai. ]{.ltx_bibblock} [Alibaba hpn: A data center network for large language model training. ]{.ltx_bibblock} [In [Proceedings of the ACM SIGCOMM 2024 Conference]{.ltx_text .ltx_font_italic}, ACM SIGCOMM '24, page 691--706, New York, NY, USA, 2024. Association for Computing Machinery. ]{.ltx_bibblock}]{#bib.bib33}
- [[\[34\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Samyam Rajbhandari, Jeff Rasley, Olatunji Ruwase, and Yuxiong He. ]{.ltx_bibblock} [Zero: Memory optimizations toward training trillion parameter models, 2020. ]{.ltx_bibblock}]{#bib.bib34}
- [[\[35\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jeff Rasley, Samyam Rajbhandari, Olatunji Ruwase, and Yuxiong He. ]{.ltx_bibblock} [Deepspeed: System optimizations enable training deep learning models with over 100 billion parameters. ]{.ltx_bibblock} [In [Proceedings of the 26th ACM SIGKDD International Conference on Knowledge Discovery & Data Mining]{.ltx_text .ltx_font_italic}, KDD '20, page 3505--3506, New York, NY, USA, 2020. Association for Computing Machinery. ]{.ltx_bibblock}]{#bib.bib35}
- [[\[36\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Noam Shazeer, \*Azalia Mirhoseini, \*Krzysztof Maziarz, Andy Davis, Quoc Le, Geoffrey Hinton, and Jeff Dean. ]{.ltx_bibblock} [Outrageously large neural networks: The sparsely-gated mixture-of-experts layer. ]{.ltx_bibblock} [In [International Conference on Learning Representations]{.ltx_text .ltx_font_italic}, 2017. ]{.ltx_bibblock}]{#bib.bib36}
- [[\[37\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Mohammad Shoeybi, Mostofa Patwary, Raul Puri, Patrick LeGresley, Jared Casper, and Bryan Catanzaro. ]{.ltx_bibblock} [Megatron-lm: Training multi-billion parameter language models using model parallelism, 2020. ]{.ltx_bibblock}]{#bib.bib37}
- [[\[38\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Hugo Touvron, Louis Martin, Kevin Stone, Peter Albert, Amjad Almahairi, Yasmine Babaei, Nikolay Bashlykov, Soumya Batra, Prajjwal Bhargava, Shruti Bhosale, Dan Bikel, Lukas Blecher, Cristian Canton Ferrer, Moya Chen, Guillem Cucurull, David Esiobu, Jude Fernandes, Jeremy Fu, Wenyin Fu, Brian Fuller, Cynthia Gao, Vedanuj Goswami, Naman Goyal, Anthony Hartshorn, Saghar Hosseini, Rui Hou, Hakan Inan, Marcin Kardas, Viktor Kerkez, Madian Khabsa, Isabel Kloumann, Artem Korenev, Punit Singh Koura, Marie-Anne Lachaux, Thibaut Lavril, Jenya Lee, Diana Liskovich, Yinghai Lu, Yuning Mao, Xavier Martinet, Todor Mihaylov, Pushkar Mishra, Igor Molybog, Yixin Nie, Andrew Poulton, Jeremy Reizenstein, Rashi Rungta, Kalyan Saladi, Alan Schelten, Ruan Silva, Eric Michael Smith, Ranjan Subramanian, Xiaoqing Ellen Tan, Binh Tang, Ross Taylor, Adina Williams, Jian Xiang Kuan, Puxin Xu, Zheng Yan, Iliyan Zarov, Yuchen Zhang, Angela Fan, Melanie Kambadur, Sharan Narang, Aurelien Rodriguez, Robert Stojnic, Sergey Edunov, and Thomas Scialom. ]{.ltx_bibblock} [Llama 2: Open foundation and fine-tuned chat models, 2023. ]{.ltx_bibblock}]{#bib.bib38}
- [[\[39\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Colin Unger, Zhihao Jia, Wei Wu, Sina Lin, Mandeep Baines, Carlos Efrain Quintero Narvaez, Vinay Ramakrishnaiah, Nirmal Prajapati, Pat McCormick, Jamaludin Mohd-Yusof, Xi Luo, Dheevatsa Mudigere, Jongsoo Park, Misha Smelyanskiy, and Alex Aiken. ]{.ltx_bibblock} [Unity: Accelerating DNN training through joint optimization of algebraic transformations and parallelization. ]{.ltx_bibblock} [In [16th USENIX Symposium on Operating Systems Design and Implementation (OSDI 22)]{.ltx_text .ltx_font_italic}, pages 267--284, Carlsbad, CA, July 2022. USENIX Association. ]{.ltx_bibblock}]{#bib.bib39}
- [[\[40\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ashish Vaswani, Noam Shazeer, Niki Parmar, Jakob Uszkoreit, Llion Jones, Aidan N Gomez, Ł ukasz Kaiser, and Illia Polosukhin. ]{.ltx_bibblock} [Attention is all you need. ]{.ltx_bibblock} [In I. Guyon, U. Von Luxburg, S. Bengio, H. Wallach, R. Fergus, S. Vishwanathan, and R. Garnett, editors, [Advances in Neural Information Processing Systems]{.ltx_text .ltx_font_italic}, volume 30. Curran Associates, Inc., 2017. ]{.ltx_bibblock}]{#bib.bib40}
- [[\[41\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Weiyang Wang, Manya Ghobadi, Kayvon Shakeri, Ying Zhang, and Naader Hasani. ]{.ltx_bibblock} [Rail-only: A low-cost high-performance network for training llms with trillion parameters, 2024. ]{.ltx_bibblock}]{#bib.bib41}
- [[\[42\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yujie Wang, Youhe Jiang, Xupeng Miao, Fangcheng Fu, Shenhan Zhu, Xiaonan Nie, Yaofeng Tu, and Bin Cui. ]{.ltx_bibblock} [Improving automatic parallel training via balanced memory workload optimization. ]{.ltx_bibblock} [[IEEE Transactions on Knowledge and Data Engineering]{.ltx_text .ltx_font_italic}, 36(8):3906--3920, 2024. ]{.ltx_bibblock}]{#bib.bib42}
- [[\[43\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Qizhen Weng, Wencong Xiao, Yinghao Yu, Wei Wang, Cheng Wang, Jian He, Yong Li, Liping Zhang, Wei Lin, and Yu Ding. ]{.ltx_bibblock} [MLaaS in the wild: Workload analysis and scheduling in Large-Scale heterogeneous GPU clusters. ]{.ltx_bibblock} [In [19th USENIX Symposium on Networked Systems Design and Implementation (NSDI 22)]{.ltx_text .ltx_font_italic}, pages 945--960, Renton, WA, April 2022. USENIX Association. ]{.ltx_bibblock}]{#bib.bib43}
- [[\[44\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Qizhen Weng, Lingyun Yang, Yinghao Yu, Wei Wang, Xiaochuan Tang, Guodong Yang, and Liping Zhang. ]{.ltx_bibblock} [Beware of fragmentation: Scheduling GPU-Sharing workloads with fragmentation gradient descent. ]{.ltx_bibblock} [In [2023 USENIX Annual Technical Conference (USENIX ATC 23)]{.ltx_text .ltx_font_italic}, pages 995--1008, Boston, MA, July 2023. USENIX Association. ]{.ltx_bibblock}]{#bib.bib44}
- [[\[45\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Wencong Xiao, Romil Bhardwaj, Ramachandran Ramjee, Muthian Sivathanu, Nipun Kwatra, Zhenhua Han, Pratyush Patel, Xuan Peng, Hanyu Zhao, Quanlu Zhang, Fan Yang, and Lidong Zhou. ]{.ltx_bibblock} [Gandiva: Introspective cluster scheduling for deep learning. ]{.ltx_bibblock} [In [13th USENIX Symposium on Operating Systems Design and Implementation (OSDI 18)]{.ltx_text .ltx_font_italic}, pages 595--610, Carlsbad, CA, October 2018. USENIX Association. ]{.ltx_bibblock}]{#bib.bib45}
- [[\[46\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Wencong Xiao, Shiru Ren, Yong Li, Yang Zhang, Pengyang Hou, Zhi Li, Yihui Feng, Wei Lin, and Yangqing Jia. ]{.ltx_bibblock} [AntMan: Dynamic scaling on GPU clusters for deep learning. ]{.ltx_bibblock} [In [14th USENIX Symposium on Operating Systems Design and Implementation (OSDI 20)]{.ltx_text .ltx_font_italic}, pages 533--548. USENIX Association, November 2020. ]{.ltx_bibblock}]{#bib.bib46}
- [[\[47\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Dian Xiong, Li Chen, Youhe Jiang, Dan Li, Shuai Wang, and Songtao Wang. ]{.ltx_bibblock} [Revisiting the time cost model of allreduce. ]{.ltx_bibblock} [[arXiv preprint arXiv:2409.04202]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib47}
- [[\[48\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ran Yan, Youhe Jiang, Xiaonan Nie, Fangcheng Fu, Bin Cui, and Binhang Yuan. ]{.ltx_bibblock} [Hexiscale: Accommodating large language model training over heterogeneous environment. ]{.ltx_bibblock} [[arXiv preprint arXiv:2409.01143]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib48}
- [[\[49\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ran Yan, Youhe Jiang, and Binhang Yuan. ]{.ltx_bibblock} [Flash sparse attention: An alternative efficient implementation of native sparse attention kernel. ]{.ltx_bibblock} [[arXiv preprint arXiv:2508.18224]{.ltx_text .ltx_font_italic}, 2025. ]{.ltx_bibblock}]{#bib.bib49}
- [[\[50\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ An Yang, Baosong Yang, Binyuan Hui, Bo Zheng, Bowen Yu, Chang Zhou, Chengpeng Li, Chengyuan Li, Dayiheng Liu, Fei Huang, Guanting Dong, Haoran Wei, Huan Lin, Jialong Tang, Jialin Wang, Jian Yang, Jianhong Tu, Jianwei Zhang, Jianxin Ma, Jianxin Yang, Jin Xu, Jingren Zhou, Jinze Bai, Jinzheng He, Junyang Lin, Kai Dang, Keming Lu, Keqin Chen, Kexin Yang, Mei Li, Mingfeng Xue, Na Ni, Pei Zhang, Peng Wang, Ru Peng, Rui Men, Ruize Gao, Runji Lin, Shijie Wang, Shuai Bai, Sinan Tan, Tianhang Zhu, Tianhao Li, Tianyu Liu, Wenbin Ge, Xiaodong Deng, Xiaohuan Zhou, Xingzhang Ren, Xinyu Zhang, Xipin Wei, Xuancheng Ren, Xuejing Liu, Yang Fan, Yang Yao, Yichang Zhang, Yu Wan, Yunfei Chu, Yuqiong Liu, Zeyu Cui, Zhenru Zhang, Zhifang Guo, and Zhihao Fan. ]{.ltx_bibblock} [Qwen2 technical report, 2024. ]{.ltx_bibblock}]{#bib.bib50}
- [[\[51\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Li Zhang, Youhe Jiang, Guoliang He, Xin Chen, Han Lv, Qian Yao, Fangcheng Fu, and Kai Chen. ]{.ltx_bibblock} [Efficient mixed-precision large language model inference with turbomind. ]{.ltx_bibblock} [[arXiv preprint arXiv:2508.15601]{.ltx_text .ltx_font_italic}, 2025. ]{.ltx_bibblock}]{#bib.bib51}
- [[\[52\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Lianmin Zheng, Zhuohan Li, Hao Zhang, Yonghao Zhuang, Zhifeng Chen, Yanping Huang, Yida Wang, Yuanzhong Xu, Danyang Zhuo, Eric P. Xing, Joseph E. Gonzalez, and Ion Stoica. ]{.ltx_bibblock} [Alpa: Automating inter- and Intra-Operator parallelism for distributed deep learning. ]{.ltx_bibblock} [In [16th USENIX Symposium on Operating Systems Design and Implementation (OSDI 22)]{.ltx_text .ltx_font_italic}, pages 559--578, Carlsbad, CA, July 2022. USENIX Association. ]{.ltx_bibblock}]{#bib.bib52}
:::

::: {.ltx_pagination .ltx_role_newpage}
:::

::::: {#A1 .section .ltx_appendix}
## [Appendix A ]{.ltx_tag .ltx_tag_appendix}Related Works {#appendix-a-related-works .ltx_title .ltx_title_appendix}

::: {#A1.p1 .ltx_para}
[LLMs training.]{.ltx_text .ltx_font_bold} LLMs have become a significant workload in the field of machine learning \[[31](#bib.bib31){.ltx_ref}, [38](#bib.bib38){.ltx_ref}, [50](#bib.bib50){.ltx_ref}\]. The first step of developing an LLM needs to train a large transformer model with trillions of tokens, so the computing infrastructure continues to evolve to adapt to the challenging workload. For example, efficient parallelization strategies are searched by model parallelizers \[[52](#bib.bib52){.ltx_ref}, [23](#bib.bib23){.ltx_ref}, [39](#bib.bib39){.ltx_ref}\], training frameworks specialized for training scalability are built to orchestrate large-scale worker nodes \[[18](#bib.bib18){.ltx_ref}, [37](#bib.bib37){.ltx_ref}, [24](#bib.bib24){.ltx_ref}\], and high-performance operators are developed to maximize the utilization of hardware accelerators \[[9](#bib.bib9){.ltx_ref}, [6](#bib.bib6){.ltx_ref}, [12](#bib.bib12){.ltx_ref}\]. However, those works are orthogonal to the optimization proposed in this paper, as the physical network topology is only visible at the cluster scheduling layer. Therefore, our work is transparent to the underlying infrastructure codes and the speedup is achieved on top of existing effort to accelerate the training performance.
:::

::: {#A1.p2 .ltx_para}
[Deep learning job schedulers.]{.ltx_text .ltx_font_bold} Job scheduling systems for deep learning tasks have been widely deployed by companies \[[7](#bib.bib7){.ltx_ref}, [45](#bib.bib45){.ltx_ref}, [44](#bib.bib44){.ltx_ref}, [2](#bib.bib2){.ltx_ref}, [25](#bib.bib25){.ltx_ref}\]. However, none of deep learning jobs have come even close to the scale and importance of LLM pre-training jobs. Arnold addresses this gap by providing a solution tailored to scheduling LLM pre-training jobs, complementing existing schedulers. Recent studies have begun to explore the characteristics of LLM workloads in GPU clusters \[[13](#bib.bib13){.ltx_ref}\]. In contrast, our work specifically targets the optimization of LLM pre-training performance.
:::
:::::

:::::: {#A2 .section .ltx_appendix}
## [Appendix B ]{.ltx_tag .ltx_tag_appendix}Limitation and Future Works {#appendix-b-limitation-and-future-works .ltx_title .ltx_title_appendix}

::: {#A2.p1 .ltx_para}
[Failure recovery.]{.ltx_text .ltx_font_bold} On hardware failure, the optimal placement of LPJs will inevitably change, but it is too expensive to solve the MIP and then migrate to the new placement. A potential approach is to increase the number of GPUs of communication groups in the initial scheduling as backups, which only run preemptive jobs and can be replaced with failure nodes when needed.
:::

::: {#A2.p2 .ltx_para}
[Other network topology.]{.ltx_text .ltx_font_bold} The characterization results, while derived from our in-house data center environments, exhibit broad applicability to CLOS-based network topologies, which represent the predominant network architecture in modern data center deployments. By varying the affinity parameters, one can effectively trade-off the balance of DP and PP groups in their own data centers. Therefore, the characterization methodology and the scheduling algorithm are generalizable to other data centers for large-scale LPJs.
:::

::: {#A2.p3 .ltx_para}
On the other hand, newer network architectures \[[33](#bib.bib33){.ltx_ref}, [26](#bib.bib26){.ltx_ref}, [41](#bib.bib41){.ltx_ref}\] dedicated to LLM training are emerging, which necessitates careful consideration of scheduling algorithms. Our proposed scheduling approach takes a pioneering step by aligning communication patterns with data center topology for LPJs, and its effectiveness has been evaluated in real-world production cluster. We leave further exploration of this direction to future work.
:::
::::::

:::::: {#A3 .section .ltx_appendix}
## [Appendix C ]{.ltx_tag .ltx_tag_appendix}Analytical Estimation for Communication Volume {#appendix-c-analytical-estimation-for-communication-volume .ltx_title .ltx_title_appendix}

::: {#A3.p1 .ltx_para}
To estimate the communication volume of pre-training jobs, we adopt an analytical model for GPT-based variants. We use the same notation from previous work \[[24](#bib.bib24){.ltx_ref}\] by denoting the vocabulary size $V$, global batch size $gb$, micro-batch size $mb$, sequence length $s$, hidden dimension $h$, the number of layers $l$, the DP size $dp$, the PP size $pp$, the number of VPP size $vp$, number of microbatches $m$. We have:
:::

::: {#A3.p2 .ltx_para}
  -- ------------------------ -- -----------------------------------------------------
     $$m=\frac{gb}{mb*dp}$$      [(11)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- ------------------------ -- -----------------------------------------------------
:::

::: {#A3.p3 .ltx_para}
- [[•]{.ltx_tag .ltx_tag_item}]{#A3.I1.i1}

  ::: {#A3.I1.i1.p1 .ltx_para}
  DP groups. GPUs within the same group replicate the model weights and exchange parameters as well as gradients, so the communication volume can be computed using Equation [[12]{.ltx_text .ltx_ref_tag}](#A3.E12 "In 1st item ‣ Appendix C Analytical Estimation for Communication Volume ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}.
  :::

  ::: {#A3.I1.i1.p2 .ltx_para}
    -- ------------------------------------------------------------------------------------ -- -----------------------------------------------------
       $$DP-volume=h*(V+s)+l/pp*(4h^{2}+2h+\underbrace{8h^{2}+7h}_{\text{dense layer}})$$      [(12)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
    -- ------------------------------------------------------------------------------------ -- -----------------------------------------------------
  :::

  ::: {#A3.I1.i1.p3 .ltx_para}
  For MoE models, we can replace the number of parameters with MoE layers accordingly.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#A3.I1.i2}

  ::: {#A3.I1.i2.p1 .ltx_para}
  PP groups. GPUs exchange intermediate activation to adjacent PP stages, and thus we apply Equation [[13]{.ltx_text .ltx_ref_tag}](#A3.E13 "In 2nd item ‣ Appendix C Analytical Estimation for Communication Volume ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} to estimate the data volumes.
  :::

  ::: {#A3.I1.i2.p2 .ltx_para}
    -- ------------------------ -- -----------------------------------------------------
       $$PP-volume=2*mb*s*h$$      [(13)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
    -- ------------------------ -- -----------------------------------------------------
  :::
:::
::::::

::::: {#A4 .section .ltx_appendix}
## [Appendix D ]{.ltx_tag .ltx_tag_appendix}Sensitivity to Shared Load {#appendix-d-sensitivity-to-shared-load .ltx_title .ltx_title_appendix}

<figure id="A4.F11" class="ltx_figure">
<img src="/html/2509.15940/assets/x15.png" id="A4.F11.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="582" height="234" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 11</span>: </span><span class="ltx_text" style="font-size:90%;">Bandwidth interference.</span></figcaption>
</figure>

::: {#A4.p1 .ltx_para}
Since GPU clusters are usually multi-tenant to improve resource utilization, we also study the interference between inter-node communication quantitatively. Before we bring our cluster online, we perform large-scale stress test on our cluster by running NCCL tests. We record the time series of bus bandwidths and show in Figure [[11]{.ltx_text .ltx_ref_tag}](#A4.F11 "Figure 11 ‣ Appendix D Sensitivity to Shared Load ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}. The stress test consists of $3$ jobs, requesting thousands of GPUs each, and spanning across multiple minipods in the cluster. Each job performs all-to-all communication with a message size of 2GB constantly, which simulates jobs performing extensive inter-node communication on a busy cluster because all-to-all generates large amount of flows in the network.
:::

::: {#A4.p2 .ltx_para}
We can observe performance fluctuation for all three jobs. For example, after job 1 is launched at $01:22$, job 3 has a slight performance degradation of $0.5$GB/s ($3\%$). The maximum performance degradation is up to $5\%$ for job 3 during the stress test period. This suggests that jobs spanning a larger number of minipods not only suffer from increased bandwidth loss but are also more exposed to interference from other workloads in the cluster.
:::
:::::

::::: {#A5 .section .ltx_appendix}
## [Appendix E ]{.ltx_tag .ltx_tag_appendix}Ada Lovelace GPUs {#appendix-e-ada-lovelace-gpus .ltx_title .ltx_title_appendix}

::: {#A5.p1 .ltx_para}
We repeat the characterization experiment in another GPU cluster, where each node is equipped with L20 GPUs, to ensure our finding is not limited to H800 GPUs, and we briefly summarize the results in Table [[2]{.ltx_text .ltx_ref_tag}](#A5.T2 "Table 2 ‣ Appendix E Ada Lovelace GPUs ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref}.
:::

<figure id="A5.T2" class="ltx_table">
<table class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<thead class="ltx_thead">
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_l ltx_border_r ltx_border_t">Best placement</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_r ltx_border_t">Model size</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_r ltx_border_t">speedup</th>
</tr>
</thead>
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_l ltx_border_r ltx_border_t">DP-aligned</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">dense 7b</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">1.4%</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_b ltx_border_l ltx_border_r ltx_border_t">PP-aligned</td>
<td class="ltx_td ltx_align_center ltx_border_b ltx_border_r ltx_border_t">dense 14b</td>
<td class="ltx_td ltx_align_center ltx_border_b ltx_border_r ltx_border_t">0.5%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 2</span>: </span><span class="ltx_text" style="font-size:90%;">Results on L20 GPUs cluster.</span></figcaption>
</figure>

::: {#A5.p2 .ltx_para}
We observe that DP-aligned can yield greater speedups for certain model configurations. This is likely because during training, L20 GPU uses a $8$-bit data format, which halves the communication volumes between PP stages. However, DP groups still use $32$-bit for parameter and gradient synchronization, so the communication volumes remain unchanged. As a result, DP group communication can become the dominant overhead, making a placement strategy that prioritizes DP groups more beneficial. However, as the size of the model grows, the bottleneck shifts back to the communication of the PP group, so the PP-aligned is preferential.
:::
:::::

:::: {#A6 .section .ltx_appendix}
## [Appendix F ]{.ltx_tag .ltx_tag_appendix}Communication Matrix {#appendix-f-communication-matrix .ltx_title .ltx_title_appendix}

<figure id="A6.F12" class="ltx_figure">
<img src="/html/2509.15940/assets/x16.png" id="A6.F12.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="415" height="176" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 12</span>: </span><span class="ltx_text" style="font-size:90%;">Example placement (96 GPUs and 12 nodes) of a LPJ.</span></figcaption>
</figure>

::: {#A6.p1 .ltx_para}
The example placement in Figure [[12]{.ltx_text .ltx_ref_tag}](#A6.F12 "Figure 12 ‣ Appendix F Communication Matrix ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} has a DP group size of $6$, and PP group size of $2$. Nodes in different colors represent they are placed in different minipods. The number is the rank of the node. In this example, DP group is aligned and the communication of PP group must cross spine switches.
:::
::::

::::::: {#A7 .section .ltx_appendix}
## [Appendix G ]{.ltx_tag .ltx_tag_appendix}Queue Management {#appendix-g-queue-management .ltx_title .ltx_title_appendix}

::: {#A7.p1 .ltx_para}
Algorithm [[1]{.ltx_text .ltx_ref_tag}](#alg1 "Algorithm 1 ‣ Appendix G Queue Management ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} illustrates our queuing policy. It manages the job queue and reserve resources for the imminent LPJ. It also employs an ML-driven job completion time (JCT) predictor to balance the trade-off of queuing delay and resource utilization
:::

<figure id="alg1" class="ltx_float ltx_float_algorithm ltx_framed ltx_framed_top">
<div class="ltx_listing">
<div id="alg1.l1" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">1:</span></span><span class="math inline"><em>J</em></span>; //job’s configurations and metadata
</div>
<div id="alg1.l2" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">2:</span></span><span class="math inline"><em>V</em></span>; //physical view of the cluster
</div>
<div id="alg1.l3" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">3:</span></span><span class="math inline"><em>Q</em></span>; //job queue, sorted by priority and arrival time
</div>
<div id="alg1.l4" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">4:</span></span><span class="math inline"><em>O</em></span>; //delay list
</div>
<div id="alg1.l5" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">5:</span></span><span class="ltx_text ltx_font_bold">function</span> <span class="ltx_text ltx_font_smallcaps">scheduler</span>(<span class="math inline"><em>J</em></span>, <span class="math inline"><em>V</em></span>, <span class="math inline"><em>Q</em></span>)
</div>
<div id="alg1.l6" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">6:</span></span>  <span class="ltx_text ltx_font_bold">while</span> <span class="math inline"><em>T</em><em>r</em><em>u</em><em>e</em></span> <span class="ltx_text ltx_font_bold">do</span>
</div>
<div id="alg1.l7" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">7:</span></span>   <span class="math inline"><em>O</em> ← ∅</span>
</div>
<div id="alg1.l8" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">8:</span></span>   <span class="ltx_text ltx_font_bold">while</span> <span class="math inline"><em>Q</em> ≠ ∅</span> <span class="ltx_text ltx_font_bold">do</span>
</div>
<div id="alg1.l9" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">9:</span></span>     <span class="math inline"><em>J</em> ← <em>Q</em>.<em>p</em><em>o</em><em>p</em>()</span>
</div>
<div id="alg1.l10" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">10:</span></span>     <span class="ltx_text ltx_font_bold">if</span> <span class="math inline"><em>p</em><em>r</em><em>e</em><em>e</em><em>m</em><em>p</em><em>t</em><em>a</em><em>b</em><em>l</em><em>e</em>(<em>J</em>)</span> <span class="ltx_text ltx_font_bold">then</span>
</div>
<div id="alg1.l11" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">11:</span></span>      <span class="math inline"><em>s</em><em>c</em><em>h</em><em>e</em><em>d</em>(<em>J</em>, <em>V</em>)</span>
</div>
<div id="alg1.l12" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">12:</span></span>     <span class="ltx_text ltx_font_bold">else</span> <span class="ltx_text ltx_font_bold">if</span> <span class="math inline"><em>J</em>.<em>r</em><em>e</em><em>q</em><em>u</em><em>e</em><em>s</em><em>t</em> &lt; <em>V</em>.<em>f</em><em>r</em><em>e</em><em>e</em>_<em>r</em><em>e</em><em>s</em><em>o</em><em>u</em><em>r</em><em>c</em><em>e</em></span> <span class="ltx_text ltx_font_bold">then</span>
</div>
<div id="alg1.l13" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">13:</span></span>      <span class="math inline"><em>s</em><em>c</em><em>h</em><em>e</em><em>d</em>(<em>J</em>, <em>V</em>)</span>
</div>
<div id="alg1.l14" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">14:</span></span>     <span class="ltx_text ltx_font_bold">else</span> <span class="ltx_text ltx_font_bold">if</span> <span class="math inline"><em>p</em><em>r</em><em>e</em><em>d</em>_<em>J</em><em>C</em><em>T</em>(<em>J</em>) &lt; <em>a</em><em>r</em><em>r</em><em>i</em><em>v</em><em>a</em><em>l</em>_<em>t</em><em>i</em><em>m</em><em>e</em></span> <span class="ltx_text ltx_font_bold">then</span>
</div>
<div id="alg1.l15" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">15:</span></span>      <span class="math inline"><em>s</em><em>c</em><em>h</em><em>e</em><em>d</em>(<em>J</em>, <em>V</em>)</span>
</div>
<div id="alg1.l16" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">16:</span></span>     <span class="ltx_text ltx_font_bold">else</span>
</div>
<div id="alg1.l17" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">17:</span></span>      <span class="math inline"><em>O</em>.<em>a</em><em>d</em><em>d</em>(<em>J</em>)</span>
</div>
<div id="alg1.l18" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">18:</span></span>     <span class="ltx_text ltx_font_bold">end</span> <span class="ltx_text ltx_font_bold">if</span>
</div>
<div id="alg1.l19" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">19:</span></span>   <span class="ltx_text ltx_font_bold">end</span> <span class="ltx_text ltx_font_bold">while</span>
</div>
<div id="alg1.l20" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">20:</span></span>   <span class="math inline"><em>Q</em> ← <em>O</em></span>
</div>
<div id="alg1.l21" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">21:</span></span>   <span class="math inline"><em>s</em><em>l</em><em>e</em><em>e</em><em>p</em>(<em>i</em><em>n</em><em>t</em><em>e</em><em>r</em><em>v</em><em>a</em><em>l</em>)</span>
</div>
<div id="alg1.l22" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">22:</span></span>  <span class="ltx_text ltx_font_bold">end</span> <span class="ltx_text ltx_font_bold">while</span>
</div>
<div id="alg1.l23" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span class="ltx_text" style="font-size:80%;">23:</span></span><span class="ltx_text ltx_font_bold">end</span> <span class="ltx_text ltx_font_bold">function</span>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_float"><span class="ltx_text ltx_font_bold">Algorithm 1</span> </span> Job scheduling policy</figcaption>
</figure>

::: {#A7.p2 .ltx_para}
[JCT predictor.]{.ltx_text .ltx_font_bold} The JCT prediction enables opportunistically scheduling short-lived jobs to the reserved resources as long as they can finish before the arrival of LPJ. This helps improve resource utilization and decrease queuing delay. The prediction is based on metadata associated with jobs, such as the number of requested CPUs and GPUs, the requested amount of drives, the department of task owners, etc. Although estimating the exact JCT is inherently difficult, we adopt a coarse-grained forecasting strategy, which classifies the JCT into different time intervals.
:::

::: {#A7.p3 .ltx_para}
To train the JCT predictor, we retrieve historical trace data from the database. Then, we pre-process the data, such as removing jobs that are early killed by users, and divide the JCT into 10-minute intervals. We then train models to predict the interval into which incoming jobs may fall by the metadata associated with the jobs. We tried both a deep neural network (DNN) and a gradient boosting predictor (GBM)\[[20](#bib.bib20){.ltx_ref}\], and found that GBM achieves higher performance, likely due to its ability to handle categorical variables.
:::

<figure id="A7.F13" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="A7.F13.sf1" class="ltx_figure ltx_figure_panel">
<img src="/html/2509.15940/assets/x17.png" id="A7.F13.sf1.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="831" height="468" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(a)</span> </span><span class="ltx_text" style="font-size:90%;">GBM: RMSE 1.61.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="A7.F13.sf2" class="ltx_figure ltx_figure_panel">
<img src="/html/2509.15940/assets/x18.png" id="A7.F13.sf2.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="830" height="466" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(b)</span> </span><span class="ltx_text" style="font-size:90%;">DNN: RMSE: 2.12.</span></figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 13</span>: </span><span class="ltx_text" style="font-size:90%;">JCT prediction.</span></figcaption>
</figure>

::: {#A7.p4 .ltx_para}
To demonstrate effectiveness, we extract 4-month trace data and divide them into a training set ($90\%$) and a test set ($10\%$). Figure [[13]{.ltx_text .ltx_ref_tag}](#A7.F13 "Figure 13 ‣ Appendix G Queue Management ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} shows an example of prediction in the test set. We apply randomized grid search to optimize hyper-parameters and also use bagging to determine uncertainty estimation. We observe that the RMSE is $1.61$ in the test set, and recent studies suggest that such prediction could help scheduling decision \[[43](#bib.bib43){.ltx_ref}, [3](#bib.bib3){.ltx_ref}\].
:::
:::::::

:::::: {#A8 .section .ltx_appendix}
## [Appendix H ]{.ltx_tag .ltx_tag_appendix}Evaluation of Queue Management {#appendix-h-evaluation-of-queue-management .ltx_title .ltx_title_appendix}

<figure id="A8.F14" class="ltx_figure">
<img src="/html/2509.15940/assets/x19.png" id="A8.F14.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="497" height="318" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 14</span>: </span><span class="ltx_text" style="font-size:90%;">Allocation and retention rate over time since a LLM pre-training job is planned.</span></figcaption>
</figure>

::: {#A8.p1 .ltx_para}
We collect job traces and replay in the trace-driven simulator. Figure [[14]{.ltx_text .ltx_ref_tag}](#A8.F14 "Figure 14 ‣ Appendix H Evaluation of Queue Management ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} shows the allocation and retention rate over time. The allocation rate is determined by dividing the total number of nodes by the nodes with allocated jobs. The retention rate measures how many planned nodes for the LLMs job are occupied by other jobs, which inevitably requires manual preemption when the LLMs job arrives. At $18:00$, Arnold is told the arrival time of the LLMs job at $22:00$, and so it triggers the code path to plan and reserve resources accordingly. We use a default bin-pack algorithm to schedule other jobs.
:::

::: {#A8.p2 .ltx_para}
The allocation rate is $0.9$ initially and then gradually decreases to below $0.5$ due to the need to reserve $1200$+ nodes in the cluster. At the beginning of the imminent period, the retention rate is relatively high and matches the allocation rate because the previous allocated jobs are not aware of the incoming LLMs pre-training job. Then the retention rate decreases faster than the allocation rate since nodes have been reserved, and is close to $0$ at the end of the period, showing the effectiveness of the scheduling policy.
:::

::: {#A8.p3 .ltx_para}
For the reserving-and-packing policy \[[43](#bib.bib43){.ltx_ref}\], it does not offer strong semantics for reservation (i.e. best effort). Thus, the scheduler will not be able to generate a feasible solution as the LLM job arrives, not to mention optimal placement (orange line). The JCT prediction navigates the trade-off space between resource utilization and guarantees by scheduling opportunistically short-lived jobs to the reserved zone. In its absence, both queuing delays and resource idle times increase, as indicated by the green line.
:::
::::::

:::::::: {#A9 .section .ltx_appendix}
## [Appendix I ]{.ltx_tag .ltx_tag_appendix}Break-down analysis {#appendix-i-break-down-analysis .ltx_title .ltx_title_appendix}

<figure id="A9.F15" class="ltx_figure">
<img src="/html/2509.15940/assets/x20.png" id="A9.F15.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="497" height="185" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 15</span>: </span><span class="ltx_text" style="font-size:90%;">Breakdown analysis of the ablation experiment.</span></figcaption>
</figure>

::: {#A9.p1 .ltx_para}
We compare aggregated kernel-level metrics by summing the duration for each type of kernel. Figure [[10]{.ltx_text .ltx_ref_tag}](#S7.F10 "Figure 10 ‣ 7.2 Cluster Experiment ‣ 7 Evaluation ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} illustrates the speedup achieved by Arnold (in green) and the slowdown (in red) of the full-scale experiment. We only report kernels whose difference is significant. The most significant speedup is the broadcast kernel ($10\%$), which is the optimized P2P implementation of our communication library. However, the speedup is slightly offset by the slowdown of a reduce-scatter kernel and even a computational kernel. The slowdown contradicts our expectations, as Arnold's scheduling also reduces the spread of DP groups. Moreover, we have not changed other configurations, so the slowdown of the GEMM kernel is unexpected.
:::

::: {#A9.p2 .ltx_para}
After thorough investigation, we suspect the slowdown is due to the interference between GPUs' streams. Due to hybrid parallelism, GPUs maintain multiple streams that issue operations concurrently during training. Although overlapping computation with communication indicates good performance optimization, it also causes resource contention and interference.
:::

::: {#A9.p3 .ltx_para}
[Network topology affects computation kernels.]{.ltx_text .ltx_font_bold} To investigate the counter-intuitive results, we isolate the impact of streams by modifying NCCL. For example, we add additional environmental variables such as [NCCL_DP_MIN_NCHANNELS]{.ltx_text .ltx_font_italic} to have fine-grained controls on the DP stream. We disable channel auto-tuning and rerun jobs with and without setting the NCCL variable. Figure [[15]{.ltx_text .ltx_ref_tag}](#A9.F15 "Figure 15 ‣ Appendix I Break-down analysis ‣ Efficient Pre-Training of LLMs via Topology-Aware Communication Alignment on More Than 9600 GPUs"){.ltx_ref} shows the breakdown analysis. Communication kernels have speedups by setting the NCCL variable, whereas computation kernels have slowdown. Since the only change is the NCCL variable, it indicates if we allocate more GPU SMs to communication, computation kernels suffer from performance loss for less available SMs.
:::

::: {#A9.p4 .ltx_para}
In production training, the NCCL variables are dynamically auto-tuned, so given that network topology-optimized scheduling influences the communication of DP and PP groups, ultimately it causes variations in the computation kernels.
:::

::: {#A9.p5 .ltx_para}
:::
::::::::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

::: ar5iv-footer
[◄](/html/2509.15939){.ar5iv-nav-button .ar5iv-nav-button-prev} [![ar5iv homepage](/assets/ar5iv.png){height="40"}](/){.ar5iv-home-button} [Feeling\
lucky?](/feeling_lucky){.ar5iv-text-button} [](/land_of_honey_and_milk){rel="nofollow" aria-hidden="true" tabindex="-1"} [Conversion\
report](/log/2509.15940){.ar5iv-text-button .ar5iv-severity-ok} [Report\
an issue](https://github.com/dginev/ar5iv/issues/new?template=improve-article--arxiv-id-.md&title=Improve+article+2509.15940){.ar5iv-text-button target="_blank"} [View original\
on arXiv](https://arxiv.org/abs/2509.15940){.ar5iv-text-button .arxiv-ui-theme}[►](/html/2509.15941){.ar5iv-nav-button .ar5iv-nav-button-next}
:::

[[]{.color-scheme-icon}](javascript:toggleColorScheme() "Toggle ar5iv color scheme"){.ar5iv-toggle-color-scheme} [Copyright](https://arxiv.org/help/license){.ar5iv-footer-button target="_blank"} [Privacy Policy](https://arxiv.org/help/policies/privacy_policy){.ar5iv-footer-button target="_blank"}

::: ltx_page_logo
Generated on Tue Oct 7 00:02:51 2025 by [[L[a]{.ltx_font_smallcaps style="position:relative; bottom:2.2pt;"}T[e]{.ltx_font_smallcaps style="font-size:120%;position:relative; bottom:-0.2ex;"}]{style="letter-spacing:-0.2em; margin-right:0.1em;"}[XML]{style="font-size:90%; position:relative; bottom:-0.2ex;"}![Mascot Sammy](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAsAAAAOCAYAAAD5YeaVAAAAAXNSR0IArs4c6QAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB9wKExQZLWTEaOUAAAAddEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIFRoZSBHSU1Q72QlbgAAAdpJREFUKM9tkL+L2nAARz9fPZNCKFapUn8kyI0e4iRHSR1Kb8ng0lJw6FYHFwv2LwhOpcWxTjeUunYqOmqd6hEoRDhtDWdA8ApRYsSUCDHNt5ul13vz4w0vWCgUnnEc975arX6ORqN3VqtVZbfbTQC4uEHANM3jSqXymFI6yWazP2KxWAXAL9zCUa1Wy2tXVxheKA9YNoR8Pt+aTqe4FVVVvz05O6MBhqUIBGk8Hn8HAOVy+T+XLJfLS4ZhTiRJgqIoVBRFIoric47jPnmeB1mW/9rr9ZpSSn3Lsmir1fJZlqWlUonKsvwWwD8ymc/nXwVBeLjf7xEKhdBut9Hr9WgmkyGEkJwsy5eHG5vN5g0AKIoCAEgkEkin0wQAfN9/cXPdheu6P33fBwB4ngcAcByHJpPJl+fn54mD3Gg0NrquXxeLRQAAwzAYj8cwTZPwPH9/sVg8PXweDAauqqr2cDjEer1GJBLBZDJBs9mE4zjwfZ85lAGg2+06hmGgXq+j3+/DsixYlgVN03a9Xu8jgCNCyIegIAgx13Vfd7vdu+FweG8YRkjXdWy329+dTgeSJD3ieZ7RNO0VAXAPwDEAO5VKndi2fWrb9jWl9Esul6PZbDY9Go1OZ7PZ9z/lyuD3OozU2wAAAABJRU5ErkJggg==)](http://dlmf.nist.gov/LaTeXML/){.ltx_LaTeXML_logo target="_blank"}
:::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
