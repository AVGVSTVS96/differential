function greet(name: string, greeting: string = "Hello"): string {
  return `${greeting}, ${name}!`;
}

function add(a: number, b: number): number {
  return a + b;
}

function multiply(a: number, b: number): number {
  return a * b;
}

const message = greet("World", "Hi");
console.log(message);
