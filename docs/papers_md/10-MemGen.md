      Under review as a conference paper at ICLR 2026




000
001
      M EM G EN : W EAVING G ENERATIVE L ATENT M EMORY
002   FOR S ELF -E VOLVING AGENTS
003
004
005    Anonymous authors
006    Paper under double-blind review
007
008
009                                                A BSTRACT
010
011            Agent memory shapes how Large Language Model (LLM)-powered agents, akin
012            to the human brain, progressively refine themselves through environment inter-
013            actions. Existing paradigms remain constrained: parametric memory forcibly
014            adjusts model parameters, and retrieval-based memory externalizes experience
015
               into structured databases, yet neither captures the fluid interweaving of reason-
               ing and memory that underlies human cognition. To address this gap, we pro-
016
               pose MemGen, a dynamic generative memory framework that equips agents with
017
               a human-esque cognitive faculty. It consists of a memory trigger, which moni-
018            tors the agent’s reasoning state to decide explicit memory invocation, and a mem-
019            ory weaver, which takes the agent’s current state as stimulus to construct a la-
020            tent token sequence as machine-native memory to enrich its reasoning. In this
021            way, MemGen enables agents to recall and augment latent memory throughout rea-
022            soning, producing a tightly interwoven cycle of memory and cognition. Exten-
023            sive experiments across eight benchmarks show that MemGen surpasses leading
024            external memory systems such as ExpeL and AWM by up to 38.22%, exceeds
025            GRPO by up to 13.44%, and exhibits strong cross-domain generalization abil-
026
               ity. More importantly, we find that without explicit supervision, MemGen sponta-
               neously evolves distinct human-like memory faculties, including planning mem-
027
               ory, procedural memory, and working memory, suggesting an emergent trajec-
028
               tory toward more naturalistic forms of machine cognition. Codes are available at
029            http://anonymous.4open.science/r/mem-gen.
030
031
032   1 I NTRODUCTION
033   The ascent of Large Language Model (LLM)-powered agents marks a paradigm shift across diverse
034   domains (Luo et al., 2025b; Yang et al., 2024b; Qian et al., 2025; Singh et al., 2025; Pantiukhin et al.,
035   2025; Ren et al., 2025). Pivotal to this success is the concept of agent memory (Zhang et al., 2024b;
036   Wu et al., 2025b), which enables LLM agents to learn progressively from environmental interac-
037   tions (Zhang et al., 2025a; Qiu et al., 2025b). Crucially, this conception of agent memory extends
038   beyond that of conversational agents (i.e., personalized memory (Wu et al., 2025b)), whose primary
039
      role is to sustain coherence across long-horizon, multi-turn dialogues (Chhikara et al., 2025; Xu
      et al., 2025a; Packer et al., 2024; Zhong et al., 2023). Rather, the scope of this paper is primarily on
040
      enabling agents to internalize experience, simulate human-like cognitive iteration, and progressively
041
      enhance problem-solving competence (Gao et al., 2025).
042
043   The memory serving as this self-evolving engine typically manifests in two dominant paradigms.
044   The first is (I) parametric memory, which internalizes experiences by directly updating agents’ pa-
045
      rameters (Yao et al., 2024; Yang et al., 2023; Zeng et al., 2023; Chen et al., 2024b; 2025). While this
      approach can yield substantial performance gains, its reliance on parameter modification inevitably
046
      entails catastrophic forgetting, i.e., the erosion of general knowledge (Dou et al., 2024). Conversely,
047
      the second paradigm is (II) retrieval-based memory, which externalizes past experiences into a
048   structured database, such as (i) raw trajectories (Luo et al., 2025a; Zhang et al., 2025a; Zhao et al.,
049   2024), (ii) high-level experiences (Zhao et al., 2024; Tang et al., 2025; Fang et al., 2025; Wang et al.,
050   2024c), and (iii) condensed skills like reusable APIs (Zheng et al., 2025) or MCP boxes (Qiu et al.,
051   2025b;a). Although this non-invasive approach circumvents catastrophic forgetting, its efficacy is
052   fundamentally tethered to context engineering. It adheres to a rigid execution pipeline, providing
053   retrieved context to the agent without achieving the fluid, seamless integration characteristic of truly
      internalized memory (Su et al., 2025).


                                                         1
      Under review as a conference paper at ICLR 2026




054
       Parametric Memory                                 Latent Memory (for MemGen)
055
056                       Rollout                                        Task Query:     Which U.S. state has the longest
           Agent                                                          coastline? You can search the Internet.
                                          Experience
057        Policy                            Data                     I will use                Refine
                          Reward
                                                                                      Latent              Latent         Latent
                                                                                                                                   The final
                                                                      an iterative              search             ...
058                                                                   search                    keyword
                                                                                                                                   answer is
                                                                                                                                   \boxed{Alaska}.
                    SFT, GRPO, DPO, etc.                              paradigm...               as [..]
059
                                                              Think      Experience              External                         Formatting
060                                                       Paradigm         Retreive            Knowledge                          Experience
061    Retrieval-based Memory
           Task                                           Latent
062                                 Previous              Space
           Query                    Trajectories
063          In-context    A-Mem, Mem0, AWM, etc.
064          prompting              Condensed
065                                 Insights & Exp.            Internal Memory                          External Memory
           Agent
                                                           Seen tasks: PopQA, TriviaQA ..            RAG system:        ...
066        Policy          Expel, Mobile-Agent-E, etc.
067
068      Figure 1: The comparison among parametric memory, retrieval-based memory and MemGen.
069
      Given these deficiencies, latent memory offers a compelling alternative, leveraging latent states as a
070
      machine-native, high-density medium for memory. Existing approaches either use the (i) key-value
071
      (KV) cache to maintain dynamic memory set (Gim et al., 2024; Jin et al., 2024; Hongkang Yang
072   et al., 2024), yet which is primarily confined to addressing long-context issues, or (ii) latent token
073   embeddings to store agent experiences (Wang et al., 2024b; 2025b), which still rely on invasive
074   LLM parameter updates. More fundamentally, these methods diverge from human cognition in two
075   critical dimensions: they lack the seamless interleaving of reasoning and memory, a process where
076   thought and memory dynamically reshape one another, and remain largely retrieval-based, fetching
077   memories by embedding similarity (Wang et al., 2024b) rather than generatively reconstructing them
078   into novel, coherent insights. This leads to the pivotal research question that motivates our work:
079    How can we architect agent memory as a dynamic cognitive faculty, capable of fluid, reconstructive pro-
080    cesses that interweave seamlessly with reasoning?
081
082   To address this challenge, we introduce MemGen, a dynamic and generative memory framework de-
083
      signed to endow any LLM agent with a more human-esque cognitive faculty. At its core, MemGen
      continuously monitors an agent’s cognitive state, enabling it to dynamically invoke a generative pro-
084
      cess that synthesizes a bespoke latent memory at any critical juncture during its reasoning process.
085
      Practically, MemGen comprises two synergistic components: a reinforcement learning (RL)-trained
086   ♣ memory trigger, which acts as a metacognitive monitor to discern the opportune moments for
087   explicit memory invocation; and a ♠ memory weaver, which takes the agent’s current state as a
088   stimulus to draw upon relevant implicit parametric memory (potentially augmented with externally
089   retrieved information) and then reconstructs this synthesis into a succinct, machine-native latent
090   memory. With the reasoning core fixed, MemGen inherently mitigates catastrophic forgetting when
091   exposed to new data, and equips agents with a generative memory deeply integrated with reasoning.
092
093
      Experimental Observation. Extensive experiments across nine benchmarks and four baseline
094   categories demonstrate that MemGen delivers ❶ substantial performance gains, with improve-
095   ments of up to 31.7% on ALFWorld (Shridhar et al., 2021) and 27.1% on KodCode (Xu et al.,
096   2025d) with Qwen3-8B, surpassing parametric memory (REINFORCE++, +5.8%) and the GRPO
097   method (+5.32%);❷ strong cross-domain generalization, where training in the math domain not
098   only avoids degradation elsewhere but also boosts performance in science reasoning (+6.06%) and
099   code generation (+5.1%); and ❸ continual learning ability, maintaining stable performance in
100   previously trained domains even after fine-tuning on three additional ones.
101
      Analysis & Interpretation. Beyond quantitative evaluation, we also analyze the functional be-
102   havior of the latent memories learned by MemGen. Through post-hoc interventions examining the
103   impact of removing specific latent memory on different agent failure modes, we found that differ-
104   ent latent memory tokens play distinct computational roles within the agent’s reasoning process,
105    including ❶ planning memory, where certain latent tokens specifically support high-level task
106   planning, ❷ procedural memory, where some latent memory tokens facilitate the agent’s recall
107   of task-specific procedural skills, such as tool usage and answer formatting, and ❸ working mem-
      ory, where certain tokens help the agent maintain coherence and understanding over long contexts


                                                              2
      Under review as a conference paper at ICLR 2026




108
      within a single task session. These specializations strongly reveal that MemGen endows the agent
109   with precise, functionally distinct memory.
110
111   2    R ELATED W ORK
112
113   LLM & Agent Memory. As outlined in §1, existing memory mechanisms designed to evolve the
114   problem-solving capacity of LLM agents can be broadly categorized into three classes: (I) para-
115
      metric memory, which either integrates past experiences directly into agent parameters through
      finetuning, as in FireAct (Chen et al., 2023), AgentLumos (Yin et al., 2024), and others (Zhang
116
      et al., 2024a; Fu et al., 2025), or maintains them in external parameter modules (Tack et al., 2024;
117
      Wang et al., 2024a); (II) retrieval-based memory, which abstracts prior experiences into trans-
118   ferable knowledge (Zhang et al., 2025a; Zhao et al., 2024), or distills them into reusable tools and
119   skills (Zheng et al., 2025; Wang et al., 2025c; Qiu et al., 2025b;a); and (III) latent memory, which
120   leverages implicit representations to encode and retrieve experience (Wang et al., 2024b; 2025b;
121   Hu et al., 2025b; Liu et al., 2024; Sun et al., 2025). Our MemGen falls within the latent memory
122   paradigm, yet distinguishes itself from prior approaches through its interweaving of reasoning and
123   memory, as well as generative reconstruction mechanism.
124
      Latent Computation. Our method is also closely related to latent computation, wherein latent
125   states are employed to intervene in or reshape the LLM’s reasoning process (Zhu et al., 2025).
126   Prominent paradigms include: (I) architecturally enabling native latent reasoning, exemplified
127   by Coconut (Hao et al., 2024), CODI (Shen et al., 2025b), LatentR3 (Zhang et al., 2025b) and
128   CoLaR (Tan et al., 2025), which render the LLM’s inference process inherently latent and machine-
129   native; and (II) employing latent computation to steer LLM generation, as in LaRS (Xu et al.,
130   2023), LatentSeek (Li et al., 2025a), SoftCoT (Xu et al., 2025c;b), and Coprocessor (Liu et al.,
131   2024), which leverage latent representations to modulate the quality of generated outputs.
132
133   3    P RELIMINARY
134   Notation. We formalize the agent’s interaction within an environment E. An agent, powered by
135   an LLM parameterized by θ, is denoted as πθ . For a given task x, the agent’s interaction unfolds as
136   a high-level trajectory, denoted as follows τ = (s0 , a0 , s1 , a1 , . . . , sT ), where st represents the state
137   of the environment and at is the high-level action taken by the agent. More internally, each action
138   at is essentially a sequence of tokens, at = (zt,1 , zt,2 , . . . , zt,Lt ), generated autoregressively by
139   the LLM. The generation of the j-th token is conditioned on the current state st and all previously
140   generated tokens within that action:
141                                             zt,j ∼ πθ (· | st , zt,<j ).                                      (1)
142
143   After an entire action sequence at is generated, it is executed in the environment, which transitions
144   the state from st to st+1 . The success of the trajectory τ is evaluated by a reward function R(τ ).
145   Problem Formalization Given a history of past experiences H = {(xi , τi )}N   i=1 , the objective is to
146   leverage this history to maximize the agent’s performance on new tasks. The policy πθ and a memory
147   system M are thus jointly optimized to maximize the expected reward over a task distribution D:
148                                                                
                                          max Ex∼D, τ ∼πθ,M R(τ ) ,                                       (2)
149                                           θ,M
150
151
      during which M is to produce a memory representation, m, which conditions the agent’s policy. The
      action at any timestep t is thus sampled as at ∼ πθ (· | st , mt ), where mt is the inserted memory at
152
      that step. Crucially, the nature and timing of memory generation, which we denote as the function
153
      fM , vary across different paradigms. We express the generation of the memory mt as:
154
155                                             mt = fM (st , H, m<t ),                                           (3)
156
      which accommodates diverse memory invocation granularities. For task-level memory (e.g., Ex-
157
      pel (Zhao et al., 2024) and G-Memory (Zhang et al., 2025a)), fM is invoked only at t = 0, and
158   mt = m0 for all subsequent steps. For step-level memory (e.g., AgentKB (Tang et al., 2025)),
159   fM is invoked at every step t to update the memory. In parametric memory, the influence of H is
160   compiled into θ, rendering memory generation implicit in the model parameters. Our work, which
161   introduces dynamic latent memory, focuses on designing a more fine-grained fM that decides for
      itself the optimal moments to regenerate mt at the token level during the agent’s reasoning process.

                                                             3
      Under review as a conference paper at ICLR 2026




162




      Explicit Reason
163                     Reaonser   Task Query: Hello, please find a flight from JFK to LAX for next Friday
                                   (should be earlier than 12am), and book me a ride from the airport.
164
165                                Okay, to achieve this, I'll start                           Wait, I can                    It seems that
                                                                                      Latent                                                     Latent
                                   by searching for the cheapest                               do this
                                                                                                                  Env.        the API is
166                                                                                                               Obs.
                                   flight on Google Flights.                                   faster by..                    not usable.
167
168                                                                                                                                     Context
169                      Trigger         Context                                                    Notations



                                                                                                                                                          Memory Insert
                                                                No Memory
                                                                                                  Prompt        Completion
170                                                             Invocation
                                                                                                  tokens        tokens
171                                       ...                                                     Latent memory tokens




      Implicit Memory
                                                                                                                                Token Sequence
172                                Token Sequence                                                     Latent token
                                                                Memory                         Latent sequence
173                                                             Invocation                                                      Memory
                                                                                                                              Invocation




                                                                             Memory Insert
174
175                      Weaver        Memory                                                   External Mem DB                       Memory
                                       Weaver Model                                                                                   Weaver Model
176                                                                                               Raw Trajectories
                                                              Parametric
177                                                            Memory
                                           Latent Memory Tokens                                  High-level abstracts        Latent Memory Tokens
178
179                                             ...                                                                             ...

180
181                                  Figure 2: The overview of our proposed MemGen.
182
183   4                 M ETHODOLOGY
184   4.1 MemGen: I NTERLEAVING M EMORY AND R EASONING
185
      Just as a person is the sum of their past experiences (McAdams, 2001), memory profoundly shapes
186   an agent’s behavior (Xiong et al., 2025). Yet existing agent memory systems lack the adaptive
187   fluidity characteristic of human cognition. Human reasoning and recollection operate in a continuous
188   interplay, whereas most agent memory frameworks retrieve information only once at task initiation
189   and append it to the query in a coarse and static manner. MemGen is designed to close this gap by
190   enabling memory to participate dynamically in the reasoning process.
191
      As shown in Figure 2, the reasoning process in an agent equipped with MemGen unfolds autore-
192   gressively, driven by a frozen core LLM, the reasoner πθ . For a given state st , πθ generates the
193   action at = (zt,1 , . . . , zt,Lt ). MemGen continuously monitors the token-by-token generation process
194   and performs on-demand memory insertion. At each token-generation step j, a memory trigger
195   Ttrigger monitors the reasoner’s internal cognitive state to determine if a moment of reflection, i.e., a
196   memory invocation, is necessary. Specifically, as the reasoner generates the token sequence zt,<j , it
197   produces a corresponding sequence of hidden state vectors, Ht,<j = (ht,1 , . . . , ht,j−1 ), where each
198   ht,k ∈ Rdmodel . The trigger takes the current hidden states Ht,<j as a representation of the reasoner’s
199   current metacognitive state and computes an invocation probability:
200                                             pj = σ (Ttrigger (ht,1 , . . . , ht,j−1 )) ,                                                                   (4)
201
      from which a binary decision, dj ∼ Bernoulli(pj ) ∈ {INVOKE, SKIP}, is sampled. If the decision
202
      is to [SKIP], πθ proceeds with its standard autoregressive generation, i.e., zt,j ∼ πθ (· | st , zt,<j ).
203
      However, if the decision is to INVOKE, the reasoning process is momentarily paused. This summons
204   the second core component of our framework: the memory weaver Wweaver , which takes the same
205   cognitive state Ht,<j as a stimulus to perform a generative act of recollection. It synthesizes a
206   bespoke, machine-native latent memory, formalized as Mt ∈ RK×dmodel with fixed length K:
207
                                    Mt := [mt,1 , mt,2 , · · · , mt,K ] = Wweaver (Ht,<j ),                                                                    (5)
