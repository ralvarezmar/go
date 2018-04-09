Mi práctica está implementada con relojes vectoriales. El reloj vectorial está formado por un vector con el id del agente y un contador.
Cuando el agente crea el archivo inicializa el reloj con 1:1 (en el caso de que su identificador sea 1) y añade la marca de CREATE. 
Si este agente añade otra marca simplemente se incrementa el contador y si es otro agente se añade su id al reloj con contador 1. 

Para hacer pruebas con varios clientes primero usé ./clicustod y ./srvcustod en vez de $HOME/clicustod y $HOME/clicustod. Así cada uno estaba en un directorio y podía probar con distintos ids. 

Para hacer pruebas:
Primero lanzo el servidor con localhost y cualquier puerto

Acto seguido lanzo el cliente(localhost y mismo puerto) con el flag -g para que me devuelva un agente. 

Después de conseguir un identificador para el agente ya puedo pedir objetos con -c(si el agente no tiene id y pide objeto, el servidor no le suministra un id para el objeto)

Después añado marcas con -m. Para añadir marcas con varios identificadores borraba el archivo myid y pedía otro al servidor. 

Para probar que devolvía bien los errores modificaba directamente el objeto borrando marcas y ejecutando el cliente con -a para que lo mandara al servidor. 
Para comprobar que el objeto que subo es correcto o no, obtengo la última marca del objeto que hay en el servidor y el que manda el cliente y comparo: recorro ambos arrays y voy comparando campo por campo. Si es correcto lo subo, si no devuelvo error (imprimo en servidor y se lo mando al cliente) y desecho objeto.


Para finalizar probé en diferentes ordenadores, uno ejecutando el server y otro el cliente (ambos conectándose con la ip del servidor)