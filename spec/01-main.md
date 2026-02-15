# Spec
## General

Dieses Projekt trägt den Namen `ets-gruppenadressgenerator`.
Inhaltlich handelt es sich um ein Golang Programm, welches gewisse Eingabedatein einliest und daraus eine XML Datei erzeugt.

## Business Scope

Die erzeugte XML Datei ist für die Software ETS (derzeitige Version ETS6) gedacht, mit welcher man KNX Systeme konfiguriert. Die erzeugte Datei enthält die sogenannten Gruppenadressen, welche eine dreistufige Hierarchie haben:

1. Hauptgruppe (max. 32 Stück)
2. Mittelgruppe (max. 8 Stück)
3. Gruppenadressen (max. 256 Stück)

Die Gruppenadressen werden in der ETS genutzt um den KNX Geräten ihre Funktionen zuzuweisen.

## Requirements

* Die Adressen aller Adressblöcke beginnen mit dem Index `0`.
* Modelliere die Räume (Eingabeschlüssel `rooms`, inkl. Stockwerk) als Hauptgruppen.
* Füge zu den gefundenen Räumen die folgenden zentralen Hauptgruppen hinzu:
   * Zentral
   * UG-Zentral
   * EG-Zentral
* Modelliere die Gewerke (Eingabeschlüssel `trades`) als Mittelgruppe für die Räume. Die zentralen Hauptgruppen benötigen keine Mittelgruppen.
* Entnehme die Gewerke aus der Eingabedatei.
* Die Objekte eines Gewerks haben immer einen `type`. Dieser bestimmt, welche Gruppenadressen für dieses Objekt notwendig sind. Die Objekte und deren `type` befinden sich ebenfalls in der Eingabedatei.
* Das Feld `type` eines Objekts lässt sich mit der Datei `knx-object-types.json` auflösen. Hierin sind die Funktionen definiert, die ein KNX Objekt eines gewissen Types benötigt.
* Der Name einer Gruppenadresse setzt sich aus dem der Hauptgruppe, Mittelgruppe sowie dem Objekt und dessen Funktion zusammen. Zwischen den drei Hierarchiestufen befindet sich ein `_` als Delimiter. 

## Technical design

### Input
Die derzeitige angedachte Eingabedatei ist `input.json` welche in diesem Projekt zu finden ist. 

### Output
Die erwartete Ausgabestruktur kannst du der Datei `spec/group-addresses-example1.xml` in diesem Projekt entnehmen. Hierin finden sich einige Haupgruppen und in einer Haupgruppe auch Mittelgruppen und eine konkrete Gruppenadresse.