208
209   where the memory is generated not merely from the parametric knowledge encoded within Wweaver
210   but may also incorporate cues retrieved from external memory databases (detailed implementation
211   is elaborated in §4.3). Crucially, Mt is not a verbatim restatement of prior content but a selective
212
      reconstruction, filtered and integrated through Wweaver , akin to the memory consolidation process in
      human brain (Spens & Burgess, 2024). Once formed, the latent memory is woven seamlessly into
213
      the reasoner’s ongoing dynamics: its hidden states are prepended to Ht,<j , upon which the reasoner
214
      resumes generation conditioned on this enriched context,
215
                                                      zt,j ∼ πθ (· | st , zt,<j , Mt ).                                                                        (6)


                                                                       4
      Under review as a conference paper at ICLR 2026




216
      This iterative cycle of generation, monitoring, invocation, weaving, and reintegration elevates rea-
217   soning from a linear unfolding to a recursive dialogue with memory, all without altering the frozen
218   reasoner πθ , and thereby preserving its general capabilities. In the following sections, we detail the
219   implementations of the memory trigger (▷ §4.2) and the memory weaver (▷ §4.3).
220
221   4.2   M EMORY T RIGGER : L EARNING TO I NVOCATE M EMORY
222
      In this section, we describe the concrete instantiation and training recipe of the memory trigger
223
      Ttrigger . Recall from §4.1 that Ttrigger serves as a metacognitive monitor, observing the evolving rea-
224   soning state of the frozen reasoner πθ and producing a binary decision dj ∈ {INVOKE, SKIP}(0 ↔
225   SKIP, 1 ↔ INVOKE) that determines whether the memory weaver should be invoked at token j.
226
227   Instantiation. We instantiate Ttrigger as a lightweight LoRA adapter attached to the reasoner πθ .
228
      At the decoding step j of the timestep t, it receives the sequence of all hidden states, Ht,<j ∈
229   R(j−1)×dmodel ; conditioned on this context, Ttrigger outputs the action probability P (dj = INVOKE) ∈
230
      [0, 1]. We use hidden states rather than textual context as input because, following prior work in
      latent reasoning (Hao et al., 2024; Shen et al., 2025a), latent embeddings retain richer context-
231
      sensitive information that would otherwise be lost after softmax decoding. To avoid excessive
232
      computational overhead, we adopt a sentence-granularity activation strategy, inspired by recent
233   studies on LLM interpretability (Anthropic, 2025; Chen et al., 2024a), which find that interventions
234   between sentences can more effectively guide LLMs’ reasoning path. Specifically, we define a
235   delimiter token set D (e.g., commas, periods) and let the trigger act only when the current token
236   falls in D. The invocation decision is computed as:
237                                                         
238                                                           0                 if zj ∈
                                                                                      / D,
                             dj = Bernoulli(pj ), pj =                                                   (7)
239                                                           Ttrigger (Ht,<j ) if zj ∈ D,
240
241   which ensures that Ttrigger is invoked only at semantically significant boundaries, preserving decod-
242
      ing efficiency. We validate that MemGen does not incur excessive inference delay in Appendix E.3.3.
243   Training Recipe. The memory trigger is trained via reinforcement learning, motivated by the
244   need to balance two competing desiderata: ensuring that critical latent memories are invoked to
245   improve task performance, while avoiding unnecessary or spurious invocations that could disrupt
246   reasoning or incur computational overhead. Given a batch of seen tasks H = {(xi , τi )}N           i=1 , the
247   frozen reasoner πθ generates candidate trajectories while the memory weaver Wweaver remains fixed.
248   At each activated step, the trigger selects an action d˜j ∈ {INVOKE, SKIP} and receives a reward
249   r(τi ) reflecting the quality of the resulting trajectory with respect to the task objective. To encourage
250   sparse yet strategically critical memory invocation, we introduce a reward-adaptive penalty, which
251
      discourages unnecessary activations while preserving essential ones, into the objective:
252                                                          h                X                       i
                                 max Eτ ∼π ,d̃∼T ϕ               R(τi ) − λ         max(0, d˜i,j − p̄) ,       (8)
253                               ϕ      i     θ   trigger
                                                                              i,j
254
255   where p̄ is computed as the mean activation probability across high-reward trajectories, i.e., those
256   with reward exceeding the batch median:
257
                               1       X 1 X
258                   p̄ =                        d˜i,j , Hhigh = {i : R(τi ) ≥ mediank (R(τk ))},             (9)
                             |Hhigh | i∈H |τi | j
259                                     high

260
261
      where ensures that Ttrigger learns to invoke memory selectively at key decision points, maximizing
      task reward while maintaining computational efficiency.
262
263
      4.3   M EMORY W EAVER : S YNTHESIZING AND I NSERTING L ATENT M EMORY
264
265   In this section, we elaborate on the weaver Wweaver , the memory carrier within the MemGen frame-
266
      work. When the agent assimilates new experiences, this information is exclusively internalized into
      the parameters of Wweaver , leaving the core reasoner πθ entirely unmodified. At junctures where the
267
      reasoner requires experiential support, a context-dependent hook activates the weaver to synthesize
268
      and externalize pertinent knowledge as a usable memory. To be more specific, recall from Equa-
269   tion (5) that after the Ttrigger signals the need for memory at step j, Wweaver accepts Ht,<j (as the
      hook) and generates a latent token sequence Mt (as the memory) for πθ .

                                                                     5
      Under review as a conference paper at ICLR 2026




270
      Instantiation. We instantiate Wweaver using anthoer LoRA adapter attached to πθ . Formally, given
271
      the incoming hook Ht,<j ∈ R(j−1)×dmodel , the weaver outputs a latent memory matrix: Mt =
272      θ′
      Wweaver (Ht,<j ) ∈ RK×dmodel , where K denotes the fixed length of the latent memory sequence and
273    ′
      θ are the trainable LoRA parameters. The synthesized Mt is first projected through a linear layer
274   to align it with the reasoner’s token-embedding space, and is then prepended to the current hidden
275   states of πθ to guide subsequent token generation, as described in Equation (6).
276
277   Training Recipe. The training of Wweaver proceeds over a batch of past trajectories H =
278
      {(xi , τi )}N
                  i=1 . Distinct from conventional agent tuning, which directly integrates experiential data
      into the parameters of πθ (Chen et al., 2025; Yin et al., 2024), MemGen internalizes experiential
279
      knowledge solely into Wweaver , which ensures that πθ ’s general capabilities remain intact.
280
281   Crucially, this separation makes MemGen agnostic to optimization strategies and compatible with
282   diverse LLM backbones. Whether employing supervised fine-tuning (SFT) or RL-based objectives
283   such as GRPO or DAPO, the weaver can be updated under a unified goal: optimizing the generation
                                                                                        W ,T
284   process of latent memory so as to maximize downstream reward. Formally, let Πθ θ′ (· | x) denote
285   the process of rolling out a trajectory for a task x by πθ in conjunction with weaver Wθ′ and trigger
286
      T . Given a reward functional R, the objective updates only θ′ by maximizing the expected reward:
                                                                                            
287                                max E(xi ,τi )∼H E       W ′ ,T                 R(xi , τ ) ,         (10)
                                    θlora               τ ∼Πθ θ      (·|xi )
288
289   where the gradients from R are propagated solely to θ′ , thereby equipping the weaver to supply
290   precisely the memories that improve end-to-end performance without altering πθ . Equation (10)
291   enables Wweaver to absorb diverse experiential signals and externalize them as dynamic, context-
292   sensitive latent memories, independent of the architectural or training paradigm of the base reasoner.
293   In practice, we first train the memory weaver using a random inserter as a lightweight surrogate for
294   the trigger, and then freeze the trained weaver and proceed to train the trigger. For a complete
295   description of the training procedure, please refer to Appendix D.1.
296   Integration with Retrieval-based Memory. Although the memory generation above primarily
297   draws on the weaver’s parametric knowledge, it can be combined with external memory sources.
298   When triggered, any retrieval-based system (e.g., MemoryBank, ExpeL) can provide textual mem-
299   ory, which is merged with the hook Ht,<j and fed into W to produce latent memory. This allows W
300   to integrate internal knowledge and external information, supplying the reasoner with richer memory
301
      support. Implementation details and results are placed in Appendix F.
302
303   5     E XPERIMENTS
304   In this section, we conduct extensive experiments to answer the following research questions:
305
      • (RQ1) Can MemGen surpass both parametric and retrieval-based memory?
306
      • (RQ2) Is the memory learnt by MemGen generalizable across task domains? And why?
307
      • (RQ3) Can MemGen facilitate continual learning and mitigate catastrophic forgetting?
308
      • (RQ4) Does MemGen implicitly evolve human-like memory hierachy?
309
310
      5.1   E XPERIMENTAL S ETUP
311
312   Evaluation and Benchmarks. Our evaluation covers nine datasets from five domains, including
313   ❶ web search: TriviaQA (Joshi et al., 2017) and PopQA (Mallen et al., 2023); ❷ embodied action:
314
      ALFWorld (Shridhar et al., 2021); ❸ math reasoning: AQuA (Ling et al., 2017), GSM8K (Cobbe
      et al., 2021), and MATH (Hendrycks et al., 2021); ❹ scientific reasoning: GPQA (Rein et al., 2023);
315
      and ❺ coding: KodCode (Xu et al., 2025d) and BigCodeBench (Jain et al., 2024).
316
317   Baselines. We compare MemGen against twelve baselines, categorized into four groups: (I)
318   Prompt-based methods: Vanilla model, CoT (Wei et al., 2023); (II) Parametric memory, where
319   experiential knowledge directly modifies model parameters via: SFT, GRPO (DeepSeek-AI et al.,
320
      2025), REINFORCE (Williams, 1992), REINFORCE++ (Hu et al., 2025a), Agent-FLAN (Chen
      et al., 2024b); (III) Retrieval-based memory, where processing tasks sequentially and storing the
321
      experiences in an external database, represented by MemoryBank (Zhong et al., 2023), ExpeL (Zhao
322
      et al., 2024), Agent Workflow Memory (AWM) (Wang et al., 2024c); and (IV) Latent computation,
323   where leveraging latent tokens as carriers of experiential knowledge, including SoftCoT (Xu et al.,
      2025c) and Co-processor (Liu et al., 2024).

                                                           6
      Under review as a conference paper at ICLR 2026




324   Table 1: Results on SmolLM3-3B and Qwen3-8B. All values represent the performance metric for
325   each task (e.g., accuracy %). We highlight the best and second best results.
326
      Backbone     Method         ALFWorld   TrivialQA   PopQA   KodCode   BigCodeBench   GPQA    GSM8K   MATH
327
                   Vanilla         18.96       10.47      8.23    37.05       35.96        9.35   47.63   16.22
328                CoT             17.60       12.88      9.95    38.45       39.42       20.70   58.91   56.33
329                SFT             32.36       55.25     37.22    59.25       40.79       19.70   63.48   45.65
                   GRPO            55.35       65.88     45.16    68.48       72.44       22.73   80.03   61.23
330
                   REINFORCE       53.13       63.20     46.81    65.53       67.14       23.44   82.03   58.75
331                REINFORCE++     53.95       63.20     44.10    65.90       68.80       22.73   81.50   59.89
332   SmolLM3-3B   Agent-FLAN      34.00       56.70     39.50    56.80       37.20       17.80   59.60   36.84
                   ExpeL           36.18       46.20     28.16    51.14       40.22       15.15   56.23   38.11
333                MemoryBank      32.80       43.30     25.81    44.50       31.80       10.20   58.30   43.53
334                AWM             40.50       49.80     29.60      -           -           -       -       -
335                SoftCoT         35.03       50.38     34.90    59.20       39.10       17.22   56.34   44.62
                   Co-processor    38.36       53.28     38.96    56.25       45.40       20.10   57.60   38.81
336
                   MemGen SFT      50.60       68.13     42.34    62.65       42.99       26.75   70.42   57.44
337                MemGen GRPO     63.60       79.30     58.60    72.85       74.24       25.20   83.47   63.65
338                Vanilla         58.93       52.18     34.13    49.10       33.33       38.18   89.48   79.82
339                CoT             57.10       53.80     33.20    51.25       35.59       35.15   87.67   78.24

340                SFT             83.59       74.55     51.12    64.75       41.33       40.33   90.76   81.35
                   GRPO            85.60       76.15     58.90    73.35       70.24       39.54   92.30   83.54
341                REINFORCE       82.10       75.22     57.96    72.11       70.20       37.12   91.25   83.27
342                REINFORCE++     84.80       75.90     58.30    72.90       71.88       37.68   91.90   85.24
      Qwen3-8B     Agent-FLAN      80.32       70.32     50.08    62.99       43.40       39.50   87.60   80.05
343
                   ExpeL           78.97       65.54     40.33    57.20       34.23       35.15   86.20   77.40
344                MemoryBank      70.41       60.56     41.60    56.39       40.61       35.66   90.35   80.35
345                AWM             80.33       69.30     43.69      -           -           -       -       -
                   SoftCoT         75.60       59.42     39.42    63.28       38.27       39.60   86.30   76.23
346                Co-processor    73.28       61.42     45.55    64.90       42.19       39.15   76.23   79.20
347
                   MemGen SFT      85.82       77.22     54.65    66.15       40.35       43.23   91.25   83.30
348                MemGen GRPO     90.60       80.65     62.30    76.16       75.56       40.24   93.20   88.24
349
350   Implementation Details. We select LLM backbones of varying sizes and sources, including
351   Qwen-2.5-1.5B (Yang et al., 2024a), HuggingFace’s SmolLM3-3B (HuggingFace, 2025), and
352   Qwen3-8B (Yang et al., 2025). The length of each latent memory sequence K is set among
353   {2, 4, 8}. As described in Equation (10), MemGen does not rely on a specific optimization algorithm,
354   so we implement two variants: MemGen SFT and MemGen GRPO , in which the weaver is updated using
355   SFT and GRPO signals. Details on these variants are provided in Appendix C. More training setup
356   and parameter configurations are listed in Appendix D.
357
358
      5.2   M AIN R ESTULS
359   [For RQ1] MemGen provides high-performing memory across domains. As shown in Tables 1
360   and 4, existing baselines exhibit clear limitations in cross-domain adaptivity. Retrieval-based memo-
361   ries (e.g., ExpeL, MemoryBank, AWM) occasionally surpass parazmetric tuning in embodied action;
362   for instance, AWM reaches 36.18% on ALFWorld with SmolLM3-3B, exceeding SFT by 3.15%.
363
      Yet their effectiveness deteriorates on reasoning-intensive tasks: ExpeL achieves only 8.12% on
      GPQA+Qwen2.5-1.5B, and even underperforms the vanilla model by 6.9% on TriviaQA, under-
364
      scoring its heavy reliance on backbone capacity. Parametric finetuning methods display the opposite
365
      tendency: they excel in structured domains such as code generation, where REINFORCE++ reaches
366   63.33% on KodCode with Qwen2.5-1.5B, but remain weak in knowledge-intensive reasoning,
367   with GPQA below 14%. In contrast, MemGen consistently advances performance across all do-
368   mains. For example, on ALFWorld+SmolLM3-3B, MemGen SFT and MemGen GRPO attain 50.60%
369   and 63.60%, improving over vanilla by 31.64% and 44.64%, respectively. Similar gains appear
370   with the larger Qwen3-8B, where MemGen GRPO achieves +27.06% on KodCode and +28.17%
371   on PopQA, surpassing GRPO by up to 3.4%. Overall, the dynamic memory insertion of MemGen
372   delivers substantial improvements across diverse task domains.
373   [For RQ2] MemGen Exhibits Strong Cross-Domain Generalization. To evaluate whether the
374   memory learned by MemGen can transfer across tasks, we train MemGen on one dataset and test it
375   on several others. We include two out-of-domain datasets, ScienceWorld (Wang et al., 2022) and
376   FEVER (Thorne et al., 2018), to further probe this. As shown in Figures 3, 9 and 10, baselines such
377   as SFT and MemoryBank achieve gains within their training domains (e.g., on ALFWorld, SFT
      +14.1% and MemoryBank +5.4% compared with vanilla), yet fail to generalize, with performance


                                                          7
      Under review as a conference paper at ICLR 2026




378
379
380
381
382
383
384
385
386
387
      Figure 3: The generalization study of MemGen. We train MemGen SFT on one dataset (ALFWorld or
      TriviaQA) and evaluate it on four datasets (TriviaQA, ALFWorld, ScienceWorld, and FEVER).
388
389   dropping sharply on FEVER by 16.2%. In contrast, MemGen not only attains substantial in-domain
390
      improvements (24.55% → 58.16% on KodCode, Figure 10), but also exhibits effective transfer:
      when trained on KodCode, performance on MATH rises from 36.6% → 54.2%. Having empirically
391
      validated MemGen’s generalizability, a question naturally arises: what underlies this capability?
392
393   [For RQ2] The Memory Trigger Intelligently De-                 120
                                                                                                   GSM8K (Avg: 86.07)
      termines When to Activate Memory Insertion, Mit-




                                                                           Memory Invocation Count
394                                                                                                GPQA (Avg: 46.87)
                                                                     100
