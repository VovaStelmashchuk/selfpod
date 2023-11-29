package app

import (
	"reflect"
	"testing"
)

func TestSimpleCaseExtractPTagsBeforeBR(t *testing.T) {
	htmlContent := `<p>First Paragraph</p><p>Second Paragraph</p><br><p>Third Paragraph</p>`
	expected := []string{"First Paragraph", "Second Paragraph"}

	result, err := extractPTagsBeforeBR(htmlContent)
	if err != nil {
		t.Errorf("extractPTagsBeforeBR returned an unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("extractPTagsBeforeBR returned %v, expected %v", result, expected)
	}
}

func TestRealCaseExtractPTagsBeforeBR(t *testing.T) {
	htmlContent := `<p>00:00:00 - Вступ. Вова знайшов нову роботу. Як Вова ходив по&nbsp;співбесідах.</p><p>00:30:31 - Фінансова подушка і збереження.</p><p>00:39:06 - Поради про те, як не проводити співбесіду.</p><p>00:44:11 - Говоримо про перший&nbsp;доповідь&nbsp;на Dou Mobile Meetup.</p><p>00:52:41 - Говоримо про BetterMe та їх&nbsp;доповідь. І трохи про Wix.</p><br><p>Коментарі та побажання можна залишити в нашому <a href="https://t.me/androidstory_chat" rel="noopener noreferrer" target="_blank">телеграм чаті</a>.</p><p><a href="https://www.patreon.com/androidstory" rel="noopener noreferrer" target="_blank">Support the show</a></p><p><br></p> <a target='_blank' rel='noopener noreferrer' href="https://open.acast.com/public/patreon/fanSubscribe/4793863">Тут для вас є ще більше нашого контенту</a><br /><hr><p style='color:grey; font-size:0.75em;'> Hosted on Acast. See <a style='color:grey;' target='_blank' rel='noopener noreferrer' href='https://acast.com/privacy'>acast.com/privacy</a> for more information.</p>`

	expected := []string{
		"00:00:00 - Вступ. Вова знайшов нову роботу. Як Вова ходив по співбесідах.",
		"00:30:31 - Фінансова подушка і збереження.",
		"00:39:06 - Поради про те, як не проводити співбесіду.",
		"00:44:11 - Говоримо про перший доповідь на Dou Mobile Meetup.",
		"00:52:41 - Говоримо про BetterMe та їх доповідь. І трохи про Wix.",
	}

	result, err := extractPTagsBeforeBR(htmlContent)
	if err != nil {
		t.Errorf("extractPTagsBeforeBR returned an unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("extractPTagsBeforeBR returned %v, expected %v", result, expected)
	}
}
