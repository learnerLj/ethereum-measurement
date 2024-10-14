# Ethereum Measurement

## Getting Started
We have set up a one-click script that can handle almost all complex operations. Please run it in an Ubuntu 22.04 environment. All scripts are located in the `script` folder, and the entry point is `start-measure.sh`.

First, you need to configure the initial settings in `script/env.sh`:
- `CHAINPATHS`: A list of folders for storing blockchain data. The number of folders corresponds to the number of nodes to be deployed.
- `RSYS_CONFIG`, `LOGROTATE_CONFIG`, `CRON_CONFIG`: These are log processing configuration files that will be written automatically, so there’s no need to worry about them.
- `GETH`: The executable path of the customized client.
- `PRYSM`: The executable path of the consensus client.

So, you need to compile the modified client before starting the measurement. A few more steps need to be set up:
1. Install the MongoDB database and set up an account with permissions to create databases.
2. Configure `instruct.env` in the `measurement` folder, especially setting up `MONGODB_URI` to ensure it has the necessary access permissions. The `MONGODB_NAME` can be set to any name.

Next, compile the client. Ensure you have installed Golang version 1.22 or higher and the `make` build suite. Run `make geth` in the root directory of the repository. If successful, you will get the following output:
`Run "./build/bin/geth" to launch geth.`
This is the path to the customized geth. Change this relative path to a global one and fill it in the `GETH` configuration in `script/env.sh`.

For the Prysm client, use the following command to download it: 
`curl https://raw.githubusercontent.com/prysmaticlabs/prysm/master/prysm.sh --output prysm.sh && chmod +x prysm.sh`, and configure it in the `script/env.sh` file as well.

Run `bash script/start-measure.sh` to start the process. It will automatically configure the client settings and system settings as mentioned in the appendix. It will also assign ports to avoid conflicts. Each node's data will be stored in the specified folder.
Please note that each client has four critical ports that cannot be used by other programs or be reused by different nodes (however, TCP and UDP ports will not conflict):
1. TCP node communication: Our configuration starts incrementing from port 30300
   - `--port 30305`
2. UDP node discovery: Our configuration starts incrementing from port 30300
   - `--discovery.port 30305`
3. RPC Service: Our configuration starts incrementing from port 8540
   - `--http.port 8547`
4. HTTP-JWT Consensus: Our configuration starts incrementing from port 8600
   - `--authrpc.port 8553`

Please note that the nodes will begin synchronization immediately, which will consume a large amount of memory and disk space. Once fully synchronized, a single node uses about 12GB of memory and 1TB of disk space. If you don't need to continuously collect data, stop the process as soon as the tests are complete.

## Designed Architecture
The main packet aggregation and analysis occur in the `measurement` folder, which handles database connections and real-time analysis. There are also many packet interaction hooks in the `node`, `p2p`, and other folders, used to collect packet requests and responses, as well as client states.