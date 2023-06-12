# Busy RDS - A Simulation Tool for Database Workloads

## Introduction
The Busy RDS tool is designed to help engineers simulate various use-cases related to Amazon Relational Database Service (RDS). This application simplifies the process of setting up and testing different features of RDS, enhancing the speed at which engineers can work and innovate.

## Purpose
Amazon RDS provides a plethora of features to enhance the resilience and availability of databases. However, setting up the pre-requisite activities for these features can be time-consuming. The Busy RDS tool helps engineers test these features efficiently, enabling faster development and deployment cycles.

## Key Features
Busy RDS supports the following simulations:

1. **Generating Dummy Data** 
2. **Simulating RDS Failover**

---

## Usage Instructions

### 1. Generating Dummy Data 
This simulation creates dummy data for a specific database server. The process involves creating a new database, building a new table, and then populating it with dummy data based on the `test_run` parameter.

#### Configuration
- **database.dsn** : Connection string to the database.
- **test_run** : Number of dummy data records to create.

#### Example Output
```
...
2023/06/13 01:43:45 Insert: LfErcYv@qiimqdI.ru success
2023/06/13 01:43:46 Insert: FGLjTjq@TkrFxnE.ru success
2023/06/13 01:43:46 Insert: shQoxsk@LaaJdPt.info success
2023/06/13 01:43:47 Insert: IWdGPXH@bindYrA.info success
2023/06/13 01:43:47 Insert: GZRfvXZ@Jnrhsuj.com success
2023/06/13 01:43:48 Insert: TtBVhtI@pGhIvxh.net success
...
```

#### Execution Steps
1. Create a `config.json` file with the required configuration information.
2. Start the application using the command: ```go run main.go -usecase 1```

---

### 2. RDS Failover Simulation
This simulation tests the behavior of Aurora during a failover event. It starts by writing data continuously at a normal day's request per second (RPS) rate. In the event of a failover, the simulation will attempt to retry the insert operation until a connection to the new writer is established.

#### Configuration
- **database.dsn** : Connection string to the database.
- **test_run** : Number of dummy data records to create.
- **rps** : Number of requests per second (typical RPS for the application).
- **max_retry** : Number of times the application should attempt to write in case of failure.
- **delay_retry** : Time the application should wait before retrying.

#### Example Output
```
...
2023/06/13 01:43:55 Insert: EkSJfxm@uPYqRyB.net success
2023/06/13 01:43:56 Insert: uIiJsXW@LuxRPqV.net success
2023/06/13 01:43:58 Failed to insert: fJhsibD@UaHySaY.net. Error: dial tcp 13.228.85.227:5432: connect: connection refused. Retrying (1/20)...
2023/06/13 01:43:59 Failed to insert: fJhsibD@UaHySaY.net. Error: dial tcp 13.228.85.227:5432: connect: connection refused. Retrying (2/20)...
2023/06/13 01:44:00 Failed to insert: fJhsibD@UaHySaY.net. Error: dial tcp 13.228.85.227:5432: connect: connection refused. Retrying (3/20)...
2023/06/13 01:44:01 Failed to insert: fJhsibD@UaHySaY.net. Error: dial tcp 13.228.85.227:5432: connect: connection refused. Retrying (4/20)...
2023/06/13 01:44:02 Failed to insert: fJhsibD@UaHySaY.net. Error: dial tcp 13.228.85.227:5432: connect: connection refused. Retrying (5/20)...
2023/06/13 01:44:03 Failed to insert: fJhsibD@UaHySaY.net. Error: dial tcp 13.228.85.227:5432: connect: connection refused. Retrying (6/20)...
2023/06/13 01:44:04 **DownTime: 6330ms**
2023/06/13 01:44:04 Insert: fJhsibD@UaHySaY.net success
2023/06/13 01:44:05 Insert: ddZxldB@HJEiRLg.info success
...
```

#### Execution Steps
1. Create a `config.json` file with the required configuration information.
2. Start the application using the command: ```go run main.go -usecase 2```
