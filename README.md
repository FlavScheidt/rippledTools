# rippledTools
Automatized tasks for configuring and analyzing rippled validators

We assume that every node in the network can access each other by the name.

# Config Cluster

Configures the cluster topology according to the desired type of test. It is basically an UNL configuration tool
Creates the rippled.cfg and validators.txt file for each node and sends to each one fo them via ssh.

Usage is:

```
$./generate_config_rippled.cfg <test type>

```

## Test types
**general** - generates a fully connected topology, similar to the real XRPL network

**unl** - generates the trust overlay based in pre-defined UNLs (see [Configuration](https://github.com/FlavScheidt/rippledTools/edit/main/README.md#configuration))

**validator** - generates a topology where each node chooses the validators they trust separately

## Configuration

#### Validators keys
We assume that the rippled keys for each validator were created previously. All the keys must be stored in separate .txt files on ConfigCluster/keys with the name of the file being the name of the node

#### General configuration files
**ClusterConfig.csv**  - lists the nodes and their info
  #####Format:
  ```
    IP,name,rippledKey,unl
 ```
**rippled.cfg** - master model for the rippled.cfg, with all an empty "validators" session.
**validators.txt** - master model for the validators.txt file

#### Topology Speficiations
The topology specifications for the **general** test must be in the ConfigCluster/unl directory. The **unl** test is configured directly on the ClusterConfig.csv, the field "unl" specifies to which UNL a node is subscribed to. The path for the file containing which node belongs to each UNL (UNL_DIR) must be set on generate_config_rippled.sh 


# New Run
Deletes databses and logs for a fresh restart (rmeotely).

Usage:
```
$ ./prepareNewRun
```
## Configuration
**nodes.txt** names of the nodes

# InspectSync
Script for inspecting if the nodes are synchronized or not (via ssh)

Usage:
```
./inspect.sh
```
Returns:

```
datetime nodeName lastClosedLedger index timeClosed parentLedger openIndex
```

# ExportLogs
Export logs from rippled, sntRippled and gossipGoSnt to a given destination

Usage:
```
./export,sh
```

##Configuration
Directly on export.sh


