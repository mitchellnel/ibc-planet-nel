# ibcplanetnel

**ibcplanetnel** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

Goals:

-   Use IBC to create and send packets between blockchains
-   Navigate between blockchains using the Cosmos SDK and the Ignite CLI Relayer
-   Create a basic blog post and save the post on another blockchain

## What is IBC?

The Inter-Blockchain Communication protocol (IBC) allows blockchainsg to talk to each other.

IBC handles transport across different sovereign blockchains.

This end-to-end, connection-oriented, stateful protocol provides reliable, ordered, and authenticated communication between heterogenous blockchains.

The [IBC protocol in the Cosmos SDK](https://docs.cosmos.network/master/ibc/overview.html) is the standard for the interaction between two blockchains.

The IBCmodule interface defines how packets and messages are constructed to be interpreted by the sending and the receiving blockchain.

The IBC relayer lets you connect between sets of IBC-enabled chains.

This tutorial teaches us how to create two blockchains, and then start and use the relayer with Ignite CLI to connect the two blockchains.

This tutorial covers essentials like modules, IBC packets, relayer, and the lifecycle of packets routed through IBC.

We will create a blockchain app with a blog module to write posts on other blockchains that contain the Hello World message.

## Get started

First, clone this repository:

```bash
git clone https://github.com/mitchellnel/ibc-planet-nel
```

### Start up the blockchains

Open an additional two terminal windows.

On one window, invoke:

```bash
ignite chain serve -c earth.yml
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

On the other window, invoke:

```bash
ignite chain serve -c mars.yml
```

### Remove existing Relayer and Ignite CLI configurations

If you previously used the relayer, follow these steps to remove existing relayer and Ignite CLI configurations:

-   Stop your blockchains
-   Delete previous configuration files:

```bash
rm -rf ~/.ignite/relayer
```

If existing relayer configurations do not exist, the command returns `no matches found`
 and no action is taken.

### Configure and start the Relayer

First, we’ll configure the relayer.

Use the Ignite CLI `configure` command with the `--advanced` option:

```bash
ignite relayer configure -a \
  --source-rpc "http://0.0.0.0:26657" \
  --source-faucet "http://0.0.0.0:4500" \
  --source-port "blog" \
  --source-version "blog-1" \
  --source-gasprice "0.0000025stake" \
  --source-prefix "cosmos" \
  --source-gaslimit 300000 \
  --target-rpc "http://0.0.0.0:26659" \
  --target-faucet "http://0.0.0.0:4501" \
  --target-port "blog" \
  --target-version "blog-1" \
  --target-gasprice "0.0000025stake" \
  --target-prefix "cosmos" \
  --target-gaslimit 300000
```

When prompted, press Enter to accept the default values for Source Account and Target Account.

The output should look like:

```
------
Setting up chains
------

? Source Account default
? Target Account default

🔐  Account on "source" is default(cosmos1c9ltr9gq2jgk0u9da8sz7tdga5eghheepy2car)

received coins from a faucet
 |· (balance: 100000stake,5token)

🔐  Account on "target" is default(cosmos1c9ltr9gq2jgk0u9da8sz7tdga5eghheepy2car)

received coins from a faucet
 |· (balance: 100000stake,5token)

⛓  Configured chains: earth-mars
```

Then, open another terminal window, and start the relayer process:

```bash
ignite relayer connect
```

And this will output:

```
------
Paths
------

earth-mars:
    earth > (port: blog) (channel: channel-0)
    mars  > (port: blog) (channel: channel-0)

------
Listening and relaying packets between chains...
------
```

### Send packets

We can now send packets between the chains, and verify the received posts:

```bash
ibc-planet-neld tx blog send-ibc-post blog channel-0 "Hello Mars" \
  "Hello Mars, I'm Alice from Earth" \
  --from alice --chain-id earth --home ~/.earth
```

Then, to verify that this post has been received on Mars:

```bash
ibc-planet-neld q blog list-post --node tcp://localhost:26659
```

And this will output:

```
Post:
- content: Hello Mars, I'm Alice from Earth
  creator: blog-channel-0-cosmos15t4xcfrp3se2cmaklp5p6825wwdal9hwnhr3j8
  id: "0"
  title: Hello Mars
pagination:
  next_key: null
  total: "0"
```

Then, to check if the packet's receipt has been acknowledged on Earth:

```bash
ibc-planet-neld q blog list-sent-post
```

And this will output:

```
SentPost:
- chain: blog-channel-0
  creator: cosmos15t4xcfrp3se2cmaklp5p6825wwdal9hwnhr3j8
  id: "0"
  postID: "0"
  title: Hello Mars
pagination:
  next_key: null
  total: "0"
```

To test timeout, set the timeout time of a packet to 1 nanosecond, verify that the packet is timed out, and check the timed-out posts:

```bash
ibc-planet-neld tx blog send-ibc-post blog channel-0 "Sorry Mars" \
  "Sorry Mars, you will never see this post" \
  --from alice --chain-id earth --home ~/.earth --packet-timeout-timestamp 1
```

Then, check the timed out posts:

```bash
ibc-planet-neld q blog list-timedout-post
```

And this will output:

```
TimedoutPost:
- chain: blog-channel-0
  creator: cosmos15t4xcfrp3se2cmaklp5p6825wwdal9hwnhr3j8
  id: "0"
  title: Sorry Mars
pagination:
  next_key: null
  total: "0"
```

And no this post was not received by Mars:

```bash
ibc-planet-neld q blog list-post --node tcp://localhost:26659
```

We see only our first post was received:

```
Post:
- content: Hello Mars, I'm Alice from Earth
  creator: blog-channel-0-cosmos15t4xcfrp3se2cmaklp5p6825wwdal9hwnhr3j8
  id: "0"
  title: Hello Mars
pagination:
  next_key: null
  total: "0"
```

We can also send a post from Mars to Earth:

```bash
ibc-planet-neld tx blog send-ibc-post blog channel-0 "Hello Earth" \
  "Hello Earth, I'm Alice from Mars" \
  --from alice --chain-id mars --home ~/.mars --node tcp://localhost:26659
```

And verify that the post was received on Earth:

```bash
ibc-planet-neld q blog list-post
```

Outputs:

```
Post:
- content: Hello Earth, I'm Alice from Mars
  creator: blog-channel-0-cosmos1k09q34ecq78qvzlghmvmglt96e09yhzs942938
  id: "0"
  title: Hello Earth
pagination:
  next_key: null
  total: "0"
```

### Relayer logs

We can also see that every time we sent a packet and received an ACK, our relayer logged these events:

```
------
Paths
------

earth-mars:
    earth > (port: blog) (channel: channel-0)
    mars  > (port: blog) (channel: channel-0)

------
Listening and relaying packets between chains...
------

Relay 1 packets from earth => mars
Relay 1 acks from mars => earth
Relay 1 packets from mars => earth
Relay 1 acks from earth => mars
```
