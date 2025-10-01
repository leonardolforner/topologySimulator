# topologySimulator
simulador para a disciplina de Simulação e métoods analíticos da PUCRS, alunos: Leonardo Forner, Pedro Semensato, João Pedro Feijó

## Before running
before being able to build this project, it is needed that you install golang into your system
the simulator was built and tested using golang version `go1.24.6` on darwin/arm64 architecture, and it is highly recommended
for you to at least install the same version of go. As there could have been changes in std packages used to calculat ethe simulations
so results can vary slightly if not reproduced correctly.

## How to run
this project has a makefile in it to facilitate building and runnning the program, if you are on a UNIX based system (MacOS/Linux)
just run 

```
go mod vendor
make build
cd bin
sim ...
```

The simulator receives the queue topology exactly the same as the original simulator from the M3 module of the class (as the group though the
format of that yml fila was really good to represent a topology), with the exception of only accepting legal yaml files (no comments before yaml lines, and no !ARGUMENTS in the start of the yaml), see `model.yml` for an example of a legal topology.

to run the simulation, use the `sim` binary, with the -config argument pointing to the topology yaml, alongside the max rands arguments with
whatever number you deem necessary.
```
sim -config ../model.yml --max-rands=1000
```
