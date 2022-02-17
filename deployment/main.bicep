param location string = resourceGroup().location
param db_admin string = 'dbadmin'
param db_pwd string =  'Fuerte#123456'
param domain string = 'alemorhubs'
param env string = 'demo'
param short_loc string = 'eus'

// restype-domain-environment-location-instance#
var hubs_namespace_name = 'hubs-${domain}-${env}-${short_loc}'
var hub_name = 'hub-${domain}-location1'
var streaming_job_name = 'stream-${domain}-${env}-${short_loc}'
var sql_name = 'sql-${domain}-${env}-${short_loc}'
var db_name = 'hubdb'
var resource_group_name = 'stream_analytics_consumer_group'

// Event Hubs
resource namespaces_hubs_namespace_resource 'Microsoft.EventHub/namespaces@2021-11-01' = {
  name: hubs_namespace_name
  location: location
  sku: {
    name: 'Standard'
    tier: 'Standard'
    capacity: 1
  }
  properties: {
    disableLocalAuth: false
    zoneRedundant: true
    isAutoInflateEnabled: false
    maximumThroughputUnits: 0
    kafkaEnabled: true
  }
}

resource namespaces_hubs_name_RootManageSharedAccessKey 'Microsoft.EventHub/namespaces/AuthorizationRules@2015-08-01' = {
  parent: namespaces_hubs_namespace_resource
  name: 'RootManageSharedAccessKey'  
  properties: {
    rights: [
      'Listen'
      'Manage'
      'Send'
    ]
  }
}

resource namespaces_hubs_name_default 'Microsoft.EventHub/namespaces/networkRuleSets@2021-11-01' = {
  parent: namespaces_hubs_namespace_resource
  name: 'default'
  properties: {
    publicNetworkAccess: 'Enabled'
    defaultAction: 'Allow'
    virtualNetworkRules: []
    ipRules: []
  }
}

resource namespaces_hubs_name_location1_hub 'Microsoft.EventHub/namespaces/eventhubs@2021-11-01' = {
  parent: namespaces_hubs_namespace_resource
  name: hub_name
  properties: {
    messageRetentionInDays: 1
    partitionCount: 1
    status: 'Active'
  }
}

resource namespaces_hub_read_policy 'Microsoft.EventHub/namespaces/eventhubs/authorizationRules@2021-11-01' = {
  parent: namespaces_hubs_name_location1_hub
  name: '${hub_name}-read-policy'
  properties: {
    rights: [
      'Listen'
    ]
  }
}

resource namespaces_hub_write_policy 'Microsoft.EventHub/namespaces/eventhubs/authorizationRules@2021-11-01' = {
  parent: namespaces_hubs_name_location1_hub
  name: '${hub_name}-write-policy'
  properties: {
    rights: [
      'Send'
    ]
  }
}

resource namespaces_hub_Default_cg 'Microsoft.EventHub/namespaces/eventhubs/consumergroups@2021-11-01' = {
  parent: namespaces_hubs_name_location1_hub
  name: '$Default'
  properties: {}
}

resource namespaces_hub_streamanalytics_cg 'Microsoft.EventHub/namespaces/eventhubs/consumergroups@2021-11-01' = {
  parent: namespaces_hubs_name_location1_hub
  name: resource_group_name
  properties: {}
}

// SQL Server
resource sql_server_resource 'Microsoft.Sql/servers@2021-08-01-preview' = {
  name: sql_name
  location: location
  properties: {
    administratorLogin: db_admin
    administratorLoginPassword: db_pwd
    version: '12.0'
    minimalTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Disabled'
  }
}

resource sqldb_resource 'Microsoft.Sql/servers/databases@2021-08-01-preview' = {
  parent: sql_server_resource
  name: db_name
  location: location
  sku: {
    name: 'Standard'
    tier: 'Standard'
    capacity: 20
  }
  properties: {
    collation: 'SQL_Latin1_General_CP1_CI_AS'
    maxSizeBytes: 268435456000
    catalogCollation: 'SQL_Latin1_General_CP1_CI_AS'
    zoneRedundant: false
    readScale: 'Disabled'
    requestedBackupStorageRedundancy: 'Geo'
    maintenanceConfigurationId: '/subscriptions/97e6e7ea-a213-4f0e-87e0-ea14b9781c76/providers/Microsoft.Maintenance/publicMaintenanceConfigurations/SQL_Default'
    isLedgerOn: false
  }
}