395   igating Domain Conflict. After training MemGen on                                            KodCode (Avg: 21.02)
                                                                      80
396   GSM8K, we evaluate 150 samples each from GSM8K,
397   KodCode, and GPQA, visualizing the frequency with               60

398
      which the memory trigger invoked the memory weaver              40
      at each relative position in the model output. We ob-           20
399
      serve that the invocation frequency varies across domains
400                                                                    0
      and correlates directly with performance in Figure 9:              0   20        40        60         80          100
401   GSM8K exhibits the largest improvement (+19.64%)                       Relative Insertion Position (%)
402   and maximal invocations, GPQA achieves moderate gains Figure 4: Memory invocation frequency
403   (+6.06%) with medium invocations, and KodCode shows across benchmarks at inference (trained on
404   the smallest improvement (+3.1%) with the fewest invo- MemGen SFT +Qwen3-8B +GSM8K).
405   cations. This indicates that MemGen autonomously assesses, based on task-specific context, when
406   memory insertion will be beneficial, invoking the weaver less frequently in unfamiliar domains.
407   [For RQ3] MemGen Effectively Mitigates Catastrophic Forgetting. In Table 5, we sequentially
408   train on four datasets and evaluate on all benchmarks after each stage, where MemGen exhibits
409   stronger knowledge retention ability compared to baseline methods. For example, unlike SFT which
410   primarily improves performance on the most recent task (54.10% on KodCode but only 2.53% on
411
      GPQA), MemGen demonstrates more balanced cross-task generalization, attaining 38.43% on AQuA
      and 21.72% on GPQA after GSM8K training. Finally, it mitigates forgetting on earlier tasks, pre-
412
      serving 40.34% on AQuA following KodCode training compared to 27.14% for ExpeL and 28.61%
413
      for SFT, indicating a more stable continual learning ability. More analysis is in Appendix E.1.
414
415   5.3 F RAMEWORK A NALYSIS
416
      Having established the expressive capabilities of MemGen, we further investigate its underlying
417   mechanisms: what do the learned latent memories look like? Do they have specialized functions?
418
      [For RQ4] The Latent Memory Is Machine-Native and Human-Unreadable. We first visu-
419   alized the latent memory sequences learned by MemGen across different datasets using t-SNE in
420   Figures 5 and 11. As shown in Figure 5 (Left), sequences from distinct domains form separate
421   distributions, with related domains clustering closely (e.g., KodCode and BigCodeBench, GSM8K
422   and MATH). Examining latent memories within the same dataset, we observed pronounced clus-
423   tering patterns (as shown in Figure 5 (Middle and Right)). To explore potential commonalities
424   within these clusters, we forcibly decoded the latent tokens. Although the decoded sequences are
425   not human-readable, they exhibit intriguing regularities: many tokens within a cluster share struc-
426   tural conventions. For example, Cluster 0 in TriviaQA frequently follows the pattern “[...]SOC”,
427
      whereas Cluster 3 in GSM8K often adopts the format “[...] pick”. A large corpus of latent
      memory tokens is provided in Appendix G. Despite these sequences being machine-native and
428
      human-unreadable, we further investigate whether their underlying semantics can be interpreted.
429
430   [For RQ4] MemGen Implicitly Learns a Human-like Memory Hierarchy. To uncover the func-
431   tional roles of different latent memory clusters, we conducted a post-hoc intervention study. Follow-
      ing the taxonomy from (Song et al., 2025), we study eight distinct types of agent failure, including


                                                               8
      Under review as a conference paper at ICLR 2026




432                Latent Memory Visualization             Latent Memory Visualization      Latent Memory Visualization
                     (across all benchmarks)                      (on TriviaQA)                     (on GSM8K)
433
434
                                                                    INGER[…]
435                                                                                                     […]_check
                                  MATH                       гpaHИ[…]
436
                                    GSM8K
437                TriviaQA                                                                             […] kindergetAs […]
                                                                                          LTR[…]
438                                              KodCode
                                                                                 […]JaB
439                                      BCB
                 ALFWorld
440                                                                 […]SOC
                                                                                                     […]_pick

441
442
443
      Figure 5: (Left) t-SNE visualization of latent memories generated by MemGen +Qwen3-8B across datasets;
444
      (Middle and Right) Latent memory visualization within the TriviaQA and GSM8K datasets, clustered using
445   K-means. The text at each cluster center represents the common pattern shared by many memory sequences in
446   the cluster, such as Cluster 0 in GSM8K, where many sequences end with “ check”.
447
448
449
450
451
452
453
454
455
456
457
458
      Figure 6: (Left) Parameter sensitivity analysis on the latent memory length K; (Right) Effects of selectively
459   removing latent memory clusters on different agent failure modes on the TriviaQA dataset.
460
461   planning errors, tool response/parsing failures, answer formatting mistakes, etc, providing a struc-
462   tured framework to assess how memory influences performance. During evaluation, we selectively
463   removed latent tokens close to a specific cluster while keeping others intact, measuring the resulting
464   changes in these failure modes. Details on (1) the visualization process, (2) failure mode annotation,
465   and (3) token filtration are in Appendix H. As shown in Figure 6 (Right), distinct memory clusters
466
      exhibit varying influence on failure modes and can be mapped to different memory functions:
467   • Planning Memory supports high-level task planning and strategic reasoning. Removal of Cluster
468     2 substantially increases planning and compositional reasoning failures, indicating that this cluster
469     is crucial for guiding the LLM agent’s decision-making and sequencing of reasoning steps.
470   • Procedural Memory captures task-specific operational knowledge, such as tool usage and for-
471     matting ability. Cluster 3 corresponds to this role, as its removal leads to a marked increase in tool
472     response errors, parsing failures, and answer formatting mistakes.
473   • Working Memory manages the retention and effective use of prior context to maintain reason-
474     ing consistency. Clusters 1 and 4 contribute to this function: for instance, removing Cluster 1’s
475     memory tokens results in more frequent task misunderstandings and think-act inconsistency.
476   Nevertheless, these memory clusters are not entirely independent: for example, removing Cluster 1
477   also negatively affects planning ability, indicating that these memory faculties interact and jointly
478   enable the LLM to leverage past experience effectively. This analysis reveals that MemGen sponta-
479   neously organizes latent memory into a structured, human-like hierarchy.
480   Ablation Study & Sensitivity Analysis. We conduct a sensitivity analysis on the length of the
481   latent memory sequence K, as shown in Figure 6 (Left). It can be observed that as the latent token
482   length increases from 2 → 32, MemGen ’s performance correspondingly improves, likely reflecting
483   the expanded memory capacity. We then perform an ablation study on the memory trigger module
484   in Table 6, demonstrating the necessity of a dedicatedly trained trigger for effective memory invo-
485   cation. Furthermore, we analyze different training paradigms of the memory weaver in Table 7. For
      additional results and discussions, please refer to Appendix E.3.


                                                               9
      Under review as a conference paper at ICLR 2026




486
      Efficiency Analysis. To confirm that the memory insertion process of MemGen does not introduce
487   significant inference overhead, we show in Appendix E.3.3 that, while achieving up to 57.66%
488   performance improvement over vanilla LLMs, the per-query inference delay remains consistently
489   below the baseline, ranging from 24% to 94% of the vanilla LLM latency. This clearly demonstrates
490   that MemGen delivers substantial performance gains without compromising efficiency.
491
492   6   C ONCLUSION
493
      In this work, we introduced MemGen, a dynamic and generative memory framework designed for
494
      LLM Agents. By interleaving reasoning with memory synthesis through a reinforcement-learned
495   memory trigger and a generative memory weaver, MemGen transcends the limitations of paramet-
496   ric and retrieval-based paradigms. Extensive experiments showcase substantial performance gains,
497   robust cross-domain generalization, strong continual learning ability, and MemGen’s explicitly mod-
498   eled memory hierarchy (i.e., planning, procedural, and working memory). These results suggest a
499   promising path toward self-evolving LLM agents capable of fluid and reconstructive intelligence.
500
501
      E THICS S TATEMENT
502
503   This work presents a latent memory architecture designed for LLM agents. All experiments are con-
504   ducted on publicly available academic benchmarks across reasoning, scientific problem-solving, and
505   embodied tasks, without deployment in real-world decision-making or safety-critical applications.
506   Therefore, we believe this research does not pose direct ethical risks.
507
508   R EPRODUCIBILITY S TATEMENT
509
510   We aim to ensure the reproducibility of our work by releasing an anonymous repository linked in
511   the abstract, as well as detailing experimental settings (hyperparameters, LLM backbones, etc.) in
512
      §5.1 and appendix D.
513
514   R EFERENCES
515
      Anthropic. On the Biology of a Large Language Model. https://transformer-circuits.
516
        pub/2025/attribution-graphs/biology.html, 2025. [Accessed 24-08-2025].
517
518   Tianle Cai, Yuhong Li, Zhengyang Geng, Hongwu Peng, Jason D. Lee, Deming Chen, and Tri Dao.
519     Medusa: Simple llm inference acceleration framework with multiple decoding heads, 2024. URL
520     https://arxiv.org/abs/2401.10774.
521
      Baian Chen, Chang Shu, Ehsan Shareghi, Nigel Collier, Karthik Narasimhan, and Shunyu Yao.
522
        Fireact: Toward language agent fine-tuning, 2023. URL https://arxiv.org/abs/2310.
523
        05915.
524
525   Guoxuan Chen, Han Shi, Jiawei Li, Yihang Gao, Xiaozhe Ren, Yimeng Chen, Xin Jiang, Zhenguo
526     Li, Weiyang Liu, and Chao Huang. Sepllm: Accelerate large language models by compressing
527     one segment into one separator. arXiv preprint arXiv:2412.12094, 2024a.
528
      Zehui Chen, Kuikun Liu, Qiuchen Wang, Wenwei Zhang, Jiangning Liu, Dahua Lin, Kai Chen, and
529
        Feng Zhao. Agent-flan: Designing data and methods of effective agent tuning for large language
530     models, 2024b. URL https://arxiv.org/abs/2403.12881.
531
532   Zhixun Chen, Ming Li, Yuxuan Huang, Yali Du, Meng Fang, and Tianyi Zhou. Atlas: Agent tuning
533     via learning critical steps, 2025. URL https://arxiv.org/abs/2503.02197.
534
      Prateek Chhikara, Dev Khant, Saket Aryan, Taranjeet Singh, and Deshraj Yadav. Mem0: Building
535
        production-ready ai agents with scalable long-term memory, 2025. URL https://arxiv.
536     org/abs/2504.19413.
537
538   Karl Cobbe, Vineet Kosaraju, Mohammad Bavarian, Mark Chen, Heewoo Jun, Lukasz Kaiser,
539     Matthias Plappert, Jerry Tworek, Jacob Hilton, Reiichiro Nakano, et al. Training verifiers to
        solve math word problems. arXiv preprint arXiv:2110.14168, 2021.


                                                      10
      Under review as a conference paper at ICLR 2026




540
      DeepSeek-AI, Daya Guo, Dejian Yang, Haowei Zhang, Junxiao Song, Ruoyu Zhang, Runxin Xu,
541     Qihao Zhu, Shirong Ma, Peiyi Wang, Xiao Bi, Xiaokang Zhang, Xingkai Yu, Yu Wu, Z. F. Wu,
542     Zhibin Gou, Zhihong Shao, Zhuoshu Li, Ziyi Gao, Aixin Liu, Bing Xue, Bingxuan Wang, Bochao
543     Wu, Bei Feng, Chengda Lu, Chenggang Zhao, Chengqi Deng, Chenyu Zhang, Chong Ruan,
544     Damai Dai, Deli Chen, Dongjie Ji, Erhang Li, Fangyun Lin, Fucong Dai, Fuli Luo, Guangbo Hao,
545     Guanting Chen, Guowei Li, H. Zhang, Han Bao, Hanwei Xu, Haocheng Wang, Honghui Ding,
546     Huajian Xin, Huazuo Gao, Hui Qu, Hui Li, Jianzhong Guo, Jiashi Li, Jiawei Wang, Jingchang
547     Chen, Jingyang Yuan, Junjie Qiu, Junlong Li, J. L. Cai, Jiaqi Ni, Jian Liang, Jin Chen, Kai
548     Dong, Kai Hu, Kaige Gao, Kang Guan, Kexin Huang, Kuai Yu, Lean Wang, Lecong Zhang,
549
        Liang Zhao, Litong Wang, Liyue Zhang, Lei Xu, Leyi Xia, Mingchuan Zhang, Minghua Zhang,
        Minghui Tang, Meng Li, Miaojun Wang, Mingming Li, Ning Tian, Panpan Huang, Peng Zhang,
550
        Qiancheng Wang, Qinyu Chen, Qiushi Du, Ruiqi Ge, Ruisong Zhang, Ruizhe Pan, Runji Wang,
551
        R. J. Chen, R. L. Jin, Ruyi Chen, Shanghao Lu, Shangyan Zhou, Shanhuang Chen, Shengfeng
552     Ye, Shiyu Wang, Shuiping Yu, Shunfeng Zhou, Shuting Pan, S. S. Li, Shuang Zhou, Shaoqing
553     Wu, Shengfeng Ye, Tao Yun, Tian Pei, Tianyu Sun, T. Wang, Wangding Zeng, Wanjia Zhao, Wen
554     Liu, Wenfeng Liang, Wenjun Gao, Wenqin Yu, Wentao Zhang, W. L. Xiao, Wei An, Xiaodong
555     Liu, Xiaohan Wang, Xiaokang Chen, Xiaotao Nie, Xin Cheng, Xin Liu, Xin Xie, Xingchao Liu,
556     Xinyu Yang, Xinyuan Li, Xuecheng Su, Xuheng Lin, X. Q. Li, Xiangyue Jin, Xiaojin Shen, Xi-
557     aosha Chen, Xiaowen Sun, Xiaoxiang Wang, Xinnan Song, Xinyi Zhou, Xianzu Wang, Xinxia
558     Shan, Y. K. Li, Y. Q. Wang, Y. X. Wei, Yang Zhang, Yanhong Xu, Yao Li, Yao Zhao, Yaofeng
559     Sun, Yaohui Wang, Yi Yu, Yichao Zhang, Yifan Shi, Yiliang Xiong, Ying He, Yishi Piao, Yisong
560
        Wang, Yixuan Tan, Yiyang Ma, Yiyuan Liu, Yongqiang Guo, Yuan Ou, Yuduan Wang, Yue Gong,
        Yuheng Zou, Yujia He, Yunfan Xiong, Yuxiang Luo, Yuxiang You, Yuxuan Liu, Yuyang Zhou,
561
        Y. X. Zhu, Yanhong Xu, Yanping Huang, Yaohui Li, Yi Zheng, Yuchen Zhu, Yunxian Ma, Ying
562
        Tang, Yukun Zha, Yuting Yan, Z. Z. Ren, Zehui Ren, Zhangli Sha, Zhe Fu, Zhean Xu, Zhenda
563     Xie, Zhengyan Zhang, Zhewen Hao, Zhicheng Ma, Zhigang Yan, Zhiyu Wu, Zihui Gu, Zijia Zhu,
564     Zijun Liu, Zilin Li, Ziwei Xie, Ziyang Song, Zizheng Pan, Zhen Huang, Zhipeng Xu, Zhongyu
565     Zhang, and Zhen Zhang. Deepseek-r1: Incentivizing reasoning capability in llms via reinforce-
566     ment learning, 2025. URL https://arxiv.org/abs/2501.12948.
567
568
      Shihan Dou, Enyu Zhou, Yan Liu, Songyang Gao, Jun Zhao, Wei Shen, Yuhao Zhou, Zhiheng Xi,
        Xiao Wang, Xiaoran Fan, Shiliang Pu, Jiang Zhu, Rui Zheng, Tao Gui, Qi Zhang, and Xuanjing
569
        Huang. Loramoe: Alleviate world knowledge forgetting in large language models via moe-style
570
        plugin, 2024. URL https://arxiv.org/abs/2312.09979.
571
572   Runnan Fang, Yuan Liang, Xiaobin Wang, Jialong Wu, Shuofei Qiao, Pengjun Xie, Fei Huang,
573     Huajun Chen, and Ningyu Zhang. Memp: Exploring agent procedural memory. arXiv preprint
574     arXiv:2508.06433, 2025.
575
      Dayuan Fu, Keqing He, Yejie Wang, Wentao Hong, Zhuoma Gongque, Weihao Zeng, Wei Wang,
576
        Jingang Wang, Xunliang Cai, and Weiran Xu. Agentrefine: Enhancing agent generalization
577     through refinement tuning, 2025. URL https://arxiv.org/abs/2501.01702.
578
579   Yichao Fu, Peter Bailis, Ion Stoica, and Hao Zhang. Break the sequential dependency of llm infer-
580     ence using lookahead decoding, 2024. URL https://arxiv.org/abs/2402.02057.
581
      Huanang Gao, Jiayi Geng, Wenyue Hua, Mengkang Hu, Xinzhe Juan, Hongzhang Liu, Shilong
