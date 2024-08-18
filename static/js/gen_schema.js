const jsonTypes = ["string", "number", "boolean", "object", "array", "null"];
let customTypes = [];

function toPascalCase(str) {
    return str.replace(/(^\w|[a-zA-Z])\S*/g, match => {
        return match.charAt(0).toUpperCase() + match.slice(1).toLowerCase();
    }).replace(/\s+/g, '');
}

function addProperty(button, propName = '', propTypes = [{ type: '', isArray: false }]) {
    const typeSection = button.parentElement;
    const propertySection = document.createElement("div");
    propertySection.className = "property-section";

    propertySection.innerHTML = `
        <label>Property name: <input type="text" class="property-name" value="${propName}"></label>
        <div class="property-types-container"></div>
        <button onclick="addPropertyType(this)">Add another type</button>
        <button onclick="removeProperty(this)">Remove property</button>
    `;

    const typeContainer = propertySection.querySelector(".property-types-container");

    // Add the first type input without a remove button
    addPropertyTypeElement(typeContainer, propTypes[0].type, propTypes[0].isArray, false);

    // Add any additional types with remove buttons
    propTypes.slice(1).forEach(propType => {
        addPropertyTypeElement(typeContainer, propType.type, propType.isArray, true);
    });

    typeSection.appendChild(propertySection);
}

function addPropertyType(button) {
    const propertySection = button.parentElement;
    const typeContainer = propertySection.querySelector(".property-types-container");

    addPropertyTypeElement(typeContainer, '', false, true);
}

function addPropertyTypeElement(container, selectedType = '', isArray = false, canRemove = true) {
    const typeSelect = document.createElement("div");
    typeSelect.className = "property-type-element";
    typeSelect.innerHTML = `
        <select class="property-type">
            ${jsonTypes.concat(customTypes).map(t => `<option value="${t}"${t === selectedType ? ' selected' : ''}>${t}</option>`).join('')}
        </select>
        <label>Array: <input type="checkbox" class="property-array"${isArray ? ' checked' : ''}></label>
        ${canRemove ? `<button onclick="removePropertyType(this)">Remove type</button>` : ''}
    `;
    container.appendChild(typeSelect);
}

function removePropertyType(button) {
    const typeElement = button.parentElement;
    typeElement.parentElement.removeChild(typeElement);
}

function addSchemaType(typeName = '') {
    typeName = typeName || document.getElementById("new-type-name").value.trim();
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
            const propTypes = Array.from(propSection.querySelectorAll(".property-type-element")).map(typeElement => {
                const type = typeElement.querySelector(".property-type").value;
                const isArray = typeElement.querySelector(".property-array").checked;
                return { type, isArray };
            });

            if (propName) {
                const anyOf = propTypes.map(propType => {
                    if (jsonTypes.includes(propType.type)) {
                        return propType.isArray ? { "type": "array", "items": { "type": propType.type } } : { "type": propType.type };
                    } else {
                        const pascalPropType = toPascalCase(propType.type);
                        return propType.isArray ? { "type": "array", "items": { "$ref": `#/definitions/${pascalPropType}` } } : { "$ref": `#/definitions/${pascalPropType}` };
                    }
                });

                schema.definitions[pascalTypeName].properties[propName] = { "anyOf": anyOf };
                schema.definitions[pascalTypeName].required.push(propName);
            }
        });
    });

    const schemaTxt = JSON.stringify(schema, null, 2);
    document.getElementById("json-schema").textContent = schemaTxt;
    return schemaTxt;
}

function initializePageWithSchema(schemaString) {
    const schema = JSON.parse(schemaString);
    const definitions = schema.definitions;

    customTypes = [];

    document.getElementById("type-definitions").innerHTML = '';

    for (const typeName in definitions) {
        if (definitions.hasOwnProperty(typeName)) {
            const originalTypeName = typeName;
            const originalTypeNameCamel = originalTypeName.replace(/([A-Z])/g, ' $1').trim();

            addSchemaType(originalTypeNameCamel);

            const typeSection = document.querySelector(`.type-section[data-type="${originalTypeNameCamel}"]`);
            const properties = definitions[typeName].properties;
            for (const propName in properties) {
                if (properties.hasOwnProperty(propName)) {
                    const prop = properties[propName];
                    console.log(`${prop.anyOf} ${JSON.stringify(prop)}`)
                    if (prop.anyOf != undefined) {
                        const propTypes = prop.anyOf.map(typeDef => {
                            if (typeDef.type === "array") {
                                return { type: typeDef.items.type || typeDef.items.$ref.replace('#/definitions/', ''), isArray: true };
                            } else {
                                return { type: typeDef.type || typeDef.$ref.replace('#/definitions/', ''), isArray: false };
                            }
                        });
                        addProperty(typeSection.querySelector("button"), propName, propTypes);
                    } else {
                        propType = null
                        if (prop.type === "array") {
                            propType = { type: prop.items.type || prop.items.$ref.replace('#/definitions/', ''), isArray: true };
                        } else {
                            propType = { type: prop.type || prop.$ref.replace('#/definitions/', ''), isArray: false };
                        }
                        addProperty(typeSection.querySelector("button"), propName, [propType]);
                    }
                }
            }
        }
    }

    updateTypeOptions();
}
