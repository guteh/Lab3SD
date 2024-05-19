# Planificación


Maquinas y contraseñas, junto con asignacion de tareas:

Para ingresar a maquinas leer archivos que subio ayudante a aula, Tutorial de Acceso a Maquinas Virtuales

### Director - DataNode1:
dist@dist081.inf.santiago.usm.cl - xrLQTggUzKS5 IP: 10.35.169.91

### Dosh Bank: - DataNode2:
dist@dist082.inf.santiago.usm.cl - C9UucTLgg7bJ IP: 10.35.169.92

### Name Node
dist@dist083.inf.santiago.usm.cl - YH6LqgC67Qqq IP: 10.35.169.93

### Mercenarios - DataNode3:
dist@dist084.inf.santiago.usm.cl - 6K57TE62BWqQ IP: 10.35.169.94


## Conexiones:

- Director - Dosh Bank: GRPC de vuelta, RABBITMQ de ida
- Director - Mercenarios: GRPC
- Director - NameNode: GRPC
- NameNode - DataNodeX: GRPC


Cada maquina tiene este repositorio github clonado, para actualizarlo escribir *git pull* en la consola y en el directorio Lab3SD

Como cada maquina tiene el repositorio, solo es necesario correr los archivos de la maquina virtual.


Director: 
- sudo docker build -t director -f Dockerfile.Director .
- docker run -p 8080:8080 -p 8084:8084 director

NameNode: 
- sudo docker build -t namenode -f Dockerfile.NameNode .
- sudo docker run -p 8080:8080 -p 8084:8084 -p 8086:8086 -p 8085:8085 namenode

Se buildean las dos pero al momento de correr no identifican tira error bind: cannot assign requested address, cuando a cada una la corro en su maquina respectiva, y si se corre con go run Director/Director.go y go run NameNode/NameNode.go se ejecutan bien.

Hacer makefile que cree DataNode, ejecute los archivos necesarios y despues lo borre