582     Liu, Jiahao Qiu, Xuan Qi, Yiran Wu, Hongru Wang, Han Xiao, Yuhang Zhou, Shaokun Zhang,
583     Jiayi Zhang, Jinyu Xiang, Yixiong Fang, Qiwen Zhao, Dongrui Liu, Qihan Ren, Cheng Qian,
584     Zhenhailong Wang, Minda Hu, Huazheng Wang, Qingyun Wu, Heng Ji, and Mengdi Wang. A
585     survey of self-evolving agents: On path to artificial super intelligence, 2025. URL https:
586     //arxiv.org/abs/2507.21046.
587
      In Gim, Guojun Chen, Seung-seob Lee, Nikhil Sarda, Anurag Khandelwal, and Lin Zhong. Prompt
588
         cache: Modular attention reuse for low-latency inference. Proceedings of Machine Learning and
589
         Systems, 6:325–338, 2024.
590
591   Raghavv Goel, Sudhanshu Agrawal, Mukul Gagrani, Junyoung Park, Yifan Zao, He Zhang, Tian
592     Liu, Yiping Yang, Xin Yuan, Jiuyan Lu, Chris Lott, and Mingu Lee. Vocabtrim: Vocabulary
593     pruning for efficient speculative decoding in llms, 2025. URL https://arxiv.org/abs/
        2506.22694.


                                                     11
      Under review as a conference paper at ICLR 2026




594
      Sachin Goyal, Ziwei Ji, Ankit Singh Rawat, Aditya Krishna Menon, Sanjiv Kumar, and Vaishnavh
595     Nagarajan. Think before you speak: Training language models with pause tokens, 2024. URL
596     https://arxiv.org/abs/2310.02226.
597
598   Shibo Hao, Sainbayar Sukhbaatar, DiJia Su, Xian Li, Zhiting Hu, Jason Weston, and Yuandong
599     Tian. Training large language models to reason in a continuous latent space, 2024. URL https:
600     //arxiv.org/abs/2412.06769.
601
      Dan Hendrycks, Collin Burns, Saurav Kadavath, Akul Arora, Steven Basart, Eric Tang, Dawn Song,
602
        and Jacob Steinhardt. Measuring mathematical problem solving with the math dataset, 2021.
603     URL https://arxiv.org/abs/2103.03874.
604
605   Hongkang Yang Hongkang Yang, Zehao Lin Zehao Lin, Wenjin Wang Wenjin Wang, Hao Wu
606     Hao Wu, Zhiyu Li Zhiyu Li, Bo Tang Bo Tang, Wenqiang Wei Wenqiang Wei, Jinbo Wang
607     Jinbo Wang, Zeyun Tang Zeyun Tang, Shichao Song Shichao Song, Chenyang Xi Chenyang Xi,
608     Yu Yu Yu Yu, Kai Chen Kai Chen, Feiyu Xiong Feiyu Xiong, Linpeng Tang Linpeng Tang, and
609     Weinan E Weinan E. Memory3 : Language modeling with explicit memory. Journal of Ma-
610
        chine Learning, 3(3):300–346, January 2024. ISSN 2790-203X. doi: 10.4208/jml.240708. URL
        http://dx.doi.org/10.4208/jml.240708.
611
612   Jian Hu, Jason Klein Liu, Haotian Xu, and Wei Shen. Reinforce++: An efficient rlhf algorithm with
613      robustness to both prompt and reward models, 2025a. URL https://arxiv.org/abs/
614      2501.03262.
615
616   Zhiyuan Hu, Yibo Wang, Hanze Dong, Yuhui Xu, Amrita Saha, Caiming Xiong, Bryan Hooi, and
617     Junnan Li. Beyond’aha!’: Toward systematic meta-abilities alignment in large reasoning models.
618
        arXiv preprint arXiv:2505.10554, 2025b.
619   HuggingFace. SmolLM3: smol, multilingual, long-context reasoner — huggingface.co. https:
620     //huggingface.co/blog/smollm3, 2025. [Accessed 23-09-2025].
621
622   Naman Jain, King Han, Alex Gu, Wen-Ding Li, Fanjia Yan, Tianjun Zhang, Sida Wang, Armando
623     Solar-Lezama, Koushik Sen, and Ion Stoica. Livecodebench: Holistic and contamination free
624     evaluation of large language models for code, 2024. URL https://arxiv.org/abs/
625     2403.07974.
626
      Chao Jin, Zili Zhang, Xuanlin Jiang, Fangyue Liu, Xin Liu, Xuanzhe Liu, and Xin Jin.
627     Ragcache: Efficient knowledge caching for retrieval-augmented generation. arXiv preprint
628     arXiv:2404.12457, 2024.
629
630   Mandar Joshi, Eunsol Choi, Daniel Weld, and Luke Zettlemoyer. TriviaQA: A large scale distantly
631    supervised challenge dataset for reading comprehension. In Regina Barzilay and Min-Yen Kan
632    (eds.), Proceedings of the 55th Annual Meeting of the Association for Computational Linguistics
633    (Volume 1: Long Papers), pp. 1601–1611, Vancouver, Canada, July 2017. Association for Com-
634
       putational Linguistics. doi: 10.18653/v1/P17-1147. URL https://aclanthology.org/
       P17-1147/.
635
636   Hengli Li, Chenxi Li, Tong Wu, Xuekai Zhu, Yuxuan Wang, Zhaoxin Yu, Eric Hanchen Jiang,
637     Song-Chun Zhu, Zixia Jia, Ying Nian Wu, and Zilong Zheng. Seek in the dark: Reasoning via
638     test-time instance-level policy gradient in latent space, 2025a. URL https://arxiv.org/
639     abs/2505.13308.
640
641   Yuhui Li, Fangyun Wei, Chao Zhang, and Hongyang Zhang. Eagle: Speculative sampling requires
642
        rethinking feature uncertainty, 2025b. URL https://arxiv.org/abs/2401.15077.
643   Wang Ling, Dani Yogatama, Chris Dyer, and Phil Blunsom. Program induction by rationale gener-
644    ation: Learning to solve and explain algebraic word problems. arXiv preprint arXiv:1705.04146,
645    2017.
646
647   Luyang Liu, Jonas Pfeiffer, Jiaxing Wu, Jun Xie, and Arthur Szlam. Deliberation in latent space via
        differentiable cache augmentation. arXiv preprint arXiv:2412.17747, 2024.


                                                      12
      Under review as a conference paper at ICLR 2026




648
      Hanjun Luo, Shenyu Dai, Chiming Ni, Xinfeng Li, Guibin Zhang, Kun Wang, Tongliang Liu, and
649     Hanan Salam. Agentauditor: Human-level safety and security evaluation for llm agents, 2025a.
650     URL https://arxiv.org/abs/2506.00641.
651
652   Michael Luo, Naman Jain, Jaskirat Singh, Sijun Tan, Ameen Patel, Qingyang Wu, Alpay Ariyak,
653    Colin Cai, Shang Zhu Tarun Venkat, Ben Athiwaratkun, Manan Roongta, Ce Zhang, Li Erran
654
       Li, Raluca Ada Popa, Koushik Sen, and Ion Stoica. Deepswe: Training a state-of-the-art
       coding agent from scratch by scaling rl. https://pretty-radio-b75.notion.site/
655
       DeepSWE-Training-a-Fully-Open-sourced-State-of-the-Art-Coding-Agent-by-Scaling-RL-
656
       2025b. Notion Blog.
657
658   Laurens van der Maaten and Geoffrey Hinton. Visualizing data using t-sne. Journal of machine
659     learning research, 9(Nov):2579–2605, 2008.
660
      Alex Mallen, Akari Asai, Victor Zhong, Rajarshi Das, Daniel Khashabi, and Hannaneh Hajishirzi.
661
        When not to trust language models: Investigating effectiveness of parametric and non-parametric
662     memories, 2023. URL https://arxiv.org/abs/2212.10511.
663
664   Dan P McAdams. The psychology of life stories. Review of general psychology, 5(2):100–122,
665     2001.
666
      Charles Packer, Sarah Wooders, Kevin Lin, Vivian Fang, Shishir G. Patil, Ion Stoica, and Joseph E.
667     Gonzalez. Memgpt: Towards llms as operating systems, 2024. URL https://arxiv.org/
668     abs/2310.08560.
669
670   Dmitrii Pantiukhin, Boris Shapkin, Ivan Kuznetsov, Antonia Anna Jost, and Nikolay Koldunov. Ac-
671    celerating earth science discovery via multi-agent llm systems, 2025. URL https://arxiv.
672
       org/abs/2503.05854.
673   Cheng Qian, Emre Can Acikgoz, Qi He, Hongru Wang, Xiusi Chen, Dilek Hakkani-Tür, Gokhan
674     Tur, and Heng Ji. Toolrl: Reward is all tool learning needs, 2025. URL https://arxiv.
675     org/abs/2504.13958.
676
677
      Jiahao Qiu, Xinzhe Juan, Yimin Wang, Ling Yang, Xuan Qi, Tongcheng Zhang, Jiacheng Guo, Yifu
         Lu, Zixin Yao, Hongru Wang, Shilong Liu, Xun Jiang, Liu Leqi, and Mengdi Wang. Agentdistill:
678
         Training-free agent distillation with generalizable mcp boxes, 2025a. URL https://arxiv.
679
         org/abs/2506.14728.
680
681   Jiahao Qiu, Xuan Qi, Tongcheng Zhang, Xinzhe Juan, Jiacheng Guo, Yifu Lu, Yimin Wang,
682      Zixin Yao, Qihan Ren, Xun Jiang, Xing Zhou, Dongrui Liu, Ling Yang, Yue Wu, Kaixuan
683      Huang, Shilong Liu, Hongru Wang, and Mengdi Wang. Alita: Generalist agent enabling scal-
684      able agentic reasoning with minimal predefinition and maximal self-evolution, 2025b. URL
685
         https://arxiv.org/abs/2505.20286.
686   David Rein, Betty Li Hou, Asa Cooper Stickland, Jackson Petty, Richard Yuanzhe Pang, Julien
687     Dirani, Julian Michael, and Samuel R. Bowman. Gpqa: A graduate-level google-proof q&a
688     benchmark, 2023. URL https://arxiv.org/abs/2311.12022.
689
690
      Shuo Ren, Pu Jian, Zhenjiang Ren, Chunlin Leng, Can Xie, and Jiajun Zhang. Towards scientific
        intelligence: A survey of llm-based scientific agents, 2025. URL https://arxiv.org/
691
        abs/2503.24047.
692
693   Zhenyi Shen, Hanqi Yan, Linhai Zhang, Zhanghao Hu, Yali Du, and Yulan He. Codi: Compressing
694     chain-of-thought into continuous space via self-distillation. arXiv preprint arXiv:2502.21074,
695     2025a.
696
      Zhenyi Shen, Hanqi Yan, Linhai Zhang, Zhanghao Hu, Yali Du, and Yulan He. Codi: Compressing
697
        chain-of-thought into continuous space via self-distillation, 2025b. URL https://arxiv.
698
        org/abs/2502.21074.
699
700   Mohit Shridhar, Xingdi Yuan, Marc-Alexandre Côté, Yonatan Bisk, Adam Trischler, and Matthew
701    Hausknecht. Alfworld: Aligning text and embodied environments for interactive learning, 2021.
       URL https://arxiv.org/abs/2010.03768.


                                                     13
      Under review as a conference paper at ICLR 2026




702
      Joykirat Singh, Raghav Magazine, Yash Pandya, and Akshay Nambi. Agentic reasoning and tool
703     integration for llms via reinforcement learning, 2025. URL https://arxiv.org/abs/
704     2505.01441.
705
706   Peiyang Song, Pengrui Han, and Noah Goodman. A survey on large language model reasoning
707     failures. In 2nd AI for Math Workshop @ ICML 2025, 2025. URL https://openreview.
708
        net/forum?id=hsgMn4KBFG.
709   Eleanor Spens and Neil Burgess. A generative model of memory construction and consolidation.
710     Nature human behaviour, 8(3):526–543, 2024.
711
712   Weihang Su, Yichen Tang, Qingyao Ai, Junxi Yan, Changyue Wang, Hongning Wang, Ziyi Ye,
713    Yujia Zhou, and Yiqun Liu. Parametric retrieval augmented generation, 2025. URL https:
714
       //arxiv.org/abs/2501.15915.
715   Yuchang Sun, Yanxi Chen, Yaliang Li, and Bolin Ding. Enhancing latent computation in transform-
716     ers with latent tokens, 2025. URL https://arxiv.org/abs/2505.12629.
717
718   Jihoon Tack, Jaehyung Kim, Eric Mitchell, Jinwoo Shin, Yee Whye Teh, and Jonathan Richard
719
         Schwarz. Online adaptation of language models with a memory of amortized contexts, 2024.
         URL https://arxiv.org/abs/2403.04317.
720
721   Wenhui Tan, Jiaze Li, Jianzhong Ju, Zhenbo Luo, Jian Luan, and Ruihua Song. Think silently, think
722    fast: Dynamic latent compression of llm reasoning chains, 2025. URL https://arxiv.org/
723    abs/2505.16552.
724
725
      Xiangru Tang, Tianrui Qin, Tianhao Peng, Ziyang Zhou, Daniel Shao, Tingting Du, Xinming Wei,
        Peng Xia, Fang Wu, He Zhu, Ge Zhang, Jiaheng Liu, Xingyao Wang, Sirui Hong, Chenglin Wu,
726
        Hao Cheng, Chi Wang, and Wangchunshu Zhou. Agent kb: Leveraging cross-domain experience
727
        for agentic problem solving, 2025. URL https://arxiv.org/abs/2507.06229.
728
729   James Thorne, Andreas Vlachos, Christos Christodoulopoulos, and Arpit Mittal. Fever: a large-scale
730     dataset for fact extraction and verification. arXiv preprint arXiv:1803.05355, 2018.
731
      Peng Wang, Zexi Li, Ningyu Zhang, Ziwen Xu, Yunzhi Yao, Yong Jiang, Pengjun Xie, Fei Huang,
732
        and Huajun Chen. Wise: Rethinking the knowledge memory for lifelong model editing of large
733     language models, 2024a. URL https://arxiv.org/abs/2405.14768.
734
735   Ruoyao Wang, Peter Jansen, Marc-Alexandre Côté, and Prithviraj Ammanabrolu. Scienceworld:
736     Is your agent smarter than a 5th grader?, 2022. URL https://arxiv.org/abs/2203.
737     07540.
738
      Xiaoqiang Wang, Suyuchen Wang, Yun Zhu, and Bang Liu. R3 Mem: Bridging memory reten-
739     tion and retrieval via reversible compression. In Wanxiang Che, Joyce Nabende, Ekaterina
740     Shutova, and Mohammad Taher Pilehvar (eds.), Findings of the Association for Computational
741     Linguistics: ACL 2025, pp. 4541–4557, Vienna, Austria, July 2025a. Association for Compu-
742     tational Linguistics. ISBN 979-8-89176-256-5. doi: 10.18653/v1/2025.findings-acl.235. URL
743     https://aclanthology.org/2025.findings-acl.235/.
744
      Yu Wang, Yifan Gao, Xiusi Chen, Haoming Jiang, Shiyang Li, Jingfeng Yang, Qingyu Yin, Zheng
745
        Li, Xian Li, Bing Yin, et al. Memoryllm: Towards self-updatable large language models. arXiv
746
        preprint arXiv:2402.04624, 2024b.
747
748   Yu Wang, Dmitry Krotov, Yuanzhe Hu, Yifan Gao, Wangchunshu Zhou, Julian McAuley, Dan Gut-
749     freund, Rogerio Feris, and Zexue He. M+: Extending memoryllm with scalable long-term mem-
750     ory, 2025b. URL https://arxiv.org/abs/2502.00592.
751
      Zhenhailong Wang, Haiyang Xu, Junyang Wang, Xi Zhang, Ming Yan, Ji Zhang, Fei Huang, and
752     Heng Ji. Mobile-agent-e: Self-evolving mobile assistant for complex tasks, 2025c. URL https:
753     //arxiv.org/abs/2501.11733.
754
755   Zora Zhiruo Wang, Jiayuan Mao, Daniel Fried, and Graham Neubig. Agent workflow memory,
        2024c. URL https://arxiv.org/abs/2409.07429.


                                                     14
      Under review as a conference paper at ICLR 2026




756
      Jason Wei, Xuezhi Wang, Dale Schuurmans, Maarten Bosma, Brian Ichter, Fei Xia, Ed Chi, Quoc
757     Le, and Denny Zhou. Chain-of-thought prompting elicits reasoning in large language models,
758     2023. URL https://arxiv.org/abs/2201.11903.
759
760   Yifan Wei, Xiaoyan Yu, Yixuan Weng, Tengfei Pan, Angsheng Li, and Li Du. Autotir: Autonomous
761     tools integrated reasoning via reinforcement learning, 2025. URL https://arxiv.org/
762     abs/2507.21836.
763
      Ronald J Williams. Simple statistical gradient-following algorithms for connectionist reinforcement
