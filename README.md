# SnowflakeID

## 概要

分散システムでユニークなIDを生成する手法。  

IDは 64 bitで内訳は
- Sign bit: 1 bit
- Timestamp: 41 bits
  - 特定の日時（epoch）からの経過時間（ミリ秒）をタイムスタンプとして使う
  - 特定の日時はアプリケーション側で好きに決めれる
- Datacenter ID: 5 bits
  - 2 ^ 5 = 32 箇所のデータセンターまで対応
- Machine ID: 5 bits
  - 2 ^ 5 = 32 台まで対応(1データセンターあたり)
- Sequence number: 12 bits
  - 各マシン毎で管理する連番。同一ミリ秒の間はインクリメントしていき、ミリ秒毎に0にリセットされる。
  - 12 bits なので4096個(0 ~ 4095)まで対応

![bits内訳](/snowflake-bits.png)

## 制約

下記の制約がある。

1. Timestampは41bitsなので `2 ^ 41 - 1 = 2199023255551 msec` (69.7年)まで扱える
1. Sequence numberは12bitsなので、生成できるIDは1ミリ秒あたり4096個まで
1. NTPで各サーバーの時刻が同期されていること

1と2に関してはアプリケーション毎に調整可能。  
例えば
- Datacenter ID と Machine ID のビット数を少し削ってTimestampやSequence numberに割り振る
- 業務系でトラフィックはそこまで多くないけど長生きして欲しいなら、Sequence Numberを削ってTimestampに割り振る

# なにが嬉しいのか

- データベースやWEBサーバーが何台あろうが重複しないIDを生成できる。スケールしやすい。
- 数字だけの64bitsで表現できる。UUIDは数値と文字列で128bits
- timestampでソートすれば、おおよそ時系列順に並べれる
  - 「おおよそ」・・・DatacenterIDやMachineIDが異なる同一ミリ秒間に生成されたIDは判別できない（ってことだと思う）

# 自作snowflakeベンチマーク結果

```
goos: darwin
goarch: amd64
pkg: snowflake
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkGenerateID-4               4096                94.45 ns/op            0 B/op          0 allocs/op
```

# 参考
- [wikipedia](https://en.wikipedia.org/wiki/Snowflake_ID)
- [Announcing Snowflake](https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake)
