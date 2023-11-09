import sys
from hfc.fabric import Client as client_fabric
import asyncio
import base64
import hashlib
from ecdsa import SigningKey, NIST256p
from ecdsa.util import sigencode_der, sigdecode_der
from hfc.fabric_network.gateway import Gateway
from hfc.fabric_network.network import Network
from hfc.fabric_network.contract import Contract

domain = ["connection-org1"]
channel_name = "mychannel"
cc_name = "fabpki"
cc_version = "1.0"
callpeer = ['peer0.org1.example.com']

if __name__ == "__main__":

    #test if the meter creation time was informed as argument

    if len(sys.argv) != 2:
        print("Usage:",sys.argv[0],"<\"ID da estação\">")
        exit(1)

    #get city name
    idEstacao = str(sys.argv[1])

    loop = asyncio.get_event_loop()
    #creates a loop object to manage async transactions
        
    new_gateway = Gateway() # Creates a new gateway instance

    c_hlf = client_fabric(net_profile=(domain[0] + ".json"))
    user = c_hlf.get_user('org1.example.com', 'User1')
    admin = c_hlf.get_user('org1.example.com', 'Admin')
    # print(admin)
    peers = []
    peer = c_hlf.get_peer('peer0.org1.example.com')
    peers.append(peer)
    options = {'wallet': ''}

    c_hlf.new_channel(channel_name)


    #invoke the chaincode to register the meter

    response = loop.run_until_complete(c_hlf.chaincode_invoke(
            requestor=admin,
            channel_name=channel_name,
            peers=['peer0.org1.example.com', 'peer0.org2.example.com'],
            args=[idEstacao],
            fcn= 'getStationData',
            cc_name=cc_name,
            transient_map=None, # optional, for private data
            wait_for_event=True, # for being sure chaincode invocation has been commited in the ledger, default is on tx event
            #cc_pattern='^invoked*' # if you want to wait for chaincode event and you have a `stub.SetEvent("invoked", value)` in your chaincode
            ))
    print(response)