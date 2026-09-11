// Conectar a la base de datos items-api
use items-api;

// Limpiar la colección si existe
db.items.drop();

// Insertar documentos de prueba
db.items.insertMany([
	{
		_id: "1",
		title: "Laptop Gamer Lenovo Legion",
		price: 1500.00
	},
	{
		_id: "2",
		title: "Monitor Curvo Samsung 27\"",
		price: 320.50
	},
	{
		_id: "3",
		title: "Teclado Mecánico Keychron K2",
		price: 110.00
	},
	{
		_id: "4",
		title: "Mouse Inalámbrico Logitech MX Master 3S",
		price: 99.99
	}
]);