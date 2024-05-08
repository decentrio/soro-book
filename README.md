# Soro-book

## Description

Sorobook is an indexing service for stellar soroban platforms. It indexes soroban transactions and soroban events which is served in friendly json-format data via a grpc server.

This repository contains the source code for the data aggregation process:
    1. Fetching ledger data from stellar node
    2. Processing ledger data into usable soroban data
    3. Pushing soroban data into postgres backend

![sorobook-aggregator](https://hackmd.io/_uploads/Bynr9f8MC.jpg)

## How to use

In order to run the sorobook data aggregator, you need to first create an instance of postgresSQL and a stellar node
    - For instructions on how to run the postgresSQL, visit this guide
    - For instructions on how to run the stellar node, see the documentation here

Once that all setup is done, you can start the aggregator by runnining:

```
    run aggregator 
```
