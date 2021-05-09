# Vanity Arweave Wallet Generator

> _W.A.V.E_, if you squint your eyes hard enough

This is just a "Hello World" project I made to learn Golang

## Usage

```
λ wave -h
wave 0.2.0 (8e7c6876)
Usage: wave [--workers WORKERS] [--number NUMBER] [--output OUTPUT] PATTERN

Positional arguments:
  PATTERN                Regex pattern to match the wallet address

Options:
  --workers WORKERS, -w WORKERS
                         Number of workers to spawn [default: 4]
  --number NUMBER, -n NUMBER
                         Number of wallets to generate [default: 1]
  --output OUTPUT, -o OUTPUT
                         Output directory [default: ./keyfiles]
  --help, -h             display this help and exit
  --version              display version and exit
```

### Example

```
λ wave '^[^-_]+$'
Pattern: /^[^-_]+$/
Outputs: keyfiles
Workers: 4
Wallets: 1
[WORKER1] address: khBnTsl41rPqALdd0PyC1XKEjtEap1i_3Tb1qzE4kSE | match: false]
[WORKER2] address: -7SMqFlCBSAvtj73bVm7Oc5XyiF-0tMaUFJ672FLdDE | match: false]
[WORKER2] address: LKQaRd8gMzZ71RVF-Q51bykX-ByK9PCkFpWvVF1K31Q | match: false]
[WORKER3] address: WLv_zn3vMQNZhsDLrMkbjsZLbf-bgkWixj0L5Fv_sU4 | match: false]
[WORKER4] address: _aQvExk3q-SUffgrvjktm1M4uEzO_qykoGeL8LXgTaQ | match: false]
[WORKER2] address: rt8y1-v2gUHFOQ_ZG_oQljpjYeb2zIWCGxGrgvB5t5M | match: false]
[WORKER1] address: XXaN9g_M-OoUyWcLlkOybyBJdHKPBO4QZF9uW9XxSMk | match: false]
[WORKER3] address: 93LJIZAIjA13hT4YQeM7zYc3asrQ8f6Mj7xycLNtw58 | match: true]
[MATCH] address: 93LJIZAIjA13hT4YQeM7zYc3asrQ8f6Mj7xycLNtw58
[EMIT] keyfile: keyfiles/arweave-keyfile-93LJIZAIjA13hT4YQeM7zYc3asrQ8f6Mj7xycLNtw58.json
```

## License

This software is licensed under [The MIT License](./LICENSE).