764     learning. Machine learning, 8:229–256, 1992.
765
766   Mingyuan Wu, Jingcheng Yang, Jize Jiang, Meitang Li, Kaizhuo Yan, Hanchao Yu, Minjia Zhang,
767    Chengxiang Zhai, and Klara Nahrstedt. Vtool-r1: Vlms learn to think with images via rein-
768    forcement learning on multimodal tool use, 2025a. URL https://arxiv.org/abs/2505.
769    19255.
770
      Yaxiong Wu, Sheng Liang, Chen Zhang, Yichao Wang, Yongyue Zhang, Huifeng Guo, Ruiming
771
        Tang, and Yong Liu. From human memory to ai memory: A survey on memory mechanisms in
772     the era of llms, 2025b. URL https://arxiv.org/abs/2504.15965.
773
774   Zidi Xiong, Yuping Lin, Wenya Xie, Pengfei He, Jiliang Tang, Himabindu Lakkaraju, and Zhen Xi-
775     ang. How memory management impacts llm agents: An empirical study of experience-following
776     behavior, 2025. URL https://arxiv.org/abs/2505.16067.
777
      Wujiang Xu, Kai Mei, Hang Gao, Juntao Tan, Zujie Liang, and Yongfeng Zhang. A-mem: Agentic
778
       memory for llm agents, 2025a. URL https://arxiv.org/abs/2502.12110.
779
780   Yige Xu, Xu Guo, Zhiwei Zeng, and Chunyan Miao. Softcot: Soft chain-of-thought for efficient
781     reasoning with llms, 2025b. URL https://arxiv.org/abs/2502.12134.
782
783
      Yige Xu, Xu Guo, Zhiwei Zeng, and Chunyan Miao. Softcot++: Test-time scaling with soft chain-
        of-thought reasoning, 2025c. URL https://arxiv.org/abs/2505.11484.
784
785   Zhangchen Xu, Yang Liu, Yueqin Yin, Mingyuan Zhou, and Radha Poovendran. Kodcode: A di-
786     verse, challenging, and verifiable synthetic dataset for coding, 2025d. URL https://arxiv.
787     org/abs/2503.02951.
788
789   Zifan Xu, Haozhu Wang, Dmitriy Bespalov, Xuan Wang, Peter Stone, and Yanjun Qi. Latent skill
790
        discovery for chain-of-thought reasoning. arXiv preprint arXiv:2312.04684, 2023.
791   An Yang, Baosong Yang, Beichen Zhang, Binyuan Hui, Bo Zheng, Bowen Yu, Chengyuan Li,
792     Dayiheng Liu, Fei Huang, Haoran Wei, et al. Qwen2. 5 technical report. arXiv preprint
793     arXiv:2412.15115, 2024a.
794
795   An Yang, Anfeng Li, Baosong Yang, Beichen Zhang, Binyuan Hui, and Bo Zheng et al. Qwen3
796     technical report, 2025. URL https://arxiv.org/abs/2505.09388.
797
      John Yang, Carlos E. Jimenez, Alexander Wettig, Kilian Lieret, Shunyu Yao, Karthik Narasimhan,
798     and Ofir Press. Swe-agent: Agent-computer interfaces enable automated software engineering,
799     2024b. URL https://arxiv.org/abs/2405.15793.
800
801   Yi Yang, Yixuan Tang, and Kar Yan Tam. Investlm: A large language model for investment using
802     financial domain instruction tuning, 2023. URL https://arxiv.org/abs/2309.13064.
803
      Weiran Yao, Shelby Heinecke, Juan Carlos Niebles, Zhiwei Liu, Yihao Feng, Le Xue, Rithesh
804
       Murthy, Zeyuan Chen, Jianguo Zhang, Devansh Arpit, Ran Xu, Phil Mui, Huan Wang, Caim-
805    ing Xiong, and Silvio Savarese. Retroformer: Retrospective large language agents with policy
806    gradient optimization, 2024. URL https://arxiv.org/abs/2308.02151.
807
808   Da Yin, Faeze Brahman, Abhilasha Ravichander, Khyathi Chandu, Kai-Wei Chang, Yejin Choi, and
809     Bill Yuchen Lin. Agent lumos: Unified and modular training for open-source language agents,
        2024. URL https://arxiv.org/abs/2311.05657.


                                                      15
      Under review as a conference paper at ICLR 2026




810
      Hongli Yu, Tinghong Chen, Jiangtao Feng, Jiangjie Chen, Weinan Dai, Qiying Yu, Ya-Qin Zhang,
811     Wei-Ying Ma, Jingjing Liu, Mingxuan Wang, and Hao Zhou. Memagent: Reshaping long-context
812     llm with multi-conv rl-based memory agent, 2025. URL https://arxiv.org/abs/2507.
813     02259.
814
815   Aohan Zeng, Mingdao Liu, Rui Lu, Bowen Wang, Xiao Liu, Yuxiao Dong, and Jie Tang. Agenttun-
816     ing: Enabling generalized agent abilities for llms, 2023. URL https://arxiv.org/abs/
817
        2310.12823.
818   Guibin Zhang, Muxin Fu, Guancheng Wan, Miao Yu, Kun Wang, and Shuicheng Yan. G-memory:
819     Tracing hierarchical memory for multi-agent systems, 2025a. URL https://arxiv.org/
820     abs/2506.07398.
821
822   Jianguo Zhang, Tian Lan, Rithesh Murthy, Zhiwei Liu, Weiran Yao, Ming Zhu, Juntao Tan, Thai
823      Hoang, Zuxin Liu, Liangwei Yang, Yihao Feng, Shirley Kokane, Tulika Awalgaonkar, Juan Car-
824
         los Niebles, Silvio Savarese, Shelby Heinecke, Huan Wang, and Caiming Xiong. Agento-
         hana: Design unified data and training pipeline for effective agent learning, 2024a. URL
825
         https://arxiv.org/abs/2402.15506.
826
827   Yang Zhang, Wenxin Xu, Xiaoyan Zhao, Wenjie Wang, Fuli Feng, Xiangnan He, and Tat-Seng
828     Chua. Reinforced latent reasoning for llm-based recommendation, 2025b. URL https://
829     arxiv.org/abs/2505.19092.
830
      Zeyu Zhang, Xiaohe Bo, Chen Ma, Rui Li, Xu Chen, Quanyu Dai, Jieming Zhu, Zhenhua Dong,
831
        and Ji-Rong Wen. A survey on the memory mechanism of large language model based agents,
832
        2024b. URL https://arxiv.org/abs/2404.13501.
833
834   Andrew Zhao, Daniel Huang, Quentin Xu, Matthieu Lin, Yong-Jin Liu, and Gao Huang. Expel: Llm
835     agents are experiential learners, 2024. URL https://arxiv.org/abs/2308.10144.
836
837
      Boyuan Zheng, Michael Y. Fatemi, Xiaolong Jin, Zora Zhiruo Wang, Apurva Gandhi, Yueqi Song,
        Yu Gu, Jayanth Srinivasa, Gaowen Liu, Graham Neubig, and Yu Su. Skillweaver: Web agents
838
        can self-improve by discovering and honing skills, 2025. URL https://arxiv.org/abs/
839
        2504.07079.
840
841   Wanjun Zhong, Lianghong Guo, Qiqi Gao, He Ye, and Yanlin Wang. Memorybank: Enhancing large
842    language models with long-term memory, 2023. URL https://arxiv.org/abs/2305.
843    10250.
844
      Zijian Zhou, Ao Qu, Zhaoxuan Wu, Sunghwan Kim, Alok Prakash, Daniela Rus, Jinhua Zhao,
845
        Bryan Kian Hsiang Low, and Paul Pu Liang. Mem1: Learning to synergize memory and reasoning
846     for efficient long-horizon agents, 2025. URL https://arxiv.org/abs/2506.15841.
847
848   Rui-Jie Zhu, Tianhao Peng, Tianhao Cheng, Xingwei Qu, Jinfa Huang, Dawei Zhu, Hao Wang,
849     Kaiwen Xue, Xuanliang Zhang, Yong Shan, et al. A survey on latent reasoning. arXiv preprint
850     arXiv:2507.06203, 2025.
851
852
853   A    U SE OF L ARGE L ANGUAGE M ODELS
854
855   In preparing this work, we employed large language models (LLMs) to assist with: polishing the
856
      writing, conducting literature reviews, and creating visualizations.
857
858   B    A DDITIONAL R ELATED W ORKS
859
860   LLM Decoding & RL. Two additional topics that relate to our work are LLM decoding and re-
861   inforcement learning (RL). From the decoding perspective, MemGen dynamically generates and in-
862   serts latent tokens, which shares similarity with speculative decoding where a drafter model receives
863   the current decoding context and produces subsequent drafted tokens (Cai et al., 2024; Fu et al.,
      2024; Li et al., 2025b; Goel et al., 2025). . However, these methods primarily aim to accelerate


                                                       16
      Under review as a conference paper at ICLR 2026




864
      LLM inference, whereas MemGen focuses on leveraging latent states as effective carriers of mem-
865   ory. From the RL perspective, MemGen employs rule-based RL to train the memory trigger, which
866   is closely related to reinforcement learning with variable reward (RLVR), including GRPO from
867   DeepSeek-R1 (DeepSeek-AI et al., 2025) and its various derivatives (Qian et al., 2025; Wu et al.,
868   2025a; Wei et al., 2025; Fu et al., 2025). While there exist efforts combining RL with agent memory,
869   to our knowledge, most do not address self-improving memory; for example, MemAgent (Yu et al.,
870   2025) and MEM1 (Zhou et al., 2025) focus on handling long-context inputs rather than evolving
871   memory mechanisms.
872
873   Latent Memory We further extend our discussion by incorporating related work on latent memory
874   mechanisms: R3Mem (Wang et al., 2025a) introduces a reversible, hierarchical memory architec-
875
      ture that integrates efficient long-term information retention with faithful retrieval, achieving state-
      of-the-art performance in long-context modeling, retrieval-augmented generation, and prolonged
876
      interactive tasks.
877
878
879   C     O PTIMIZATION A LGORITHM D ETAILS
880
881   In this section, we provide a detailed exposition of the optimization algorithms for training the
882   memory weaver, Wweaver , as mentioned in §4.3. The core principle, as established in Equation (10),
883
      is to update only the weaver’s parameters, denoted as θ′ , while keeping the reasoner πθ frozen. This
      modularity allows MemGen to be compatible with various optimization paradigms. We detail the
884
      specific implementations for Supervised Fine-Tuning (SFT) and a reinforcement learning approach,
885
      Group Relative Policy Optimization (GRPO).
886
887
      C.1    C OMBINING MemGen WITH SFT
888
889   The objective of Supervised Fine-Tuning is to train the memory weaver to generate latent mem-
890   ories that guide the frozen reasoner πθ to replicate the behavior observed in a dataset of high-
891   quality demonstration trajectories. We leverage the provided history of past experiences, H =
892   {(xi , τi∗ )}N                 ∗
                   i=1 , where each τi is treated as an expert demonstration.
893   Each expert trajectory τi∗ consists of a sequence of states and actions, τi∗ = (si,0 , a∗i,0 , si,1 , a∗i,1 , . . .).
894   Each expert action a∗i,t is a sequence of tokens, a∗i,t = (z∗i,t,1 , z∗i,t,2 , . . . , z∗i,t,Lt ). The goal is to
895   maximize the conditional log-likelihood of generating these expert token sequences.
896
897
      During the training of the weaver, both the reasoner πθ and the memory trigger T are held fixed. At
      each token generation step j where the trigger activates, the weaver Wθ′ takes the reasoner’s hidden
898
      states Hi,t,<j as input and produces a latent memory Mi,t . The reasoner then generates the next
899
      token conditioned on this memory. The SFT objective is to adjust the weaver’s parameters θ′ to
900   maximize the probability of the ground-truth token z∗i,t,j .
901
902   Formally, the optimization problem is to minimize the negative log-likelihood of the expert trajec-
903
      tories, averaged over the dataset H. The loss function for the weaver’s parameters θ′ is defined
      as:
904                                                                                                  
                                                  i −1 X
                                                 TX    Lt
905
                    LSFT (θ′ ) = −E(xi ,τi∗ )∼H          log πθ (z∗i,t,j | si,t , z∗i,t,<j , Mi,t,j ) , (11)
906
                                                      t=0 j=1
907
      where the latent memory Mi,t,j is synthesized by the weaver at that specific step:
908
909                                               Mi,t,j = Wθ′ (Hi,t,<j ).                                            (12)
910   Note that the generation of Mi,t,j only occurs if the fixed trigger T determines an invocation is
911   needed at step j. In steps where no memory is invoked, the conditioning term Mi,t,j is omitted. The
912   gradients are computed exclusively with respect to the weaver’s parameters θ′ and used to update
913   them via gradient descent:
914                                         θ′ ← θ′ − η∇θ′ LSFT (θ′ ),                                (13)
915   where η is the learning rate. Through this process, the memory weaver learns to synthesize latent
916   memories that effectively steer the frozen reasoner’s generative process to align with the training
917   data’s behavior, thereby internalizing the knowledge from the demonstration data without corrupting
      the general capabilities of the core LLM.


                                                               17
      Under review as a conference paper at ICLR 2026




918
      C.2    C OMBING MemGen WITH GRPO
919
920   The memory weaver can also be trained using a reinforcement learning objective. We specifically
921   adapt the GRPO algorithm. The training process begins by sampling a batch of tasks from the ex-
                                                                         W ,T
922   perience history H. For each task xi , we use the policy Πθ θ′ (composed of the frozen reasoner
923   πθ and the current memory weaver Wθ′ ) to generate a group of K distinct trajectories, denoted as
924   Gi = {τi,1 , τi,2 , . . . , τi,K }. Each trajectory is generated by the agent’s interaction with the environ-
925   ment and results in a final reward R(τi,k ), evaluated by the reward function. Going forward, GRPO
926   computes a group-relative baseline by averaging the rewards of all trajectories within the group Gi :
927                                                                    K
                                                             1 X
928                                                R̄(Gi ) =     R(τi,k ).                                          (14)
929                                                          K
                                                                   k=1
930
      The advantage for a specific trajectory τi,k is then its reward relative to this baseline:
931
932                                              A(τi,k ) = R(τi,k ) − R̄(Gi ).                                     (15)
933
      This formulation allows the weaver to learn by differentiating between better and worse outcomes
934
      within a set of its own generations, promoting policies that produce trajectories with above-average
935   rewards. The final objective function, maximized with respect to the weaver’s parameters θ′ , is:
936                                       "                                                                                     #
                                                 K
937            ′                              1 X                 Wθ′ ,T
                                                                                            
                                                                                              Wθ′ ,T
                                                                                                                               
      JGRPO (θ ) = E             W ′ ,T             A(τi,k ) log Πθ      (τi,k | xi ) − β KL Πθ      (· | xi ) ∥ Πref (· | xi ) ,
938                  xi ∼H, Gi ∼Πθ θ          K
                                                k=1
939                                                                                                                 (16)
940   where Πref is a fixed reference policy and β is a coefficient controlling the strength of the KL regu-
941   larization. The gradients are computed only for the weaver’s parameters θ′ , thus refining its ability
942   to synthesize impactful latent memories while preserving the integrity of the core reasoner.
943
944   C.3    M EMORY T RIGGER T RAINING D ETAILS
945                                                                     ϕ
946
      We now detail the training procedure of the memory trigger Ttrigger  . During this stage, the reasoner
947
      πθ and the memory weaver Wθ′ remain frozen, and only the trigger parameters ϕ are updated. For
      each task xi with initial environment state s0 , we generate trajectories under the composed policy
948                                ϕ
949
      induced by πθ , Wθ′ , and Ttrigger . At each decoding step t and token position j where zt,j ∈ D, the
950
      trigger observes the hidden-state prefix Ht,<j and samples a binary action
951                                                                             ϕ
                                                                                              
                    dt,j ∼ πϕ (d | Ht,<j ) := Bernoulli pt,j ,       pt,j = σ Ttrigger (Ht,<j ) ,       (17)
952
953   where dt,j = 1 corresponds to INVOKE. If invoked, the decoder pauses and the weaver synthesizes
954   a latent memory Mt from Ht,<j , after which generation resumes based on the augmented context.
955   Let Ii denote the set of activated positions in trajectory τi , and let d̃i = {dt,j }(t,j)∈Ii be the trigger
956   decisions within that rollout. After completing the episode, the environment returns a scalar reward
957   R(τi ), which we assign uniformly to all decisions in d̃i .
958   To discourage excessively frequent invocation, we introduce a sparsity-inducing penalty. For a group
959   of rollouts generated from the same task xi , we identify the high-reward subset
960
961                                    Hhigh = {i : R(τi ) ≥ mediank R(τk )},                                       (18)
962   and compute the reference activation level
963                                                        X 1 X
                                                    1
964                                       p̄ =                     dt,j ,                                           (19)
965                                               |Hhigh |   |Ii |
                                                             i∈Hhigh         (t,j)∈Ii
966
      which estimates the typical invocation frequency among successful trajectories. The trigger is
967
      trained to maximize the reward while penalizing activations above p̄:
968
                                               h            X                     i
969                     max Eτi ∼πθ , d̃i ∼T ϕ R(τi ) − λ         max 0, dt,j − p̄ ,           (20)
970                          ϕ                     trigger
                                                                           (t,j)∈Ii
