User management in Viridian
===========================

Viridian uses Hyperledger Fabric as its blockchain/database backend. Access to the Hyperledger Fabric blockchain is granted if the user is registered with a Certificate Authority (CA) and possesses an appropriate X509 certificate.

Therefore, user management is done through configuring the CA and determining the attributes (settings) of the certificates.

While it is possible to user third party CAs, we use the Fabric CA that is part of the Hyperledger project and naturally works well with Fabric.

See the [Fabric CA User’s Guide](https://hyperledger-fabric-ca.readthedocs.io/en/release-1.4/users-guide.html) for a description on how to configure and use the Fabric CA.


Attributes in Fabric CA
-----------------------

There are some default attributes in Fabric CA certificates/identities (each identity can have multiple certificates):

- **Name:** (`hf.EnrollmentID` in the certificate or `--id.name` at the Fabric CA client CLI) The name of the identity, i.e. the user name.
- **Type:** (`hf.Type` in the certificate or `--id.type` at the CLI) is one of `peer`, `orderer` and `client`. Each component/participant in Fabric must have a certificate to sign transactions, with its respective type set in this attribute. Peers are the servers that validate transactions, orderers are the servers that create blocks and add them to the blockchain. Client is an application connecting to the blockchain network, i.e. a user.
- **Affiliation:** (`hf.Affiliation` in the certificate or `--id.affiliation` at the CLI) Here, the position of the identity in the organizational hierarchy can be set. It can be '.' for 'root' affiliation (meaning highest possible position in hierarchy, like 'super user'), a certain string, e.g. 'org1', for belonging to organization 1, 'org1.department1' for one level lower, and so on. Affiliations should be lower-case [to avoid confusion](https://hyperledger-fabric-ca.readthedocs.io/en/release-1.4/users-guide.html#registering-a-new-identity).

These three attributes are automatically registered for every identity under the names `hf.EnrollmentID` (=name), `hf.Type`, `hf.Affiliation` (see https://hyperledger-fabric-ca.readthedocs.io/en/release-1.4/users-guide.html#attribute-based-access-control).

Some extra attributes have a special meaning in Fabric CA (reserved attributes, they begin with `hf.`):

- **`hf.Registrar.Roles`:** This is a comma-separated list of types (e.g. peer, client), which can be administered by the identity.
- **`hf.Registrar.Attributes`:** This is a comma-separated list of all attributes (both reserved, i.e. beginning with `hf.`, or not reserved), which can be set by the identity.
- **`hf.Revoker`:** If `true`, the identity is allowed to revoke certificates.
- **`hf.GenCRL`:** If `true`, the identity is allowed to generate a CRL (certificate revocation list).
- **`hf.AffiliationMgr`:** If `true`, the identity is allowed to dynamically add/remove/modify affiliations registered with the server.

On top of that, arbitrary extra attributes can be set as key-value pairs. They should not start with `hf.`.

Each identity, if allowed at all, can only create new identites that have less powers than themselves.


### How to use certificate attributes in Viridian

#### Admin users

There must be at least one identity that is allowed to create and modify other identities. The admin identity might be accessible to everyone under certain conditions, but only through an application that strictly regulates what actions are allowed. Only a very limited number of people (i.e. system administrators of the Fabric CA and the application interfacing it) should have direct access to the admin identity.

It is important to control access to the creation of identities to prevent [Sybil attacks](https://en.wikipedia.org/wiki/Sybil_attack).

#### Registering new users

New identities can only be created once by each non-existent user. New users must prove to the application that they don't already have an account, either manually via communicating with a trusted moderator, or automated if possible. The proof could consist of a presentation of their legal ID/passport, whose number is stored as a hash (how exactly? Salting not possible if lookup is needed?! Use key stretching? https://en.wikipedia.org/wiki/Key_stretching). The hash of the presented ID/passport number is calculated and it is looked up if this hash has been stored already, in which case access is denied.

#### Regular users

All other identities should be only regular users and should not be allowed to create or modify any identities. If a user needs to modify their own identity (but that should rarely or never be required), one may allow them to do so via an interfacing application (see above), but they can never modify any other identity. 

The **type** of all (regular) users should always be **client**.

Of course, setting the **name** in the certificate equal to the user name in the Viridian blockchain is the most natural choice.

