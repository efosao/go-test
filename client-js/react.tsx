import React from 'react';
import r2wc from "react-to-webcomponent"
import ReactDOM from "react-dom/client";
import Select from 'react-select';

console.log("Initializing React");

type Props = {
    cls?: string;
};

const TestReactComponent = ({ cls = "" }: Props) => {
    return <div className={`p-4 bg-red-50 rounded border-2 border-red-300 ${cls}`}>Test React Component</div>;
};

const TestReactComponentWC = r2wc(TestReactComponent, React, ReactDOM, {
    props: {
        cls: "string",
    },
});

customElements.define("test-rc", TestReactComponentWC);

type Option = {
    label: string;
    value: string;
    selected: boolean;
}

const ReactSelect = ({ options }: { options: string }) => {
    let optsArray: Option[] = []
    
    try {
        optsArray = JSON.parse(options);
    } catch (error) {
        console.log({ error })
    }
    const selectedOptions = optsArray.filter(o => o.selected)

    return (
        <Select
            defaultValue={selectedOptions}
            isMulti
            onChange={selectedOptions => {
                const taggy = document.querySelector("#taggy")
                const htmx = window.__htmx;
                if (!taggy) return;
                htmx.ajax("POST", "/partials/posts/search/0/", {
                    values: {
                        tags: selectedOptions.map(o => o.value).join(","),
                    },
                    target: "#post-list",
                    swap: "outerHTML",
                });

            }}
            name="colors"
            options={optsArray}
            className="basic-multi-select"
            classNamePrefix="select"
        />
    )};


const ReactSelectWC = r2wc(ReactSelect, React, ReactDOM, {
    props: {
        options: "string"
    },
});

customElements.define("react-select", ReactSelectWC);


