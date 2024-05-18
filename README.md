# Planificación


Maquinas y contraseñas, junto con asignacion de tareas:

Para ingresar a maquinas leer archivos que subio ayudante a aula, Tutorial de Acceso a Maquinas Virtuales

### Director - DataNode1:
dist@dist081.inf.santiago.usm.cl - xrLQTggUzKS5 IP: 10.35.169.91

### Dosh Bank:
dist@dist082.inf.santiago.usm.cl - C9UucTLgg7bJ IP: 10.35.169.92

### Name Node - DataNode2:
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

hola
