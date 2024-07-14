const jsonTypes = ["string", "number", "boolean", /*"object",*/ "array"/*, "null"*/];
const customTypes = ["Main"];

function toPascalCase(str) {
    return str.replace(/(^\w|[a-zA-Z])\S*/g, match => {
        return match.charAt(0).toUpperCase() + match.slice(1).toLowerCase();
    }).replace(/\s+/g, '');
}

function addProperty(button) {
    const typeSection = button.parentElement;
    const propertySection = document.createElement("div");
    propertySection.className = "property-section";
    propertySection.innerHTML = `
        <label>Property name: <input type="text" class="property-name"></label>
        <label>Type: 
            <select class="property-type">
                ${jsonTypes.concat(customTypes).map(type => `<option value="${type}">${type}</option>`).join('')}
            </select>
        </label>
        <label>Array: <input type="checkbox" class="property-array"></label>
        <button onclick="removeProperty(this)">Remove property</button>
    `;
    typeSection.insertBefore(propertySection, button);
}

function addType() {
    const typeName = document.getElementById("new-type-name").value.trim();
    if (typeName && !customTypes.includes(typeName)) {
        customTypes.push(typeName);
        const typeSection = document.createElement("div");
        typeSection.className = "type-section";
        typeSection.setAttribute("data-type", typeName);
        typeSection.innerHTML = `
            <h3>${typeName} <button onclick="removeType(this)">Remove type</button></h3>
            <button onclick="addProperty(this)">Add property</button>
        `;
        document.getElementById("type-definitions").appendChild(typeSection);
        updateTypeOptions();
        document.getElementById("new-type-name").value = '';
    }
}

function removeProperty(button) {
    const propertySection = button.parentElement;
    propertySection.parentElement.removeChild(propertySection);
}

function removeType(button) {
    const typeSection = button.parentElement.parentElement;
    const typeName = typeSection.getAttribute("data-type");
    customTypes.splice(customTypes.indexOf(typeName), 1);
    typeSection.parentElement.removeChild(typeSection);
    updateTypeOptions();
}

function updateTypeOptions() {
    document.querySelectorAll(".property-type").forEach(select => {
        const selectedValue = select.value;
        select.innerHTML = jsonTypes.concat(customTypes).map(type => `<option value="${type}">${type}</option>`).join('');
        select.value = selectedValue;
    });
}

function generateSchema() {
    const schema = {
        "$schema": "http://json-schema.org/draft-06/schema#",
        "$ref": "#/definitions/Main",
        "definitions": {}
    };

    document.querySelectorAll(".type-section").forEach(section => {
        const typeName = section.getAttribute("data-type");
        const pascalTypeName = toPascalCase(typeName);
        schema.definitions[pascalTypeName] = {
            "type": "object",
            "additionalProperties": false,
            "properties": {},
            "required": [],
            "title": pascalTypeName
        };

        section.querySelectorAll(".property-section").forEach(propSection => {
            const propName = propSection.querySelector(".property-name").value;
            const propType = propSection.querySelector(".property-type").value;
            const isArray = propSection.querySelector(".property-array").checked;

            if (propName) {
                let property;
                if (jsonTypes.includes(propType)) {
                    property = isArray ? { "type": "array", "items": { "type": propType } } : { "type": propType };
                } else {
                    const pascalPropType = toPascalCase(propType);
                    property = isArray ? { "type": "array", "items": { "$ref": `#/definitions/${pascalPropType}` } } : { "$ref": `#/definitions/${pascalPropType}` };
                }
                schema.definitions[pascalTypeName].properties[propName] = property;
                schema.definitions[pascalTypeName].required.push(propName);
            }
        });
    });

    document.getElementById("json-schema").textContent = JSON.stringify(schema, null, 2);
}