971
      where λ > 0 controls the sparsity regularization strength.


                                                                 18
       Under review as a conference paper at ICLR 2026




972
       We optimize this objective following a GRPO-style procedure. For each task prompt, we generate a
973                              (g)
       group of G trajectories {τi }G g=1 (we set G = 8 in practice; Table 3). We compute per-trajectory
974
                    (g)
975    rewards R(τi ) and construct normalized advantages
976
977
                                                                                  v
978                  (g)                                  G                       u G
              (g) R(τi ) − µi                         1
                                                          X        (g)
                                                                                  u X     (g)
979          Ai =             ,                  µi = G         R(τi ),      σi = t 1
                                                                                    G (R(τ ) − µi )2 ,
                                                                                              i
980                   σi                                  g=1                           g=1
981
982
983    thus measuring each rollout’s performance relative to the group. The trigger parameters are then
984    updated by maximizing a GRPO-style surrogate objective:
985
986
987                             G
988                           1 X       X                  (g)    (g)  (g)                
                      max                          log πϕ dt,j | Ht,<j Ai − β KL πϕ ∥ πϕold .
989                       ϕ   G g=1          (g)
                                      (t,j)∈Ii
990
991
992
993    After each update, we set πϕold ← πϕ and repeat. This objective allows the trigger to learn to invoke
994
       memory only at critical points that improve group-relative performance while maintaining sparsity.
995
996
997
998    D     E XPERIMENTAL D ETAILS
999
1000
1001   D.1   T RAINING DATASET S ETUP
1002
1003
       For all datasets except PopQA, we utilize their official training sets to train both the memory weaver
1004
       and memory trigger. Specifically, we first train the memory weaver in isolation. Since the trigger
1005
       module is not yet available at this stage, latent memory tokens are randomly inserted at punctuation
1006   marks during training to provide diverse learning signals. Once the memory weaver is adequately
1007   trained, we proceed to train the memory trigger, enabling it to learn when to invoke memory (with
1008   the memory weaver fixed). For PopQA, which lacks a dedicated training set, we instead leverage the
1009   trained model on TriviaQA and directly evaluate on PopQA. The dataset statistics are summarized
1010   in Table 2.
1011
1012
1013                                               Table 2: Dataset Statistics.
1014
1015
                    Dataset Name          Train Size        Test Size               Note
1016
1017                ALFWorld                3.32K             134         From AgentBank-ALFWorld
1018
                    TriviaQA                4.13K            7.9K         From AgentBank-TriviaQA
                    PopQA                      -             14.3K                   -
1019
                    KodCode                   8K              2K            KodCode-Light-RL-10K
1020
                    BigCodeBench             912              228                    -
1021                GPQA                     448              198          GPQA-Diamond for testing
1022                GSM8K                   7.47K             1K                     -
1023                MATH                     1.6K             500                    -
1024
1025



                                                                 19
       Under review as a conference paper at ICLR 2026




1026
       D.2               PARAMETER C ONFIGURATIONS
1027
1028
1029
                                     Table 3: Hyperparameters used in the training of MemGen.
1030
1031
                                       Settings            Hyperparameters
1032
1033                                                       train batch size = 4
1034                                                       learning rate = 1e-5
                                                           epochs = 2
1035                                   Training(SFT)
                                                           warmup ratio = 0.1
1036
                                                           optim = adamw torch
1037
                                                           schedular = cosine
1038
1039                                                       rollout batch size = 8
1040
                                                           train batch size = 8
                                                           epochs = 2
1041
                                                           beta = 0.0
1042
                                       Training(GRPO)      num iterations = 1
1043                                                       learning rate = 1e-5
1044                                                       warmup ratio = 0.1
1045                                                       optim = adamw torch
1046                                                       schedular = cosine
1047
                                                           r = 16
1048                                                       lora alpha = 32
1049                                                       target modules = [q proj, v proj]
1050
                                       LoRA
                                                           lora dropout = 0.1
1051                                                       bias = none
1052                                                       task type = CAUSAL LM
1053                                                       adam offload
1054                                   Optimization        flash attn
1055                                                       deepspeed enable sleep
1056
1057
1058
1059   E               E XTRA R ESULTS
1060
1061
1062   Table 4: Results on Qwen2.5-1.5B. All values represent the performance metric for each task (e.g.,
1063   accuracy %). We highlight the best and second best results.
1064
1065
                      Method         ALFWorld TrivialQA PopQA KodCode BigCodeBench GPQA GSM8K MATH
1066                  Vanilla          22.54      32.10   16.08        24.55    40.35      11.62   39.51   36.63
                      CoT              18.30      28.67   18.39        32.32    38.59      15.67   56.79   45.22
1067
1068                  SFT              36.57      63.84   39.20        55.83    37.72      11.11   54.83   38.84
                      GRPO             43.55      68.21   43.15        62.11    70.34      15.65   68.10   47.42




       Qwen2.5-1.5B
1069
                      REINFORCE        43.25      66.50   41.87        60.20    67.80      12.50   67.40   46.89
1070
                      REINFORCE++      43.66      66.90   44.69        63.33    69.50      13.80   69.04   47.33
1071                  Agent-FLAN       35.80      64.28   38.90        56.21    43.83       9.35   53.02   29.82
1072
                      ExpeL            28.96      25.20   20.20        31.15    39.78       8.12   45.12   38.12
1073                  MemoryBank       27.89      38.14   22.78        37.93    35.87      13.87   47.88   30.47
1074                  AWM              30.42      55.69   32.54          -        -           -      -       -
1075                  SoftCoT          33.07      62.22   38.78        55.13    36.10       9.31   54.50   38.55
1076                  Co-processor     35.66      64.78   40.12        56.65    38.10      12.12   57.12   37.40
1077                  MemGen SFT       40.30      65.02   41.28        58.16    42.47      18.28   58.15   47.12
1078                  MemGen GRPO      54.27      73.42   49.28        65.43    72.81      18.18   73.39   53.36
1079



                                                                  20
       Under review as a conference paper at ICLR 2026




1080
       E.1     C ONTINUAL L EARNING R ESULT
1081
1082   The results in Table 5 indicate three main findings. First, MemGen exhibits stronger knowledge re-
1083   tention and forward transfer than SFT and ExpeL. For example, when trained on GPQA, MemGen
1084   reaches 47.96% on GSM8K and 28.80% on KodCode, surpassing SFT at 45.74% and 18.50% by
1085   margins of +2.22% and +10.30%, respectively. Similarly, when trained on KodCode, MemGen
1086   maintains 40.34% on AQuA and 20.09% on GPQA, whereas SFT yields 28.61% and 2.53%. Sec-
1087   ond, MemGen demonstrates more balanced cross-task generalization. Unlike SFT, which primarily
1088
       improves performance on the most recent task (e.g., 54.10% on KodCode but only 2.53% on GPQA),
       MemGen achieves competitive results across tasks. After training on GSM8K, it attains 38.43% on
1089
       AQuA and 21.72% on GPQA, both higher than SFT and ExpeL, suggesting that latent memory
1090
       captures task-invariant reasoning. Third, MemGen effectively mitigates forgetting on earlier tasks.
1091   After sequential training on KodCode, it preserves 40.34% on AQuA compared to 27.14% of ExpeL
1092   and 28.61% of SFT, showing stronger robustness to catastrophic forgetting. These findings suggest
1093   that MemGen provides a more stable and transferable mechanism for continual learning.
1094
1095   E.2     T RIGGER F REQUENCY V ISUALIZATION
1096
1097   More visualizations of the trigger frequency are displayed in Figures 7 and 8, where we paired
1098   Qwen2.5-1.5B or SmolLM3-3B with the GSM8K dataset and tested on subsets of GSM8K,
1099   KodCode, and GPQA (each subset having the same number of samples). We then tallied the fre-
1100   quency of memory trigger INVOKE judgments at each relative percentile position in the LLM output.
1101
1102
                                           120                                              GSM8K (Avg: 75.17)
1103                                                                                        GPQA (Avg: 46.87)




                 Memory Invocation Count
1104                                       100                                              KodCode (Avg: 51.70)
1105
1106                                       80
1107                                       60
1108
1109                                       40
1110
                                           20
1111
1112                                        0
1113                                             0   20         40              60              80                 100
1114                                                      Relative Insertion Position (%)
1115
1116
       Figure 7: Memory invocation frequency across benchmarks at inference (trained on MemGen
       SFT +Qwen2.5-1.5B +GSM8K).
1117
1118
1119
       E.3     F RAMEWORK A NALYSIS
1120
1121   E.3.1    A BLATION S TUDY
1122
1123   To assess the effectiveness of the proposed memory trigger, we conduct an ablation study compar-
1124   ing different memory invocation strategies, as summarized in Table 6. Specifically, Random denotes
1125
       a naı̈ve baseline where latent memory tokens are inserted at arbitrary token positions with a fixed
       probability p. All delimiters activated represents a sentence-level strategy that invokes the memory
1126
       weaver at every delimiter position without any selection mechanism. Finally, MemGen’s dedicated
1127
       Trigger corresponds to our standard approach, where the memory weaver is activated by the trained
1128   trigger Ttrigger . The results reveal several key observations. First, sentence-level intervention already
1129   improves performance compared to random invocation. For instance, activating the weaver at all
1130   delimiters yields 17.34%, 56.20%, and 64.15% on GPQA, Kodcode, and TriviaQA, respectively,
1131   consistently outperforming all random baselines (e.g., p = 0.5 achieves only 16.66%, 52.95%, and
1132   57.28%). This highlights the importance of aligning memory injection with semantic boundaries
1133   rather than distributing it across the sequence. More importantly, our trained trigger achieves the
       best overall performance, reaching 18.28%, 58.16%, and 65.02% on the three benchmarks. This


                                                                     21
       Under review as a conference paper at ICLR 2026




1134                                      160
1135                                                                                            GSM8K (Avg: 74.49)
                                          140                                                   GPQA (Avg: 48.59)




                Memory Invocation Count
1136
1137                                      120                                                   KodCode (Avg: 27.96)
1138                                      100
1139
                                          80
1140
1141                                      60
1142                                      40
1143                                      20
1144
1145
                                           0
                                                0       20          40              60               80                100
1146                                                          Relative Insertion Position (%)
1147
1148   Figure 8: Memory invocation frequency across benchmarks at inference (trained on MemGen
1149   SFT +SmolLM3-3B +GSM8K).
1150
1151
1152
                                                                 Trained on GSM8K
1153
1154
                                  GSM8K
1155
1156                     KodCode
1157
1158
1159
                                           GPQA                                                           Vanilla
                                                                                                          SFT
1160
                                          MATH                                                            MemoryBank
1161
                                                                                                          MemGen
1162
1163                                                0    10      20           30         40          50          60
1164
1165   Figure 9: The generalization study of MemGen. We train MemGen SFT on GSM8K and evaluate it on
1166   all four datasets.
1167
1168                                                            Trained on KodCode
1169
                                                                                                          Vanilla
1170                              GSM8K                                                                   SFT
1171
                                                                                                          MemoryBank
1172
                         KodCode                                                                          MemGen
1173
1174
1175                                       GPQA
1176
1177
1178                                      MATH
1179
1180                                                0    10      20           30         40          50          60
1181
1182   Figure 10: The generalization study of MemGen. We train MemGen SFT on KodCode and evaluate it
1183   on all four datasets.
1184
1185
1186   demonstrates that selective activation, i.e., deciding when and where to weave in memory, provides
1187   the most effective support for reasoning, as it balances memory utility and interference more pre-
       cisely than coarse-grained alternatives.


                                                                         22
       Under review as a conference paper at ICLR 2026




1188   Table 5: Continual learning results of Qwen2.5-1.5B-Instruct across four reasoning and
1189   programming datasets (AQuA, GPQA, GSM8K, KodCode). The model is sequentially trained on
1190   each dataset (AQuA → GPQA → GSM8K → KodCode), and after each training stage, evaluation
1191   is conducted on all four benchmarks.
1192
1193                                 Method         AQuA      GPQA          GSM8K     KodCode
                     Trained On
1194
                                     Vanilla        41.34        11.62       39.51      24.55
1195
1196                 AQuA            SFT            42.52        16.67       42.10      18.20
1197                                 ExpeL          41.73        12.67       40.16      16.30
1198                                 MemGen SFT     43.31        19.70       39.80      19.55
1199
1200                 GPQA            SFT            38.55        17.17       45.74      18.50
1201                                 ExpeL          37.24        14.35       42.67      15.20
1202                                 MemGen SFT     39.85        20.72       47.96      28.80
1203
                     GSM8K           SFT            33.46        13.13       52.31      19.45
1204
                                     ExpeL          34.89        12.42       48.78      13.65
1205
                                     MemGen SFT     38.43        21.72       55.67      19.75
1206
1207                 KodCode         SFT            28.61         2.53       24.14      54.10
1208                                 ExpeL          27.14         6.23       31.44      48.35
1209                                 MemGen SFT     40.34        20.09       53.72      52.95
1210
1211
1212
       Table 6: Ablation study of different memory invocation strategies. Random denotes a naı̈ve baseline
       where latent memory tokens are inserted at arbitrary token positions with a fixed probability p. All
1213
       delimiters activated represents a sentence-level strategy that invokes the memory weaver at every de-
1214
       limiter position without any selection mechanism. Finally, MemGen’s dedicated Trigger corresponds
1215   to our standard approach, where the memory weaver is activated by the trained trigger Ttrigger .
1216
1217
                       Memory Invocation Strategy       GPQA             Kodcode     TriviaQA
1218
1219                   Random (p = 0.2)                  15.66            54.55       63.55
1220                   Random (p = 0.5)                  16.66            52.95       57.28
1221                   Random (p = 0.8)                  12.63            53.60       62.22
1222                   All delimiters activated          17.34            56.20       64.15
1223                   MemGen’s dedicated Trigger        18.28            58.16       65.02
1224
1225
1226   E.3.2   A NALYSIS OF M EMORY W EAVER
1227
1228   We further investigate how the parameterization of the memory weaver influences downstream per-
1229   formance. As shown in Table 7, increasing the number of trainable parameters by adopting a full-
1230
       parameter SFT setup surely enhances the weaver’s memory and learning capacity, leading to higher
       task performance across benchmarks. Nevertheless, the LoRA-based instantiation already achieves
1231
       competitive results while retaining high parameter efficiency, demonstrating that even lightweight
1232
       adaptations can endow the weaver with sufficient capacity to generate effective latent memories.
1233
1234
       E.3.3   E FFICIENCY A NALYSIS
1235
1236   Table 8 reports average inference time and task performance across three benchmarks. Both SFT
1237   and MemGen drastically reduce per-task inference time compared with vanilla models due to fewer
1238   tokens required to reach correct answers. For instance, on KodCode with Qwen2.5-1.5B, MemGen
1239   SFT completes tasks in 2.94 s, a 75.4% reduction from 11.96 s for vanilla, while improving accuracy
1240   by 33.61%. On ALFWorld+Qwen3-8B, MemGen SFT adds only 1.6% more time compared with
1241   SFT (20.08 s vs 19.76 s) but increases accuracy from 83.59% to 85.82%. These results confirm that
       MemGen efficiently enhances reasoning performance without incurring significant inference delays.


                                                        23
       Under review as a conference paper at ICLR 2026




1242
1243
1244
1245
1246
1247
       Table 7: Ablation study of the latent weaver. We instantiate MemGen SFT ’s memory weaver with
1248   LoRA of different ranks as well as full-parameter SFT.
1249
1250
                      Base LLM: Qwen2.5-1.5B-Instruct         GPQA       Kodcode    TriviaQA
1251
1252                  LoRA (r = 4)                            17.16       54.85      63.04
1253                  LoRA (r = 6)                            18.18       55.25      64.55
1254                  LoRA (r = 8)                            17.67       55.75      64.10
1255                  LoRA (r = 16)                           18.28       56.16      65.02
1256                  LoRA (r = 24)                           18.67       57.20      65.40
1257
                      LoRA (r = 32)                           19.26       57.95      65.85
1258
1259                  Full SFT                                21.21       60.00      67.10
1260
1261
1262
1263
1264
1265
1266
1267
1268
1269
1270
1271   Table 8: Average per-task inference time (seconds) and task performance (accuracy %) across three
1272   benchmarks. Performance improvement (%) of MemGen SFT over vanilla LLM is also reported.
1273
1274         Model & Method                 KodCode               ALFWorld              TriviaQA
1275                                    Time (s)    Acc (%)   Time (s)    Acc (%)   Time (s)   Acc (%)
1276
                                                   Qwen2.5-1.5B
1277
1278
             Vanilla                     11.96       24.55      21.17      22.54     2.18      32.10
             SFT                          2.01       55.83      10.79      36.57     1.98      63.84
1279         MemGen SFT                  2.94       58.16       12.94     40.30      2.05      65.02
1280         Improvement over Vanilla       -       +33.61        -       +17.76       -       +32.92
1281
                                                    SmolLM-3B
1282
1283
             Vanilla                     13.12       37.05      34.82      18.96     4.26      10.47
             SFT                          3.04       59.25      12.88      32.36     3.05      55.25