resource sql_Server_AllowAllWindowsAzureIps 'Microsoft.Sql/servers/firewallRules@2021-05-01-preview' = {
  parent: sql_server_resource
  name: 'AllowAllWindowsAzureIps'
  properties: {
    startIpAddress: '0.0.0.0'
    endIpAddress: '0.0.0.0'
  }
}

resource sql_Server_ClientIp_2022_2_12_15_21_5 'Microsoft.Sql/servers/firewallRules@2021-05-01-preview' = {
  parent: sql_server_resource
  name: 'ClientIp-2022-2-12_15-21-5'
  properties: {
    startIpAddress: '76.203.170.106'
    endIpAddress: '76.203.170.106'
  }
}

// Stream Analytics
resource streamingjobs_name_resource 'Microsoft.StreamAnalytics/streamingjobs@2021-10-01-preview' = {
  name: streaming_job_name
  location: location
  properties: {
    sku: {
      name: 'Standard'
    }
    eventsOutOfOrderPolicy: 'Adjust'
    outputErrorPolicy: 'Stop'
    eventsOutOfOrderMaxDelayInSeconds: 0
    eventsLateArrivalMaxDelayInSeconds: 5
    dataLocale: 'en-US'
    compatibilityLevel: '1.2'
    contentStoragePolicy: 'SystemAccount'
    jobType: 'Cloud'
  }
  dependsOn: [
    sql_server_resource
  ]
}

resource streamingjobs_input_alias 'Microsoft.StreamAnalytics/streamingjobs/inputs@2021-10-01-preview' = {
  parent: streamingjobs_name_resource
  name: hub_name
  properties: {
    type: 'Stream'
    datasource: {
      type: 'Microsoft.EventHub/EventHub'
      properties: {
        consumerGroupName: namespaces_hub_streamanalytics_cg.name
        eventHubName: hub_name
        serviceBusNamespace: hub_name
        sharedAccessPolicyName: namespaces_hub_read_policy.name
        authenticationMode: 'ConnectionString'
      }
    }
    compression: {
      type: 'None'
    }
    partitionKey: 'id'
    serialization: {
      type: 'Json'
      properties: {
        encoding: 'UTF8'
      }
    }
  }
}

resource streamingjobs_output_MotorMessages_alias 'Microsoft.StreamAnalytics/streamingjobs/outputs@2021-10-01-preview' = {
  parent: streamingjobs_name_resource
  name: '${db_name}-MotorEvents'
  properties: {
    datasource: {
      type: 'Microsoft.Sql/Server/Database'
      properties: {
        maxWriterCount: 1
        maxBatchCount: 10000
        table: 'MotorEvents'
        server: sql_server_resource.name
        database: sqldb_resource.name
        user: db_admin
        password: db_pwd
        authenticationMode: 'ConnectionString'
      }
    }
  }
}

resource streamingjobs_output_ACMessages_alias 'Microsoft.StreamAnalytics/streamingjobs/outputs@2021-10-01-preview' = {
  parent: streamingjobs_name_resource
  name: '${db_name}-ACEvents'
  properties: {
    datasource: {
      type: 'Microsoft.Sql/Server/Database'
      properties: {
        maxWriterCount: 1
        maxBatchCount: 10000
        table: 'ACEvents'
        server: sql_server_resource.name
        database: sqldb_resource.name
        user: db_admin
        password: db_pwd
        authenticationMode: 'ConnectionString'
      }
    }
  }
}

resource streamingjobs_output_GeneratorMessages_alias 'Microsoft.StreamAnalytics/streamingjobs/outputs@2021-10-01-preview' = {
  parent: streamingjobs_name_resource
  name: '${db_name}-GeneratorEvents'
  properties: {
    datasource: {
      type: 'Microsoft.Sql/Server/Database'
      properties: {
        maxWriterCount: 1
        maxBatchCount: 10000
        table: 'GeneratorEvents'
        server: sql_server_resource.name
        database: sqldb_resource.name
        user: db_admin
        password: db_pwd
        authenticationMode: 'ConnectionString'        
      }
    }
  }
}
