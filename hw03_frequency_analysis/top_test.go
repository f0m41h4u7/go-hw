package hw03_frequency_analysis //nolint:golint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Change to true if needed
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

var textEnglish = `BPF is a highly flexible and efficient virtual machine-like construct in the Linux kernel allowing to execute bytecode at various hook points in a safe manner. It is used in a number of Linux kernel subsystems, most prominently networking, tracing and security (e.g. sandboxing).
	Although BPF exists since 1992, this document covers the extended Berkeley Packet Filter (eBPF)
version which has first appeared in Kernel 3.18 and renders the original version which is being
referred to as “classic” BPF (cBPF) these days mostly obsolete. cBPF is known to many as being the
packet filter language used by tcpdump. Nowadays, the Linux kernel runs eBPF only and loaded
cBPF bytecode is transparently translated into an eBPF representation in the kernel before program
execution.
	BPF does not define itself by only providing its instruction set, but also by offering further
infrastructure around it such as maps which act as efficient key-value stores, helper functions to
interact with and leverage kernel functionality, tail calls for calling into other BPF programs, security hardening primitives, a pseudo file system for pinning objects (maps, programs), and infrastructure for allowing BPF to be offloaded, for example, to a network card. LLVM provides a BPF back end, so that tools like clang can be used to compile C into a BPF object file, which can then be loaded into the kernel. BPF is deeply tied to the Linux kernel and allows for full  programmability without sacrificing native kernel performance. 
	Last but not least, also the kernel subsystems making use of BPF are part of infrastructure. The two main subsystems discussed throughout this document are tc and XDP where BPF programs can be adached to. XDP BPF programs are adached at the earliest networking driver stage and trigger a run of the BPF program upon packet reception. By definition, this achieves the best possible packet processing performance since packets cannot get processed at an even earlier point in soeware. However, since this processing occurs so early in the networking stack, the stack has not yet extracted metadata out of the packet. On the other hand, tc BPF programs are executed later in the kernel stack, so they have access to more metadata and core kernel functionality. Apart from tc and XDP programs, there are various other kernel subsystems as well which use BPF such as tracing (kprobes, uprobes, tracepoints, etc).`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		assert.Len(t, Top10(""), 0)
	})

	t.Run("test on english text", func(t *testing.T) {
		expected := []string{"the", "bpf", "kernel", "and", "to", "a", "in", "as", "is", "programs"}
		assert.Subset(t, expected, Top10(textEnglish))
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{"он", "а", "и", "что", "ты", "не", "если", "то", "его", "кристофер", "робин", "в"}
			assert.Subset(t, expected, Top10(text))
		} else {
			expected := []string{"он", "и", "а", "что", "ты", "не", "если", "-", "то", "Кристофер"}
			assert.ElementsMatch(t, expected, Top10(text))
		}
	})
}
