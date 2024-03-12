import htmx from "htmx.org";
import Select from "react-select";
import r2wc from "@r2wc/react-to-web-component";
// import { MenubarDemo } from "./App";

console.log("Initializing React");

type Props = {
  cls?: string;
};

const TestReactComponent = ({ cls = "" }: Props) => {
  return (
    <div className={`p-4 bg-red-50 rounded border-2 border-red-300 ${cls}`}>
      Test React Component
    </div>
  );
};

const TestReactComponentWC = r2wc(TestReactComponent, {
  props: {
    cls: "string",
  },
});

customElements.define("test-rc", TestReactComponentWC);

type Option = {
  label: string;
  value: string;
  selected: boolean;
};

const ReactSelect = ({ options }: { options: string }) => {
  let optsArray: Option[] = [];

  try {
    optsArray = JSON.parse(options);
  } catch (error) {
    console.error("Error parsing options", error);
  }
  const selectedOptions = optsArray.filter((o) => o.selected);

  return (
    <Select
      defaultValue={selectedOptions}
      isMulti
      onChange={(selectedOptions) => {
        htmx.ajax("POST", "/partials/posts/search/0/", {
          values: {
            tags: selectedOptions.map((o) => o.value).join(","),
          },
          target: "#post-list",
          swap: "innerHTML",
        });
      }}
      name="colors"
      options={optsArray}
      className="basic-multi-select"
      classNamePrefix="select"
    />
  );
};

const ReactSelectWC = r2wc(ReactSelect, {
  props: {
    options: "string",
  },
});

customElements.define("react-select", ReactSelectWC);

// const AppBarWC = r2wc(MenubarDemo);

// customElements.define("app-bar", AppBarWC);
