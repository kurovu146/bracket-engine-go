# bracket-engine-go

Tournament bracket generation library for Go. Supports 5 formats with full match linking, seeding, and bye handling.

Ported from [@kurovu146/bracket-engine](https://www.npmjs.com/package/@kurovu146/bracket-engine) (TypeScript).

## Install

```bash
go get github.com/kurovu146/bracket-engine-go@latest
```

## Formats

| Format | Function | Description |
|--------|----------|-------------|
| Single Elimination | `GenerateSingleElimination` | Knockout bracket with optional 3rd place match |
| Double Elimination | `GenerateDoubleElimination` | Winners + Losers bracket + Grand Final |
| Round Robin | `GenerateRoundRobin` | Every participant plays every other |
| Swiss | `GenerateSwiss` | Paired by standings each round (R1 by seed, R2+ placeholder) |
| Group Stage | `GenerateGroupStage` | Groups play round-robin internally |

## Quick Start

```go
package main

import (
    "fmt"
    bracket "github.com/kurovu146/bracket-engine-go"
)

func main() {
    players := []string{"Alice", "Bob", "Charlie", "Dave"}

    // Single Elimination
    matches, err := bracket.GenerateSingleElimination(players, nil)
    if err != nil {
        panic(err)
    }
    for _, m := range matches {
        fmt.Printf("%s: %v vs %v (%s)\n", m.MatchID, m.Player1ID, m.Player2ID, *m.RoundName)
    }
    // WB-R1-M1: Alice vs Dave (Semi-final)
    // WB-R1-M2: Bob vs Charlie (Semi-final)
    // WB-R2-M1: <nil> vs <nil> (Final)
}
```

## API

### Single Elimination

```go
matches, err := bracket.GenerateSingleElimination(playerIDs, &bracket.SingleEliminationOptions{
    ThirdPlaceMatch: true,
    BestOf: &bracket.BestOfConfig{
        Default: bracket.IntPtr(3),
        Final:   bracket.IntPtr(5),
    },
})
```

### Double Elimination

```go
matches, err := bracket.GenerateDoubleElimination(playerIDs, &bracket.DoubleEliminationOptions{
    GrandFinalReset: true,
    BestOf: &bracket.BestOfConfig{
        Default:    bracket.IntPtr(3),
        GrandFinal: bracket.IntPtr(5),
    },
})
```

### Round Robin

```go
matches, err := bracket.GenerateRoundRobin(playerIDs, &bracket.RoundRobinOptions{
    DoubleRoundRobin: true,
})
```

### Swiss

```go
matches, err := bracket.GenerateSwiss(playerIDs, &bracket.SwissOptions{
    NumRounds: bracket.IntPtr(5),
})
```

### Group Stage

```go
result, err := bracket.GenerateGroupStage(playerIDs, &bracket.GroupStageOptions{
    NumGroups:    bracket.IntPtr(4),
    Distribution: "snake",
})
// result.Groups  — [][]string (group assignments)
// result.Matches — []MatchSeed (all matches across groups)
```

## MatchSeed

Every generated match contains:

```go
type MatchSeed struct {
    MatchID             string      // Stable ID: "WB-R1-M1", "LB-R2-M3", "GF-M1", "G0-R1-M2"
    Round               int         // 1-based round number
    MatchNumber         int         // Sequential number across tournament
    Player1ID           *string     // nil = TBD/bye
    Player2ID           *string     // nil = TBD/bye
    BracketType         BracketType // "winners", "losers", "grand_final", "group_0", etc.
    NextMatchIndex      *int        // Array index of winner's next match
    LoserNextMatchIndex *int        // Array index of loser's next match (double elim)
    NextMatchSlot       *MatchSlot  // "player1" or "player2"
    LoserNextMatchSlot  *MatchSlot
    RoundName           *string     // "Final", "Semi-final", "LB Round 3", etc.
    IsBye               bool        // Auto-advance (no opponent)
    BestOf              *int        // Best-of hint (e.g. 3)
}
```

## Match ID Formats

| Bracket | Format | Example |
|---------|--------|---------|
| Winners | `WB-R{round}-M{match}` | `WB-R1-M1` |
| Losers | `LB-R{round}-M{match}` | `LB-R3-M2` |
| Grand Final | `GF-M1` | `GF-M1` |
| Grand Final Reset | `GF-M2` | `GF-M2` |
| 3rd Place | `3RD-M1` | `3RD-M1` |
| Round Robin | `RR-R{round}-M{match}` | `RR-R2-M3` |
| Swiss | `SW-R{round}-M{match}` | `SW-R1-M4` |
| Group N | `G{N}-R{round}-M{match}` | `G0-R1-M2` |

## Seeding

Standard tournament seeding with automatic bye handling:

- 8-player seed order: `[1, 8, 4, 5, 2, 7, 3, 6]`
- Non-power-of-2 counts: top seeds receive byes
- Byes are marked with `IsBye: true`

## Design

This library generates **root tournament formats only**. Multi-phase tournaments (e.g. Group Stage -> Single Elimination) should be orchestrated by combining multiple generators at the app layer.

Swiss R2+ matches are generated as placeholders with nil players. The app must fill them dynamically based on standings after each round.

## License

MIT
