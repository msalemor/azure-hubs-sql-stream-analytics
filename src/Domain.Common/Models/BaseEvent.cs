using System;

namespace Domain.Common.Models
{

    public class BaseEvent
    {
        public BaseEvent()
        {
            Ts = DateTime.UtcNow;
        }

        public string DeviceId { get; set; }
        public string Type { get; set; }
        public DateTime Ts { get; set; }
    }
}
