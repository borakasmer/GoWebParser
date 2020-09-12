using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;
using RabbitMQ.Client;
using Swashbuckle.AspNetCore.Annotations;

namespace GoService.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class ExchangeController : ControllerBase
    {      
        public ExchangeController()
        {
        }

        //dotnet add package RabbitMQ.Client
        //dotnet add package Newtonsoft.Json --version 12.0.3
        [HttpPost]
        [SwaggerOperation(Summary = "ExchangeType is an Enum.", Description = "<h2>ExchangeType Values :</h2> <hr><br> <b>DOLAR : 0</b> </br><b>EURO : 1</b><br><b>STERLİN : 2</b></br><b>GRAM ALTIN : 3")]
        public bool Insert([FromBody] Exchnage exchange)
        {
            try
            {
                exchange.ExchangeName = exchange.ExchageType.ToString();
                Console.WriteLine(exchange.ExchageType);
                var factory = new ConnectionFactory()
                {
                    HostName = "78.217.***.***", //MyIp
                    UserName = "test",
                    Password = "test",
                    Port = 1881,
                    VirtualHost = "/",
                }              
                using (var connection = factory.CreateConnection())
                using (var channel = connection.CreateModel())
                {
                    channel.QueueDeclare(queue: "product",
                                         durable: false,
                                         exclusive: false,
                                         autoDelete: false,
                                         arguments: null);

                    var productData = JsonConvert.SerializeObject(exchange);
                    var body = Encoding.UTF8.GetBytes(productData);

                    channel.BasicPublish(exchange: "",
                                         routingKey: "product",
                                         basicProperties: null,
                                         body: body);
                    Console.WriteLine($"{exchange.Name} is Send to the queue");
                }

                return true;                
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
                return false;
            }                        
        }
    }
}
