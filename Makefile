.PHONY: clean merc namenode director

clean:
	rm -rf DataNode
merc:
	rm -rf DataNode
	mkdir DataNode
	go run Mercenarios/Mercenarios.go

namenode:
	rm -rf DataNode
	mkdir DataNode
	go run NameNode/NameNode.go

director: 
	rm -rf DataNode
	mkdir DataNode
	go run Director/Director.go

doshbank: 
	rm -rf DataNode
	mkdir DataNode
	go run DoshBank/DoshBank.go