1284         MemGen SFT                  3.48       62.65       14.69     50.60      3.16      68.13
1285         Improvement over Vanilla       -       +25.60        -       +31.64       -       +57.66
1286
                                                     Qwen3-8B
1287
1288
             Vanilla                     16.99       49.10      55.42      58.93     8.70      52.18
             SFT                          7.24       64.75      19.76      83.59     5.98      74.55
1289         MemGen SFT                  7.56       66.15       20.08     85.82      6.25      77.22
1290         Improvement over Vanilla       -       +17.05        -       +26.89       -       +25.04
1291
1292
1293
1294
1295



                                                        24
       Under review as a conference paper at ICLR 2026




1296
       E.3.4      A BLATION OF M EMORY W EAVER : PAUSE T OKEN
1297
1298   To further assess the necessity of the memory weaver, we introduce a strong and conceptually
1299   aligned baseline based on the pause token mechanism (Goyal et al., 2024). Pause tokens are a
1300   small set of learned special tokens P = {<pause>1 , . . . , <pause>K } that instruct the model to
1301   suspend outward generation while continuing to update its internal hidden states. Formally, inserting
1302
       a pause token <pause> at position j forces the model to compute a new hidden state
1303                                         hj = fθ (hj−1 , <pause>),
1304
       while suppressing any semantic output. A sequence of K pause tokens thus yields K steps of
1305
       latent computation over the current cognitive state. We train a set of K pause tokens (where K
1306   exactly equals the number of generated latent tokens by memory weaver) under the same protocol
1307   in Equation (10).
1308
1309   Results.     We report the ablation results as follows:
1310
1311   Table 9: The performance comparison between MemGen and a variant by replacing mem-
1312   ory weaver with learnt pause tokens on three benchmarks.  The backbone is set as
1313   Qwen2.5-1.5B-Instruct.
1314
1315                         Method                      Kodcode      TriviaQA      GPQA
1316                         Trigger + Pause Token         49.50          56.30      13.80
1317                         Trigger + Weaver              58.16          65.02      18.28
1318
1319   As shown in Table 9, although pause tokens provide moderate gains, indicating that allowing the
1320   model brief intervals of internal processing is beneficial, they consistently underperform the full
1321   weaver-equipped system. We attribute this to the fixed, context-independent nature of pause tokens:
1322   they cannot reconstruct or integrate task-relevant information in a targeted manner. In contrast, the
1323   memory weaver produces context-sensitive latent memory vectors tailored to the reasoner’s current
1324   cognitive state, leading to substantially stronger performance.
1325
1326   F     I NTEGRATION WITH R ETRIEVAL - BASED M EMORY
1327
1328   F.1     F ORMALIZING THE I NTEGRATION P ROCESS
1329
       While the primary mechanism of MemGen leverages the parametric knowledge encapsulated within
1330
       the memory weaver Wweaver , the framework is designed to be extensible, allowing for seamless
1331   integration with external, retrieval-based memory systems. This hybrid approach enables the weaver
1332   to synthesize latent memories that are informed by both its internalized experiences and a vast corpus
1333   of external information, thereby providing a richer and more comprehensive context to the reasoner
1334   πθ . Let Mext denote an external memory database, and let R(·) be a retrieval function that, given a
1335   natural language query, returns a set of relevant textual memory snippets. When the memory trigger
1336   Ttrigger determines an invocation is necessary at the token-generation step j (i.e., dj = INVOKE), the
1337   natural language text generated thus far serves as the query for the external memory system. This
1338   query, denoted as qt,j , is produced by decoding the sequence of tokens generated up to that point:
1339                                            qt,j = Decode(zt,<j ).                                        (21)
1340
       The retrieval process is then formalized as:
1341
1342                                              Ct = R(qt,j ; Mext ),                                       (22)
1343   where Ct = {c1 , c2 , . . . , cP } is a set of P retrieved textual snippets. These snippets are subsequently
1344   encoded into a sequence of embeddings, Et ∈ RLc ×dmodel , where Lc is the total length of the encoded
1345   text. This allows the weaver to process the retrieved information in its native latent space.
1346
       This retrieved information is subsequently merged with the reasoner’s internal cognitive state Ht,<j .
1347
       The combined context is then fed into the memory weaver Wweaver to produce the final latent mem-
1348
       ory. The invocation of the weaver, as described in Equation (5), is thus modified to:
1349
                                             Mt = Wweaver ([Ht,<j ; Et ]),                                    (23)


                                                           25
       Under review as a conference paper at ICLR 2026




1350
       where [·; ·] denotes the concatenation of the hidden state sequences. This integrated process allows
1351   Wweaver to reconstruct both internal parametric knowledge and externally retrieved information into
1352   a compact, potent latent memory Mt for the reasoner.
1353
1354
       F.2     E XPERIMENTAL R ESULTS
1355
1356   As shown in Table 10, even when MemGen’s own parametric memory is disabled (which means
1357
       that merely the retrieved textual snippets are fed into Wweaver ), MemGen significantly enhances the
       retrieval baseline, boosting performance on ALFWorld from 36.18% to 45.60% and on PopQA
1358
       from 28.16% to 39.50%. This demonstrates that MemGen serves as a powerful synthesizer, not
1359
       merely appending but actively reconstructing externally retrieved information into a more potent,
1360   task-relevant latent memory, thereby underscoring the necessity of a generative memory mechanism.
1361   Most notably, when the weaver’s internal parametric memory is enabled alongside external retrieval,
1362   performance is dramatically amplified, achieving scores of 75.90% on ALFWorld and 76.40% on
1363   TriviaQA. This highlights the powerful synergy achieved when the memory weaver is empowered
1364   to dynamically integrate its own internalized experiences with external textual knowledge.
1365
1366   Table 10: Performance results when integrating MemGen with a retrieval-based memory baseline,
1367   ExpeL. The backbone LLM is set as SmolLM3-3B.
1368
                     Method                                   ALFWorld     TriviaQA    PopQA
1369
                     Vanilla LLM                                18.96        10.47       8.23
1370
                     ExpeL                                      36.18        46.20      28.16
1371                 MemGen + ExpeL (w/o parametric memory)     45.60        53.20      39.50
1372                 MemGen + ExpeL (w/ parametric memory)      75.90        76.40      60.23
1373
1374
1375
       G      L ATENT M EMORY T OKEN D EMONSTRATION
1376
1377
             Cases from TriviaQA
1378
1379
             Case 1
1380         system prompt
1381         Answer the given question. You must conduct reasoning inside <think> and </think> first
                   every time you get new information. After reasoning, if you find you lack some
1382              knowledge, you can call a search engine by <search> query </search> and it will
1383              return the top searched results between <information> and </information>. You can
                  search as many times as your want. If you find no further external knowledge
1384              needed, you can directly provide the answer inside <answer> and </answer>, without
1385               detailed illustrations. For example, <answer> Beijing </answer>.

1386         user prompt
1387         Which Lloyd Webber musical premiered in the US on 10th December 1993?

1388         assistant output
1389         [UPPORT...’,eniable certif]<think> I need to find out which Lloyd Webber musical
                  premiered in the US on 10th December 1993. </think>
1390         [essengeryyyyMMddELCOME certif]<search> Lloyd Webber musical premiered in the US on 10
1391              th December 1993 </search>

1392         user prompt
1393         <observation> Doc 1(Title: "Andrew Lloyd Webber") opened at the Adelphi Theatre in
                  London on 12 July 1993, and ran for 1,529 performances ... </observation>
1394
1395         assistant output
             [essengeryyyyMMddELCOME
1396              certif]<think> The search results indicate that "Sunset Boulevard" is the Lloyd
1397              Webber musical that premiered in the US on 10th December 1993. </think>
             [UPPORT...’,eniable certif]<answer> Sunset Boulevard </answer>
1398
1399         Case 2
             system prompt
1400         Answer the given question. You must conduct reasoning inside <think> and </think> first
1401               every time you get new information. After reasoning, if you find you lack some
                  knowledge, you can call a search engine by <search> query </search> and it will
1402              return the top searched results between <information> and </information>. You can
1403              search as many times as your want. If you find no further external knowledge




                                                        26
       Under review as a conference paper at ICLR 2026




1404
1405           needed, you can directly provide the answer inside <answer> and </answer>, without
                detailed illustrations. For example, <answer> Beijing </answer>.
1406
1407      user prompt
          Who had a 70s No 1 hit with Kiss You All Over?
1408
1409      assistant output
          [UPPORT...’,eniable certif]<think> I need to search who had a 70s No 1 hit with the
1410           song "Kiss You All Over." </think>
1411      [essengeryyyyMMddELCOME
               certif]<search> who had a 70s No 1 hit with Kiss You All Over? </search>
1412
1413      user prompt
          <observation> Doc 1(Title: "Kiss You All Over") Exile’s only big hit in the pop rock
1414           market (they would later have great success on the country music charts)... </
1415           observation>

1416      assistant output
1417      [essengerUPPORT’options certif]<think> The song "Kiss You All Over" was a 70s No 1 hit
               with the band Exile. </think>
1418      [UPPORT...’,eniable certif]<answer> Exile </answer>
1419
          Case 3
1420      system prompt
1421      Answer the given question. You must conduct reasoning inside <think> and </think> first
                every time you get new information. After reasoning, if you find you lack some
1422           knowledge, you can call a search engine by <search> query </search> and it will
1423           return the top searched results between <information> and </information>. You can
               search as many times as your want. If you find no further external knowledge
1424           needed, you can directly provide the answer inside <answer> and </answer>, without
1425            detailed illustrations. For example, <answer> Beijing </answer>.

1426      user prompt
1427      What was the name of Michael Jackson’s autobiography written in 1988?

1428      assistant output
1429      [UPPORTUPPORTeniable certif]<think> I need to search for the name of Michael Jackson’s
               autobiography written in 1988. </think>
1430      [UPPORT...’,eniable certif]<search> Michael Jackson autobiography 1988 </search>
1431
          user prompt
1432      <observation> Doc 1(Title: "Michael Jackson") of 4.4 million people. In 1988, Jackson
1433           released his only autobiography, ""Moonwalk"", which took four years to complete
               and sold 200,000 copies... </observation>
1434
1435      assistant output
          [cmsUPPORT’options Certif]<think> The name of Michael Jackson’s autobiography written
1436           in 1988 is "Moonwalk". </think>
1437      [UPPORTUPPORTeniable certif]<answer> Moonwalk </answer>

1438      Case 4
1439      system prompt
          Answer the given question. You must conduct reasoning inside <think> and </think> first
1440            every time you get new information. After reasoning, if you find you lack some
1441           knowledge, you can call a search engine by <search> query </search> and it will
               return the top searched results between <information> and </information>. You can
1442           search as many times as your want. If you find no further external knowledge
1443           needed, you can directly provide the answer inside <answer> and </answer>, without
                detailed illustrations. For example, <answer> Beijing </answer>.
1444
1445      user prompt
          Which volcano in Tanzania is the highest mountain in Africa?
1446
1447      assistant output
          [UPPORT...’,eniable certif]<think> I need to search for the highest mountain in Africa,
1448            which is located in Tanzania. </think>
1449      [UPPORTUPPORTeniable certif]<search> highest mountain in Africa in Tanzania </search>

1450      user prompt
1451      <observation> Doc 1(Title: "Mount Kilimanjaro") Mount Kilimanjaro Mount Kilimanjaro or
               just Kilimanjaro (), with its three volcanic cones, ""Kibo"", ""Mawenzi"", and ""
1452           Shira"", is a dormant volcano in Tanzania... </observation>
1453
          assistant output
1454      [essengeryyyyMMddeniable certif]<think> The highest mountain in Africa located in
1455           Tanzania is Mount Kilimanjaro. </think>
          [UPPORT...’,eniable certif]<answer> Mount Kilimanjaro </answer>
1456
1457



                                                    27
       Under review as a conference paper at ICLR 2026




1458
          Cases from KodCode
1459
1460      Case 1
1461      MODEL INPUT:
          Write an efficient and correct Python function to solve the following problem.Question:
1462            Write a function that takes a list of integers and returns a new list containing
1463           only the even integers from the original list, sorted in non-decreasing order. The
                list can contain negative numbers and is not necessarily sorted. If the input
1464           list is empty or contains no even integers, the function should return an empty
1465           list.

1466      MODEL OUTPUT:
1467      [keyword-kindërgetAs-slide]def even_sorted(lst):
          [.keyword pick] """
1468      [LTRetical] Returns a new list containing only the even integers from the original list
1469           ,
             sorted in non-decreasing order.
1470      [.keyword pick] """
1471      [LTRetical] even_numbers = [num for num in lst if num % 2 == 0]
          [.keyword pick] even_numbers.sort()
1472      [.keywordetical] return even_numbers<|im_end|>
1473
          Case 2
1474      MODEL INPUT:
1475      Write an efficient and correct Python function to solve the following problem.Question:
                Write a python function to find the first maximum and the second maximum number
1476           in an array. If the array has less than two unique elements, return ‘None‘ for the
1477            missing values.

1478      MODEL OUTPUT:
1479      [keyword-kindërgetAs-slide]def find_first_max_second_max(arr):
             """
1480      [.keywordetical] Returns the first and second maximum numbers in the array.
1481         If the array has less than two unique elements,[.keyword_pick] returns None for the
                  missing values.
1482      [.keyword pick] """
1483         if len(arr) < 2:
          [.keyword pick] return None,[.keyword pick] None
1484
1485      [.keywordetical] unique_elements = set(arr)
          [.keyword pick] if len(unique_elements) < 2:
1486      [.keywordetical] return None,[.keyword pick] None
1487
          [.keywordetical] first_max = max(unique_elements)
1488         second_max = max(unique_elements, default=None)
1489
             return first_max, second_max<|im_end|>
1490
1491      Case 3
          MODEL INPUT:
1492      Write an efficient and correct Python function to solve the following problem.Question:
1493            Given a list of integers, write a function that returns the length of the longest
                contiguous subarray where all elements have the same value. The function should
1494           also return the value of the elements in this subarray. Consider the output to be
1495           a tuple where the first element is the length of the longest subarray and the
               second element is the value of the elements in the subarray.
1496
1497      For example, consider the following list as input,
          numbers = [1, 2, 2, 3, 3, 3, 2, 2, 2, 2, 1, 1]
1498
1499      Expected Output: (4, 2)

1500      MODEL OUTPUT:
1501      [kindgetAs?-kindërgetAs-slide]def longest_subarray_with_same_value(numbers):
             """
1502      [.keyword pick] Returns the length of the longest contiguous subarray where all
1503           elements have the same value.
          [.keyword pick] The function also returns the value of the elements in this subarray.
1504      [LTR pick] """
1505      [LTRetical] if not numbers:
          [.keyword pick] return (0, None)
1506
1507         max_length = 1
          [.keyword pick] current_length = 1
1508         max_value = numbers[0]
1509         current_value = numbers[0]

1510      [.keyword pick] for i in range(1,[LTR pick] len(numbers)):
1511      [LTR pick] if numbers[i] == current_value:




                                                      28
       Under review as a conference paper at ICLR 2026




1512
1513               current_length += 1
                else:
1514               if current_length > max_length:
1515      [LTR pick] max_length = current_length
                      max_value = current_value
1516               current_length = 1
1517               current_value = numbers[i]

1518         if current_length > max_length:
1519            max_length = current_length
                max_value = current_value
1520
1521         return (max_length, max_value)<|im_end|>

1522
1523
1524      Cases from GSM8K
1525
          Case 1
1526      MODEL INPUT:
1527      Solve the math problem with proper reasoning, and make sure to put the FINAL ANSWER
               inside \boxed{}.
1528      Question: Billy sells DVDs. He has 8 customers on Tuesday. His first 3 customers buy
1529           one DVD each. His next 2 customers buy 2 DVDs each. His last 3 customers don’t buy
                any DVDs. How many DVDs did Billy sell on Tuesday?
1530
1531      MODEL OUTPUT:
          [ even. s. even even ...
1532           even] Billy’s first 3 customers buy 3 * 1 = <<3*1=3>>3 DVDs.
1533      [-gallery.]His next 2 customers buy 2 * 2 = <<2*2=4>>4 DVDs.
          [-routing.]Billy sells 3 + 4 = <<3+4=7>>7 DVDs on Tuesday.\boxed{7}<|im_end|>
1534
1535      Case 2
          MODEL INPUT:
1536      Solve the math problem with proper reasoning, and make sure to put the FINAL ANSWER
1537           inside \boxed{}.
          Question: Henry made two stops during his 60-mile bike trip. He first stopped after 20
1538           miles. His second stop was 15 miles before the end of the trip. How many miles did
1539            he travel between his first and second stops?

1540      MODEL OUTPUT:
1541      [..... a even even even... even
               ]Henry traveled 60-15=<<60-15=45>>45 miles before his second stop.
1542      [-android.]He traveled 20 miles before his first stop,[ecies.] so he traveled
1543           45-20=<<45-20=25>>25 miles between his first and second stops.\boxed{25}<|im_end|>
          [ even.]The final answer is 25<|im_end|>
1544
1545      Case 3
          MODEL INPUT:
