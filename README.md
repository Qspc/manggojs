# FabricNetwork-2.x

Youtube Channel: https://www.youtube.com/watch?v=SJTdJt6N6Ow&list=PLSBNVhWU6KjW4qo1RlmR7cvvV8XIILub6

Network Topology

Three Orgs(Peer Orgs)

    - Each Org have one peer(Each Endorsing Peer)
    - Each Org have separate Certificate Authority
    - Each Peer has Current State database as couch db

One Orderer Org

    - Three Orderers
    - One Certificate Authority

Steps:

1. Clone the repo
2. Run Certificates Authority Services for all Orgs
3. Create Cryptomaterials for all organizations
4. Create Channel Artifacts using Org MSP
5. Create Channel and join peers
6. Deploy Chaincode
   1. Install All dependency
   2. Package Chaincode
   3. Install Chaincode on all Endorsing Peer
   4. Approve Chaincode as per Lifecycle Endorsment Policy
   5. Commit Chaincode Defination
7. Create Connection Profiles
8. Start API Server
9. Register User using API
10. Invoke Chaincode Transaction
11. Query Chaincode Transaction

**installation**

aplikasi :

- viscode
- postman
- npm
- curl
- nodejs
- docker
- python
- go lang
- git

package :

- express
- mongoose
- cors
- jsonwebtoken
- mongodb
- nodemon
- dotenv
- body-parser
- bcrypt
- fabric
- log4js
- setting package golang

nanya diza :

post api login tambah role --> nunggu testing diza
liat kodingan bang dani ---> kyanya local
server linux pak ahsyar --> nunggu resp pak ahsyar
argumen isasset dll masuk ke fe

- key & record
- reject : benihID, kuantitas, prevID, reason

pertanyaan2 yg masih bingung :

- user contract tuh dipake gak si?
- buat sistem dploy kalo local gimna ya?
- di controller invoke ada dua argumen yg tidak dikirim tapi diterima request (userName dan role), itu dapat darimana?
- arg setiap fungsi cuma satu yg dikirim, padhal di kodingan beda2
- informasi role pada fungsi 'GetUserByID' harus dirubah
