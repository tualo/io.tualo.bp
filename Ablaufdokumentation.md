# Ablaufdokumentation der Bildverarbeitung

## 1. Bewegungsunschärfe-Erkennung

Die Funktion `continuous.Process` nimmt ein Bild entgegen und prüft dieses darauf,
ob Bewegungsunschärfe vorhanden ist.


## 2. Kontour-Extraktion

Die Funktion `continuous.Paper` extrahiert die größte Kontour und prüft dabei, 
ob diese einer Seite entspricht. Die äußere Kontour wird markiert.
Falls eine passende Kontour gefunden wird, wird diese an die Funktion `Sheet` übergeben.

## 3. Seitenanalyse

Die Funktion `continuous.Sheet` analysiert die gefundene Seite und versucht dabei, 
den Typ des Stimmzettels zu identifizieren (`tesseract.Static().DetectBallotpaperType`).
Basierend auf dem ermittelten Typ wird die Skalierung für nachfolgende Bildverarbeitungsschritte bestimmt. 
Mithilfe der Skalierung können die Millimeterangaben des Backends verwendet werden, 
um beispielsweise die Kreuzregionen zu bestimmen.

### 3.1 Barcode-Erkennung

Das Bild wird nach Barcodes durchsucht. 
Spezielle Barcodes für Kisten und Stapeltrenner werden dabei herausgefiltert.

**TODO:** Implementierung einer Prüfsummen-Validierung

Wird ein Paginationscode (Code 128, mindestens 5 Zeichen) gefunden, wird die Seite 
auf die definierten Regionen von Interesse (ROI) - die Kreuzregionen - untersucht.

Entsprechend dem identifizierten Stimmzetteltyp wird jede Region von Interesse 
an `continuous.SheetRois` übergeben.

## 4. Kreuzregionen-Analyse

Die Funktion `continuous.SheetRois` untersucht jede Region auf Kreuzmarkierungen.
Dabei wird eine Kreuzmarkierung mittels verschiedener Verfahren analysiert:
- Einfache Bildverarbeitung
- Optional: Lokales neuronales Netzwerk
- Optional: Neuronaler-Netzwerk-Service (extern)