1546      Solve the math problem with proper reasoning, and make sure to put the FINAL ANSWER
1547           inside \boxed{}.Question: A basket contains 25 oranges among which 1 is bad, 20%
               are unripe, 2 are sour and the rest are good. How many oranges are good?
1548
1549      MODEL OUTPUT:
          [. even.. pre. even even even... even
1550           ]The number of unripe oranges is 20/100*25 = <<20/100*25=5>>5
1551      [-animate.]The number of sour oranges is 2
          [-login.]The number of bad oranges is 1
1552      [-login.]The number of good oranges is 25 - 5 - 2 - 1 = <<25-5-2-1=17>>17\boxed{17}<|
1553           im_end|>

1554
1555
1556      Cases from GSM8K Cluster
1557
          Cluster 1
1558      Solve the math problem with proper reasoning, and make sure to put the FINAL ANSWER
1559           inside \boxed{}.Question: Josh decides to try flipping a house. He buys a house
               for $80,000 and then puts in $50,000 in repairs. This increased the value of the
1560           house by 150%. How much profit did he make?
1561      [..... pre. even even even... even ]The house was worth $80

1562      Solve the math problem with proper reasoning, and make sure to put the FINAL ANSWER
1563           inside \boxed{}.Question: James decides to run 3 sprints 3 times a week. He runs
               60 meters each sprint. How many total meters does he run a week?
1564      [..... out. even even even.. even even ]He runs 3*60=<<3*60=180>>180 meters a day
1565



                                                     29
       Under review as a conference paper at ICLR 2026




1566
1567       Solve the math problem with proper reasoning, and make sure to put the FINAL ANSWER
                inside \boxed{}.Question: Kylar went to the store to buy glasses for his new
1568            apartment. One glass costs $5, but every second glass costs only 60% of the price.
1569             Kylar wants to buy 16 glasses. How much does he need to pay for them?
           [.... pre... even... even ]The first 16 glasses cost 16 * 5 = $<<16*5=80>>80
1570
1571       Solve the math problem with proper reasoning, and make sure to put the FINAL ANSWER
                inside \boxed{}.Question: Toulouse has twice as many sheep as Charleston.
1572            Charleston has 4 times as many sheep as Seattle. How many sheep do Toulouse,
1573            Charleston, and Seattle have together if Seattle has 20 sheep?
           [. even even. s. even even even.. even even
1574            ]Charleston has 4*20=<<4*20=80>>80 sheep.
1575
           Solve the math problem with proper reasoning, and make sure to put the FINAL ANSWER
1576            inside \boxed{}.Question: John drives for 3 hours at a speed of 60 mph and then
1577            turns around because he realizes he forgot something very important at home. He
                tries to get home in 4 hours but spends the first 2 hours in standstill traffic.
1578            He spends the next half-hour driving at a speed of 30mph, before being able to
1579            drive the remaining time of the 4 hours going at 80 mph. How far is he from home
                at the end of those 4 hours?
1580       [...... even even even... even
1581            ]He drives 60*3 = <<60*3=180>>180 miles in the first 3 hours.

1582       Cluster 2
1583       Janet’s ducks lay 16 eggs per day,[-animate.] and she eats 3 eggs for breakfast every
                morning
1584
1585       Wendi feeds her 20 chickens 15 cups of feed in the morning,[-animate.] 25 cups in the
                afternoon
1586
1587       She sold half of what was left at the orange house,[-animate.] she had 1/2 left.

1588       x = 162/18 = <<162/18=9>>9
1589       [-animate.]Allen’s age is 11x = 11*9 = <<11*9=99>>99 years old now.

1590       He runs 20 miles a day because 60/3=<<60/3=20>>20
1591       [-animate.]He runs 10 miles the other two days because 20/2=<<20/2=10>>10

1592       Cluster 3
1593       It takes 2/2=<<2/2=1>>1 bolt of white fiber.
           [ecies.]So it takes 2+1=<<2+1=3>>3 bolts in total.
1594
1595       The total cost of the five cell phones is $150 x 5 = $<<150*5=750>>750.
           [ even.]The interest for each unit is $150 x 2/100 = $<<150*2/100=3>>3.
1596
1597       The second level is 16*2=<<16*2=32>>32 square footage.
           [ even.]The third level is 32*2=<<32*2=64>>64 square footage.
1598
1599       20% of 20 students is 20 x 20/100 = <<20*20/100=4>>4 students.
           [ecies.]The remaining students are 20 - 4 = <<20-4=16>>16 students.
1600
1601       If Artie rounds the price of petunias to the nearest dollar,[-animateecies] he will
                sell them for $2 each.
1602
1603       Cluster 4
           The groomer offers a 30% discount,[-loginecies] so the cost of the grooming is 100 -
1604            (100 * 0.[ecies.]30) = $<<100-30=70>>70
1605
           If Dana skips at 3 miles per hour,[-login.] then she runs at 3 * 4 = <<3*4=12>>12 miles
1606             per hour.
1607
           There were 140 * 2 = <<140*2=280>>280 emeralds in the chest.
1608       [-login.]There were 175 + 140 + 280 = <<175+140+280=595>>595 gems in the chest.
1609
           There are 30+60 = <<30+60=90>>90 students at Dr.[-loginecies.] Wertz’s school.
1610
1611       Cecilia will feed her puppy 1 cup of dog food every day for 180 days,[-login.] so she
                will use 180 cups of dog food in the first 180 days.
1612
1613
1614
1615   H     M EMORY F UNCTIONAL S TUDY
1616
1617   H.1   V ISUALIZATION P ROCESS
1618
1619   This section details the methodology employed to visualize the geometric structure of the latent
       memory (as shown in Figures 5 and 11) and to identify semantically distinct clusters within it. The


                                                       30
       Under review as a conference paper at ICLR 2026




1620                                      Latent Memory Visualization
1621                                        (across all benchmarks)
1622                                                     BCB
1623                                        GPQA

1624
                                                           KodCode
1625                                       MATH
1626
1627
                                         GSM8K                        ALFWorld
1628                                                 TriviaQA
1629
1630
1631
1632
                   Latent Memory Visualization                       Latent Memory Visualization
1633
                          (on TriviaQA)                                      (on GSM8K)
1634
1635
1636
1637
1638
1639
1640
1641
1642
1643
1644
1645
1646   Figure 11: (Up) t-SNE visualization of latent memories generated by MemGen +SmolLM3-3B
1647   across datasets; (Down) Latent memory visualization within the TriviaQA and GSM8K datasets,
1648   clustered using K-means.
1649
1650
1651   process involves two primary stages: obtaining a fixed-dimensional representation for each memory
1652   sequence and applying dimensionality reduction and clustering algorithms.
1653
1654   Sequence Representation. Let Dmem = {M1 , M2 , . . . , MQ } be a corpus of Q latent memory
1655
       sequences collected from MemGen’s inference, where each sequence Mi = (mi,1 , . . . , mi,K ) ∈
1656   RK×dmodel consists of K token embeddings of dimension dmodel . To obtain a single, holistic vector
1657
       representation for each sequence, we compute its mean embedding m̄i ∈ Rdmodel :
1658                                                           K
                                                         1 X
1659                                             m̄i =       mi,l .                                     (24)
1660
                                                         K
                                                               l=1
1661
       This procedure yields a set of high-dimensional representations M̄ = {m̄1 , m̄2 , . . . , m̄Q }, which
1662
       serves as the input for the subsequent analysis.
1663
1664   For visualization, we employed t-SNE (Maaten & Hinton, 2008) to project the high-dimensional
1665   set M̄ into a two-dimensional space. Formally, this mapping can be expressed as yi ∈ R2 =
1666   ft-SNE (m̄i ). These 2D data points are then utilized for the scatter plot.
1667   To quantitatively identify distinct functional groups, we applied the K-means algorithm directly to
1668   the high-dimensional mean embeddings in M̄. This procedure partitions the memory representa-
1669   tions into a predefined number of N discrete clusters, C = {C1 , C2 , . . . , CN } (we set N = 4).
1670
1671   H.2   FAILURE TAXONOMY D EFINITIONS
1672
1673   To systematically analyze the functional impact of different memory clusters, we adopt and refine
       the failure taxonomy proposed by Song et al. (2025). Below, we provide precise definitions of


                                                          31
       Under review as a conference paper at ICLR 2026




1674
       each failure category in the context of LLM-based agent behavior, which guide our annotation and
1675   evaluation process:
1676
1677   • Planning Failure. This category refers to errors in high-level task decomposition and strategic
1678     decision-making. The agent either formulates subgoals that do not align with the original objec-
1679     tive, fails to plan intermediate steps necessary for task completion, or misorders the reasoning
1680     sequence, leading to suboptimal or incomplete solutions.
1681   • Compositional Reasoning. This type of failure occurs when the agent struggles to integrate
1682     multiple pieces of information or reasoning steps into a coherent solution.
1683   • Tool Parsing Error. This failure occurs when the agent’s generated output cannot be parsed
1684     into a valid tool call. Typical patterns include malformed function signatures, missing required
1685     arguments, or unrequired arguments in the invocation, which prevent the external tool from being
1686     executed as intended.
1687   • Tool Response Error. This category refers to errors that arise after a tool has been successfully
1688     invoked. Either the tool itself returns incorrect or incomplete information, or the agent misuses or
1689     misinterprets the returned content.
1690   • Answer Formatting Failure. This category includes errors in the final presentation or structuring
1691     of the output, despite the reasoning process being largely correct. Examples include incorrect
1692     output format (e.g., unable to include the answers with in \boxed), violation of task-specific
1693     response templates, or missing required components in the final answer.
1694   • Demand Misunderstanding. This failure indicates that the agent has misinterpreted the user’s
1695     intent or the task specification. It may solve a different problem from the one posed, ignore key
1696     constraints, or pursue irrelevant objectives due to misunderstanding the instruction semantics.
1697   • Think-Act Inconsistency. This type refers to a mismatch between the agent’s internal reasoning
1698     and its subsequent action in ReAct mode. The agent may articulate a correct reasoning chain but
1699     execute a contradictory action or produce a final answer inconsistent with its prior deliberation.
1700   • False Belief captures cases where the agent maintains and reasons with inaccurate assumptions
1701     about the external environment, user state, or task context.
1702
1703   H.3    A NNOTATING FAILURE M ODES AND F ILTERING L ATENT M EMORY
1704
1705   Annotation of Failure Modes. We manually annotated agent failures on the TriviaQA dataset
1706   based on the eight failure modes discussed above. Each trajectory generated by the agent was evalu-
1707   ated by human annotators who assigned only one failure label if the agent’s behavior deviated from
1708   a successful path. Trajectories exhibiting failures that did not fall into these predefined categories
1709
       were excluded from this specific analysis.
1710
1711
       Inference-time Filtering of Latent Memory Clusters. The core of our intervention study in-
       volved selectively removing the influence of a specific latent memory cluster during the agent’s
1712
       reasoning process. Our methodology is designed to be consistent with the clustering process itself,
1713
       which operates on sequence-level representations.
1714
1715   Let C = {C1 , C2 , . . . , CN } be the set of N latent memory clusters. These clusters were derived by
1716   applying K-means to a collection of historical memory representations, where each representation
1717
       is the mean embedding of an entire latent memory sequence. For each cluster Ci ∈ C, we compute
       its centroid µi ∈ Rdmodel by averaging these historical sequence representations within the cluster.
1718
1719   During inference, when the memory weaver synthesizes a new latent memory sequence Mt =
1720   (mt,1 , . . . , mt,K ) ∈ RK×dmodel , we first compute its single-vector representation, m̄new , by averag-
1721   ing its constituent token embeddings:
1722                                                           K
                                                          1 X
1723                                            m̄new =       mt,l .                                        (25)
                                                          K
1724                                                           l=1
1725   To determine the semantic affiliation of this new sequence, we compare its mean embedding m̄new
1726   against a comprehensive reference set Ecomp = Evocab ∪ {µ1 , . . . , µN }, where Evocab ∈ RV ×dmodel
1727   is the LLM’s vocabulary embedding matrix. We then identify the set of top-k nearest neighbors to
       m̄new based on cosine similarity, denoted as Sk (m̄new ). In our experiments, we set k = 10.


                                                          32
       Under review as a conference paper at ICLR 2026




1728
       To ablate the influence of a target cluster Cj , the entire latent memory sequence Mt is filtered (i.e.,
1729   discarded and not prepended to the reasoner’s context) if the centroid of that cluster, µj , is present
1730   within this top-k set. Formally, Mt is filtered if:
1731
                                                     µj ∈ Sk (m̄new ).                                     (26)
1732
1733   This sequence-level filtering allows us to precisely ablate the contribution of a specific learned mem-
1734   ory function and observe its impact on agent behavior.
1735
1736   H.4   C OMPARISON B ETWEEN L ATENT T OKENS AND T RAJECTORIES
1737
       For each memory invocation at step j and timestep t, we construct three sets of embeddings within
1738
       the hidden representation space of the frozen reasoner πθ . First, we record the latent memory tokens
1739   generated by the weaver,
1740
                                       Mt = [mt,1 , . . . , mt,K ] ∈ RK×dmodel .
1741
       Second, we obtain the corresponding context hidden states that serve as the hook for memory gen-
1742
       eration,
1743                                Ht,<j = (ht,1 , . . . , ht,j−1 ) ∈ R(j−1)×dmodel ,
1744
       collected directly from the forward pass of πθ during online reasoning. Third, we take a past success-
1745
       ful trajectory τprev from H and replay it under teacher forcing through πθ , yielding its hidden-state
1746   sequence
1747                                    H̃traj = (h̃1 , . . . , h̃L ) ∈ RL×dmodel .
1748
1749   To enable joint visualization, we compute the mean embedding of each set:
1750                             K                                j−1                      L
                              1 X                           1 X                          1X
1751                  m̄t =       mt,k ,        h̄ctx
                                                  t =           ht,k ,        h̄traj =      h̃ℓ .
1752                          K                            j−1                           L
                                k=1                               k=1                      ℓ=1
1753   Finally, we project the three averaged vectors
1754
1755
                                                     {m̄t , h̄ctx   traj
                                                              t , h̄ }

1756   from Rdmodel into two dimensions (via t-SNE) for visualization. The resulting visualizations are
1757   shown in Figure 12.
1758
1759
1760
1761
1762
1763
1764
1765
1766
1767
1768
1769
       Figure 12: The t-SNE visualizations among the learnt latent tokens, trajectory representations and
1770
       context representations. The left is drawn from KodCode dataset, and the right is from TriviaQA
1771
       dataset.
1772
1773   As shown in Figure 12, the trajectory and context embeddings occupy more similar regions of the
1774   feature space, as both originate from transformations of human-readable tokens. In contrast, latent
1775   memory tokens form a more machine-native and human-unreadable representation, resulting in a
1776   substantially different embedding distribution.
1777
       We further quantify the geometric relationship between the three clusters by computing the Eu-
1778
       clidean distances between their centers in Rdmodel . For each invocation (t, j), we treat m̄t , h̄ctx
                                                                                                         t , and
1779   h̄traj as the centers of the latent-memory, context, and trajectory clusters, respectively, and compute
1780   the pairwise distances
1781            (t)                            (t)                           (t)
               dmem-ctx = m̄t − h̄ctx
                                  t   2
                                        ,    dmem-traj = m̄t − h̄traj 2 ,   dctx-traj = h̄ctx
                                                                                          t − h̄
                                                                                                traj
                                                                                                     2
                                                                                                       .


                                                             33
       Under review as a conference paper at ICLR 2026




1782
       Aggregating over all invocations T , we report the mean distances
1783
1784                     1 X (t)                         1 X (t)                         1 X (t)
            d¯mem-ctx =         dmem-ctx , d¯mem-traj =        dmem-traj , d¯ctx-traj =      dctx-traj ,
1785                    |T |                            |T |                            |T |
                             t∈T                              t∈T                           t∈T
1786
1787   which directly reflect how far the latent memory centers deviate from the context and replayed-
1788
       trajectory manifolds.
1789
       Table 11: Pairwise Euclidean distances between the centers of context, latent-memory (inserted),
1790
       and trajectory clusters.
1791
1792
                 Dataset      context ↔ inserted     context ↔ trajectory     inserted ↔ trajectory
1793
1794             GPQA               165.311                   64.180                 185.315
1795
                 KodCode            227.108                   96.831                 197.105
1796
1797   From Table 11, the distances reveal a consistent geometric pattern: the context and trajectory centers
1798   lie extremely close, while the generated latent memories occupy a distant region of the representa-
1799   tion space, forming a well-separated cluster from both context and trajectory. This large separation
1800
       demonstrates that the weaver does not merely compress or replay past hidden states; instead, it syn-
       thesizes novel latent structures that are parametrically reconstructed rather than retrieved, supporting
1801
       our claim that latent memory introduces genuinely new inferential content beyond the observed tra-
1802
       jectories.
1803
1804
1805
1806
1807
1808
1809
1810
1811
1812
1813
1814
1815
1816
1817
1818
1819
1820
1821
1822
1823
1824
1825
1826
1827
1828
1829
1830
1831
1832
1833
1834
1835



                                                         34
