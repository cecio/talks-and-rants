# Hackersgen 2023 - Finals

Repo for the Hackersgen 2023 Finals challenge. Docs are in Italian.

Content:

- `VM` folder: executables (both Linux and Windows) for the Virtual Machine used in the challenge

## Challenge Description (ITA)

La sfida si svolge con la stessa Virtual Machine utilizzata durante le Qualificazioni (vedi il folder per la documentazione). Ma questa volta la sfida e' leggermente differente:

la VM e' infettata da un "Virus" che va disabilitato. In fase di startup la VM vi comunica una locazione di memoria e un "Killswitch" che vi permettra' di disabilitare il Virus:

```
WARNING: System Infected!!
Identified KillSwitch (6 chars): 19XZYM
Identified Memory Address for KillSwitch (decimal): 2976
```

Il vostro compito sara' di scrivere un programma in grado di disabilitare il "Virus" tenendo conto che i parametri di Killswithc e Address saranno dinamici. 
Inoltre il "Virus" non se ne sta a guardare, ma scrive, ad ogni istruzione eseguita, un byte 0x91 (RET) ogni 8 bytes della memoria del programma, partendo dalla locazione 8. Ovviamente questo potrebbe corrompere il codice del vostro codice.
