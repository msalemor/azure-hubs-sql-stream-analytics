using System;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using Azure.Messaging.EventHubs;
using Azure.Messaging.EventHubs.Producer;
using Domain.Common.Models;
using Microsoft.Azure.WebJobs;
using Microsoft.Azure.WebJobs.Host;
using Microsoft.Extensions.Logging;

namespace GenerateMessages
{
    public static class Function1
    {
        private const string ConnectionString = "Endpoint=sb://hubs-ecloud1-demo-eus.servicebus.windows.net/;SharedAccessKeyName=hub-ecloud1-location1-write-policy;SharedAccessKey=+r4UZWiZvRl6WRsWjMl2Aj7ANwiFKbH+i8ISV+rpLO8=;EntityPath=hub-ecloud1-location1";

        [FunctionName("Function1")]
        public async static Task Run([TimerTrigger("* * * * * *")] TimerInfo myTimer, ILogger log)
        {
            log.LogInformation($"C# Timer trigger function executed at: {DateTime.Now}");
            EventHubProducerClient producerClient = null;

            try
            {
                var rn = (new Random());
                var rnd = rn.Next(1, 4);
                float temperature = rn.Next(20, 101);
                float revolutions = rn.Next(100, 301);
                float voltage = rn.Next(118, 122);
                float hertz = rn.Next(58, 62);
                float amps = rn.Next(20, 31);
                float gasPercentage = rn.Next(50, 101);
                float airTemperature = rn.Next(15,31);
                float airFlow = rn.Next(3, 6);

                object obj = null;

                if (rnd == 1)
                {
                    obj = new MotorEvent { DeviceId = "1X100M", Revolutions = revolutions, Temperature = temperature };
                }
                else if (rnd == 2)
                {
                    
                    obj = new ACEvent { DeviceId = "2X100AC", CoolantTemperature=temperature, AirFlow=airFlow, AirTemperature=airTemperature};
                }
                else
                {
                    obj = new GeneratorEvent { DeviceId = "3X100G", Hertz=hertz, Amps=amps, Voltage= voltage, GasPercentage= gasPercentage };
                }

                // Camel case serialization
                var options = new JsonSerializerOptions
                {
                    PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
                };

                // Serialize into json
                string jsonString = JsonSerializer.Serialize(obj,options);

                log.LogInformation($"Message: {jsonString}");


                // Create a producer client that you can use to send events to an event hub
                producerClient = new EventHubProducerClient(ConnectionString);

                // Create a batch of events 
                EventDataBatch eventBatch = await producerClient.CreateBatchAsync();

                // Emmit the message to Event Hubs
                if (eventBatch.TryAdd(new EventData(Encoding.UTF8.GetBytes(jsonString))))
                    await producerClient.SendAsync(eventBatch);
                else
                    throw new ApplicationException("Unable to create message");
            }
            catch(Exception ex)
            {
                log.LogError(ex, "Unable to write to eventhubs");
            }
            finally
            {
                if (producerClient!=null)
                    await producerClient.DisposeAsync();
            }
        }
    }
}
