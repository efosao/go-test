import register from "preact-custom-element";
import { useState } from "preact/hooks";

const Greeting = () => {
  const [count, setCount] = useState(0);
  const increment = () => setCount(count + 1);
  // You can also pass a callback to the setter
  const decrement = () => setCount((currentCount) => currentCount - 1);

  return (
    <div class=" flex flex-col gap-2 bg-pink-200 p-4 my-2 rounded-lg border-2 border-red-300">
      <p class="text-lg font-bold text-pink-800 dark:text-pink-800">
        Count: {count}
      </p>
      <div class="flex gap-2">
        <button class="button" onClick={increment}>
          Increment
        </button>
        <button class="button" onClick={decrement}>
          Decrement
        </button>
      </div>
    </div>
  );
};

register(Greeting, "x-greeting", ["name"], { shadow: false });

export default Greeting;
