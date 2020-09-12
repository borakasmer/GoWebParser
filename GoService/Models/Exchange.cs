using System.ComponentModel.DataAnnotations;

public enum ExchangeType{
    DOLAR,
	EURO,
	STERLİN,
    [Display(Name="GRAM ALTIN")]
	GRAMALTIN
}
public class Exchnage{
    public double Price { get; set; }
    public string Name { get; set; }
    public ExchangeType ExchageType { get; set; }
    public string ExchangeName { get; set; }
}