Output simple-usage-example.gif

Set TypingSpeed 100ms
Set Height 1000
Set Width 1000
Require slotalk

Type "tree"
Enter
Sleep 2.5s
Type "cat app1/metrics.go"
Enter
Sleep 2s
Type "slotalk init -f app1/metrics.go"
Enter
Sleep 2s
Type "cat app1/metrics.go | slotalk init -f -"
Enter
Sleep 2s
Type "slotalk init --to-file"
Enter
Sleep 1s
Type "tree"
Enter
Sleep 3.5s
Type "cat slo_definitions/app1.yaml"
Enter
Sleep 500ms
Type "cat slo_definitions/app1.yaml"
Enter
Sleep 500ms
Type "cat slo_definitions/app1.yaml"
Enter
Sleep 500ms
Sleep 3s
Ctrl+D

