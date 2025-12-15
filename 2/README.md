## Testing

### Native grep
time grep ERROR bigfile.txt > grep.out

### Distributed mygrep (3 nodes, quorum=2)
time cat bigfile.txt | ./mygrep grep ERROR \
--nodes localhost:9001,localhost:9002,localhost:9003 \
--quorum 2 > mygrep.out

diff grep.out mygrep.out