# CS2 Cheat Detection — Machine Learning Pipeline

> End-to-end machine learning pipeline for detecting cheaters in Counter-Strike 2 by analyzing replay demo files.

[![Python](https://img.shields.io/badge/Python-3.10+-3776AB?logo=python&logoColor=white)](https://python.org)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![scikit-learn](https://img.shields.io/badge/scikit--learn-1.x-F7931E?logo=scikit-learn&logoColor=white)](https://scikit-learn.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Overview

This project implements a **two-stage cheat detection system** for Counter-Strike 2 (CS2):

1. **Event-Level Classification** — A Random Forest model analyzes individual kill events to detect aimbot/wallhack patterns (aim angle snaps, kills through walls, etc.).
2. **Player-Level Classification** — A Logistic Regression model aggregates event-level predictions per player to flag suspected cheaters.

The system parses raw CS2 demo files (`.dem`), extracts gameplay features, engineers cheat-indicative signals, and runs ML models to produce cheat probability scores.

### Key Results

| Model | Metric | Value |
|:------|:-------|------:|
| Event-Level (Random Forest) | ROC AUC | ~0.65 |
| Player-Level (Logistic Regression) | Cheater F1 | 0.50 |
| Player-Level (Logistic Regression) | Accuracy | 84% |

---

## Architecture

```
CS2 Demo Files (.dem)
        │
        ├──► Go Parser (main.go)          → aimDelta extraction
        │        uses demoinfocs-golang
        │
        ├──► Python Parser (Awpy)         → kill events + visibility checks
        │
        └──► Combined CSV Dataset
                    │
                    ▼
            Feature Engineering
            (distance, log_aimDelta, aimDelta_diff, enemy_visible)
                    │
                    ▼
        ┌───────────────────────┐
        │  Event-Level Model    │
        │  (Random Forest)      │──► per-kill cheat probability
        └───────────────────────┘
                    │
                    ▼
        ┌───────────────────────┐
        │  Player-Level Model   │
        │  (Logistic Regression)│──► per-player cheat flag
        └───────────────────────┘
```

---

## Project Structure

```
.
├── CS2_Cheat_Detection_Final_Report.ipynb   # Full pipeline notebook (report)
├── combined_wallhack_aim_engineered.csv     # Processed kill-event dataset
├── go_parser_custom/                        # Custom Go demo parser
│   ├── main.go                              # Extracts aimDelta from .dem files
│   ├── parse_all.bat                        # Batch script to parse multiple demos
│   ├── go.mod / go.sum                      # Go module dependencies
├── go_parser/                               # demoinfocs-golang library (dependency)
├── *.png                                    # Generated plots & visualizations
│   ├── event_level_roc.png                  # Event-level ROC curve
│   ├── playerlevelROC.png                   # Player-level ROC curve
│   ├── feature_distributions.png            # Feature distribution histograms
│   └── demos.png / demoPrredcition.png      # Demo screenshots
└── README.md
```

---

## Getting Started

### Prerequisites

- **Python 3.10+** with pip
- **Go 1.21+** (for demo parsing)
- CS2 demo files (`.dem`)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/peszew/CS2-Cheat-Detection.git
   cd CS2-Cheat-Detection
   ```

2. **Install Python dependencies:**
   ```bash
   pip install pandas numpy scikit-learn matplotlib awpy
   ```

3. **Install Go dependencies (for demo parsing):**
   ```bash
   cd go_parser_custom
   go mod tidy
   cd ..
   ```

### Usage

#### Step 1: Parse Demo Files

Use the Go parser to extract aim data from `.dem` files:

```bash
cd go_parser_custom
go run main.go path/to/demo.dem
```

Or parse all demos in a folder:

```bash
cd go_parser_custom
parse_all.bat
```

#### Step 2: Run the Notebook

Open and run the Jupyter notebook for the full pipeline:

```bash
jupyter notebook CS2_Cheat_Detection_Final_Report.ipynb
```

> **Note:** Update the `RAW_CSV` path in the notebook to point to your local dataset.

---

## Features Extracted

| Feature | Description | Cheat Signal |
|:--------|:------------|:-------------|
| `aimDelta` | Change in aim angles over last 6 ticks | Aimbot: unnaturally fast/precise snaps |
| `enemy_visible` | Line-of-sight check (via Awpy) | Wallhack: kills through walls |
| `distance` | 3D Euclidean distance to victim | Wallhack: long-range non-visible kills |
| `log_aimDelta` | Log-transformed aimDelta | Normalizes distribution |
| `aimDelta_diff` | Per-player sequential diff of log_aimDelta | Sudden aim behavior changes |
| `weapon` | Weapon used for the kill | Contextual feature |
| `map` | Map played | Contextual feature |

---

## Results & Visualizations

### Event-Level ROC Curve
The Random Forest achieves an AUC of ~0.65, indicating moderate ability to distinguish cheat vs. legitimate kill events.

### Player-Level Classification
Using aggregated event scores (mean/max probability, cheat rate), the Logistic Regression model achieves:
- **Legit class:** Precision 0.93, Recall 0.88
- **Cheater class:** Precision 0.44, Recall 0.57

These results reflect the challenge of limited labeled data and noisy features.

---

## Future Improvements

- **More features** — shot spread, hit probability, kill timing, movement patterns
- **More data** — larger labeled dataset with confirmed cheaters
- **Advanced models** — gradient boosting, neural networks, time-series models
- **Threshold calibration** — optimize player-level thresholds via ROC analysis
- **Anomaly detection** — one-class SVM, isolation forests for unsupervised approaches
- **Real-time system** — streaming predictions during live matches

---

## Tech Stack

- **Python** — pandas, NumPy, scikit-learn, matplotlib, Awpy
- **Go** — demoinfocs-golang (CS2 demo parsing)
- **Jupyter Notebook** — analysis & reporting

---

## Author

**Piotr Szewczyk**

EMLET Course Project

---

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

---

## Acknowledgments

- [demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang) — CS2/CS:GO demo parser
- [Awpy](https://github.com/pnxenopoulos/awpy) — Python CS analytics library
- [scikit-learn](https://scikit-learn.org) — machine learning framework
