# Integrantes:
Benjamín Gutierrez 202004621-2 
Sofía Parada Hormazábal 202004671-9

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

Instrucciones de compilacion: 
Se debe ubicar en el directorio Lab3SD de cada maquina virtual, en donde se encuentran los archivos.

Para docker:
 
Director: 
- sudo docker build -t director -f Dockerfile.Director .
- docker run -p 8080:8080 -p 8084:8084 director

NameNode: 
- sudo docker build -t namenode -f Dockerfile.NameNode .
- sudo docker run -p 8080:8080 -p 8084:8084 -p 8086:8086 -p 8085:8085 namenode

Notar que make merc no se debe ejecutar hasta que namenode y director printeen su primera linea respectiva!
Luego de esto ejecutar: make merc


!!Se buildean las dos pero al momento de correr no identifican tira error bind: cannot assign requested address, cuando a cada una la corro en su maquina respectiva, y si se corre con go run Director/Director.go y go run NameNode/NameNode.go se ejecutan bien.

Existe un makefile para correr los archivos sin dockerizar, para eso se debe escribir en las terminales lo siguiente y en el siguiente orden:

dist083: make namenode
dist081: make director
dist084: make merc

Existe make clean para borrar el archivo DataNode que almacena datos de mercenarios
