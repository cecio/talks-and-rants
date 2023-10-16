# La sfida

Il pacchetto contiene un eseguibile (vm o vm.exe nei due formati Linux e Windows) con l'implementazione di una Custom Virtual Machine. La VM contiene un "Secret Message" criptato in memoria in posizione casuale contigua (riassegnata ad ogni esecuzione), con il seguente algoritmo

- XOR con chiave: `theblackpirate`
- due bitwise Rotate Left (`ROL`)

Per risolvere il challenge e' necessario scrivere un programma nel linguaggio della VM (documentato nel file "doc.png") per estrarre il "Secret Message". Sara' necessario inviare il programma (la sequenza di bytecodes) alla mail `hackersgen.challenge2023@sorint.com` : il programma, una volta inserito, dovra' stampare in esadecimale il "Secret Message" decifrato.

La Virtual Machine ha a disposizione una modalita' di debug (accessibile tramite parametro "debug") per visualizzare lo stato dei registri e dei flag ad ogni step di esecuzione.
