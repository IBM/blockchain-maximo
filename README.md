# Track Maximo Assets / Work Orders on a Blockchain Ledger

In this Code Pattern, we will demonstrate how to leverage "Automation Scripts" to track Maximo Assets within a Blockchain ledger.

This is targeted towards Maximo users that have assets that may be managed by third party vendors and regulators. The blockchain ledger provides a "single source of truth" for all interested parties.

When the reader has completed this Code Pattern, they will understand how to leverage Maximo Automation Scripts to propagate changes to a blockchain ledger.

<img src="https://i.imgur.com/lKeuzAn.png">
<!-- ![Architecture](https://i.imgur.com/lKeuzAn.png) -->


#  Components

* [Maximo](https://www.ibm.com/products/maximo). This is an asset management platform that allows enterprises to automate workflows such as maintenance, work orders, and more.

Sign up for a trial account of Maximo [here](https://www.ibm.com/account/reg/us-en/signup?formid=urx-20869)


# Flow

1. Work Order is created/updated in Maximo.
2. Work Order data is stored in blockchain ledger.
3. External parties are notified of work order via UI.
4. External parties update work order in UI, updates are stored in immutable blockchain.

# Prerequisites

* [Maximo](https://www.ibm.com/products/maximo). This is an asset management platform that allows enterprises to manage and automate workflows such as maintenance, work orders, and more.

* [Docker Engine](https://docs.docker.com/install/). This is a tool that allows for complex applications to packaged as containers. We'll use docker here to deploy a containerized Hyperledger Fabric blockchain network

* [Ngrok](https://ngrok.com/). This is a service that enables a server to be exposed via a public url. This is necessary in this example because we are running our server and blockchain ledger on a local machine.

# Steps

Follow these steps to setup and run this Code Pattern.

1. [Start Ngrok](#1-install-and-start-ngrok)
2. [Create Automation Script in Maximo](#2-create-automation-script-in-maximo)
3. [Create Work Order in Maximo](#3-create-work-order-in-maximo)
4. [Deploy Blockchain Ledger](#4-deploy-blockchain-ledger)
5. [Deploy Web Application](#5-deploy-backend)
6. [Simulate Work Order Update](#6-simulate-work-order-update)

<!-- 5. [Create a Dashboard](#4-create-dashboard) -->

## 1. Install and start Ngrok

Here we will start ngrok to expose our blockchain server on a public address. This is necessary for the hosted Maximo service to send requests to our local machine.

```
#Linux
wget https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-amd64.zip

#Mac OS X
wget https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-darwin-386.zip

# Extract ngrok to path
unzip -d /usr/local/bin/ ngrok-stable*zip

# Start ngrok server to forward traffic to port 8000
ngrok http 8000
```

You should see some output like so. Take note of the public URL (http://a6fe75e8.ngrok.io in this example).

```
ngrok by @inconshreveable                                       (Ctrl+C to quit)

Session Status                online
Session Expires               7 hours, 59 minutes
Version                       2.3.35
Region                        United States (us)
Web Interface                 http://127.0.0.1:4040
Forwarding                    http://a6fe75e8.ngrok.io -> http://localhost:8000
Forwarding                    https://a6fe75e8.ngrok.io -> http://localhost:8000

Connections                   ttl     opn     rt1     rt5     p50     p90
                              0       0       0.00    0.00    0.00    0.00
```

Replace the `url` variable with your ngrok URL (here)[https://github.com/IBM/blockchain-maximo/blob/master/maximo_automation_scripts/work_order_script.py#L57]

## 2. Create Automation Script in Maximo

Login to Maximo Dashboard

<img src="https://i.imgur.com/PdqHYup.png">

In the "Find Navigation Item" input box, search for the "Automation Scripts" option.

<img src="https://i.imgur.com/HH1poFr.png">

Hover over the "Create" button, and select "Create Script with Object Launch Point". This "Object Launch Point" enables a script to be executed when a specific object is updated.

<img src="https://i.imgur.com/gYhLBA2.png">


On the first page, give the Launch Point a name, and specify which object type should trigger the script when updated.

Also select the type of Object we want to trigger our script.

In this case, we'll select the "WORKORDER" object, as we'll want to allow third party vendors to track certain work orders.


We'll also add an optional "Object Event Condition", which allows us to define conditionals restricting which objects will trigger the script. In this case, we'll add a simple conditional "VENDOR is not null", which is saying that we only want to sync work orders that have be assigned to a vendor.

<img src="https://i.imgur.com/iNmx08T.png">



Click the "next" button. Now, enter a name and language for the script. We'll use python here.

<img src="https://i.imgur.com/5AhvIEG.png">


Click "Next" a final time, and then copy and paste the code from our included script at [maximo_automation_scripts/work_order_script.py](maximo_automation_scripts/work_order_script.py). *Confirm that the `url` variable matches your ngrok address from Step 1.*

Click the "Create" to finalize the script.

To see more on launch points, please visit the following [page](https://www.ibm.com/support/knowledgecenter/SSANHD_7.5.3/com.ibm.mbs.doc/autoscript/c_ctr_launch_points.html)

## 3. Create Work Order in Maximo

Next, we'll create a Work Order. In the "Find Navigation Item" section, search for "Work Order Tracking"

Click "New Work Order".

Click "Save Work Order". We will return to this section once the blockchain ledger and Web Application are up and running

## 4. Deploy Blockchain Ledger

Add the following entries to your `/etc/hosts` file.

```
127.0.0.1 orderer.example.com
127.0.0.1 ca.example.com
127.0.0.1 peer0.org1.example.com
```

Run the following commands in a terminal
```
cd local
./startFabric.sh
./installChaincode.sh
```

Running `docker ps` should result in output like so
```
Kalonjis-MacBook-Pro:~ kkbankol@us.ibm.com$ docker ps
CONTAINER ID        IMAGE                                                                                                     COMMAND                  CREATED             STATUS              PORTS                                            NAMES
e3e1c879eefa        dev-peer0.org1.example.com-maximo-1.23-08526223df0db7d64ee8f87b65fba71d9a5024b5345275fb79c16157019ac4cf   "chaincode -peer.add…"   2 hours ago         Up 2 hours                                                           dev-peer0.org1.example.com-maximo-1.23
dfe5aaee9809        hyperledger/fabric-peer                                                                                   "peer node start"        31 hours ago        Up 31 hours         0.0.0.0:7051->7051/tcp, 0.0.0.0:7053->7053/tcp   peer0.org1.example.com
72b28036782e        hyperledger/fabric-ccenv                                                                                  "/bin/bash -c 'sleep…"   31 hours ago        Up 31 hours                                                          chaincode
56076b5f6008        hyperledger/fabric-orderer                                                                                "orderer"                31 hours ago        Up 31 hours         0.0.0.0:7050->7050/tcp                           orderer.example.com
a3f31fadb77a        hyperledger/fabric-couchdb                                                                                "tini -- /docker-ent…"   31 hours ago        Up 31 hours         4369/tcp, 9100/tcp, 0.0.0.0:5984->5984/tcp       couchdb
b85f8b60ea6a        hyperledger/fabric-ca                                                                                     "sh -c 'fabric-ca-se…"   31 hours ago        Up 31 hours         0.0.0.0:7054->7054/tcp                           ca.example.com
9be81686e614        hyperledger/fabric-tools                                                                                  "/bin/bash"              31 hours ago        Up 31 hours                                                          cli
```


## 5. Deploy Web Application

Clone repository using the git cli

```
git clone https://github.com/IBM/blockchain-maximo
```


### Install Node.js packages

If expecting to run this application locally, please install [Node.js](https://nodejs.org/en/) and NPM. Windows users can use the installer at the link [here](https://nodejs.org/en/download/)

If you're using Mac OS X or Linux, and your system requires additional versions of node for other projects, we'd suggest using [nvm](https://github.com/creationix/nvm) to easily switch between node versions. Install nvm with the following commands

```bash
curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.11/install.sh | bash
# Place next three lines in ~/.bash_profile
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion
nvm install v8.9.0
nvm use 8.9.0
```

To run the dashboard locally, we'll need to install a few node libraries which are listed in our `package.json` file. These include the following
- [Vue.js](https://vuejs.org/): Used to simplify the generation of front-end components
- [Express.js](https://expressjs.org/): Used to provide custom api endpoints

These libraries can be installed by entering the following commands in a terminal.

```
cd backend
npm install
cd ..
cd frontend
npm install
```

After installing the prerequisites, we can start the application.

Run the following to start the backend
```
cd backend
npm start
```

In a separate terminal, run the following to start the frontend UI
```
cd frontend
npm run serve
```


Confirm you can access the Dashboard UI at [http://localhost:8080](http://localhost:8080).


<img src="https://i.imgur.com/I216GCw.png">


## 6. Simulate Work Order Update

Now that our system is up and running, we can simulate a typical workflow.

In this case, we'll take a Work Order that is reporting a Hazmat situation in a work site. Once a vendor has been assigned to the Work Order, they'll be able to see the Work Order status, asset, and other information in this UI.

Our web app demo implements the following simulated flow.

1. Hazmat Inspector verifies that there is a hazardous situation, and assesses how severe it is.
2. Hazmat Vendor begins working on fixing the situation.
3. Building Inspector checks work, and verifies when it is complete.

Begin by going back to the Maximo "Work Orders" page, and selecting your previously created work order.

Scroll down to the "Responsibility" section and click on the arrow next to "Vendor".

<img src="https://i.imgur.com/zio0UmS.png">


Click "Select Value", and select one of the listed companies.

After confirming a vendor has been added to the work order, click "Save Work Order" in the left hand menu.


Since a third party vendor is now associated with the Work Order, our automation script should be triggered. We can confirm this by checking the web application in the "Work Orders" table. If a new work order was successfully stored in the ledger, we should see a new row in the table with the "LastModifiedBy" attribute set as "Maximo", the "Status" set as "WAPPR" (Waiting for Approval), and the ID should match our original Work Order number from Maximo.

We're able to see that here.

<img src="https://i.imgur.com/WtPncoh.png">

Now let's update our work order state, and assign it to a user. We'll do this by clicking the "Update" button in the corresponding row. Enter the user id for a "HAZMAT_INSPECTOR", and a work priority number from 1-4, 4 being the most severe. We should then see the status change from "WAPPR" to "APPR" (Approved), meaning the inspector has verified the affected areas, and the vendor can begin work.

We'll also see the "LastModifiedBy" field update to that same user id that was entered.

<img src="https://i.imgur.com/cYbqQTX.png">

Next, we can assign a vendor. Enter a "HAZMAT_VENDOR" user id, and click "submit". The status will then change from "APPR" to "INPRG" (in progress).


<img src="https://i.imgur.com/9XVsCsv.png">


Finally, we can click update one last time to simulate work being complete and ready for inspection. Enter a user id matching a "BUILDING_INSPECTOR" and click submit. We'll finally see the status change to "COMP" (complete).

<img src="https://i.imgur.com/yubA88j.png">



If we click the "View History" button, we can trace back through the changes made to the work order, and see exactly who was responsible for each change.

<img src="https://i.imgur.com/OonrohQ.png">


# Learn more

<!-- * **Watson IOT Platform Code Patterns**: Enjoyed this Code Pattern? Check out our other [Watson IOT Platform Code Patterns](https://developer.ibm.com/?s=Watson+IOT+Platform). -->

<!-- * **Knowledge Center**:Understand how this Python function can load data into  [Watson IOT Platform Analytics](https://www.ibm.com/support/knowledgecenter/en/SSQP8H/iot/analytics/as_overview.html) -->

# License

This code pattern is licensed under the Apache Software License, Version 2.  Separate third party code objects invoked within this code pattern are licensed by their respective providers pursuant to their own separate licenses. Contributions are subject to the [Developer Certificate of Origin, Version 1.1 (DCO)](https://developercertificate.org/) and the [Apache Software License, Version 2](https://www.apache.org/licenses/LICENSE-2.0.txt).

[Apache Software License (ASL) FAQ](https://www.apache.org/foundation/license-faq.html#WhatDoesItMEAN)
