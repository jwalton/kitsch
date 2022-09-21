import Convert from "ansi-to-html";
import { execSync } from "child_process";
import escapeHtml from "escape-html";
import * as fs from "fs";
import { Code, Parent } from "mdast";
import path from "path";
import tmp from "tmp";
import { Transformer } from "unified";
import visit from "unist-util-visit";

const KITSCH = process.env.KITSCH || "kitsch";
const FLEX_SPACE = "--space--";

const convert = new Convert({
  fg: "#cccccc",
});

const examplePromptRegex = /^<ExamplePrompt>([^]*)<\/ExamplePrompt>$/;

function escapeMdx(text: string): string {
  let result = escapeHtml(text);

  // Replace "{" and "}" because otherwise they will be interpreted as jsx.
  result = result.replace(/([{}])/g, "{'$1'}");

  // Replace " " with "&nbsp;".
  result = result.replace(/ /g, "&nbsp;");

  return result;
}

function execCombined(
  command: string,
  options: { env?: NodeJS.ProcessEnv | undefined }
): string {
  const out = tmp.fileSync({ postfix: ".stdout" });
  execSync(command, { ...options, stdio: ["ignore", out.fd, out.fd] });
  return fs.readFileSync(out.name, { encoding: "utf-8" });
}

function runKitschPrompt(exmaple: string): string {
  const configParts = exmaple.split("---");

  let demoContext: string;
  let kitschConfig: string;

  if (configParts.length === 2) {
    demoContext = configParts[0];
    kitschConfig = configParts[1];
  } else {
    demoContext = "";
    kitschConfig = configParts[0];
  }

  demoContext += `\nflexibleSpaceReplacement: '${FLEX_SPACE}'`;

  // Fix the CWD to be the root of the docs folder.
  kitschConfig = kitschConfig.replace(
    /\${CWD}/g,
    path.join(__dirname, "..", "..")
  );

  // Copy the example to a temporary file
  const demoFile = tmp.fileSync({ postfix: ".yaml" });
  fs.writeFileSync(demoFile.name, demoContext);

  let options = "";
  if (kitschConfig.trim()) {
    const configFile = tmp.fileSync({ postfix: ".yaml" });
    fs.writeFileSync(configFile.name, kitschConfig);
    options += ` --config "${configFile.name}"`;
  }

  // Run kitsch
  try {
    const output = execCombined(
      `${KITSCH} prompt ${options} --demo "${demoFile.name}"`,
      { env: { ...process.env, FORCE_COLOR: "3" } }
    );

    let html = escapeMdx(output);

    // Split up lines at flexible spaces.
    const lines = html.split("\n");
    const splitLines = lines.map((line) => {
      const parts = line.split(FLEX_SPACE);
      if (parts.length > 1) {
        return `<div class="kitschPromptLine"><span class="kitchPromptLinePart">${parts
          .map(convertToHtml)
          .join('</span><span class="kitchPromptLinePart">')}</span></div>`;
      } else {
        return convertToHtml(line);
      }
    });

    html = splitLines.join("");

    return html;
  } catch (err) {
    return escapeMdx("Error running example: " + err.stack);
  }
}

export function convertToHtml(ansi: string) {
  // Convert ANSI to HTML.
  let html = convert.toHtml(ansi);

  // Fix each style tag to be a JSX style tag.
  html = html.replace(/style="([^"]*)"/g, (match, style) => {
    const reactStyle = style
      .split(";")
      .map((s) => {
        const [key, value] = s.split(":");
        return `"${key}": "${value}"`;
      })
      .join(",");

    return `style={{${reactStyle}}}`;
  });
  return html;
}

export default function kitschPromptPreview(): Transformer {
  const transformer: Transformer = async (ast) => {
    // visit(ast, "jsx", (node: Literal, index: number, parent: Parent) => {
    //   const match = examplePromptRegex.exec(node.value);
    //   if (!match) {
    //     return;
    //   }

    //   node.value = `<ExamplePrompt>\n${runKitschPrompt(
    //     match[1]
    //   )}</ExamplePrompt>`;
    // });

    // TODO: Add `import` for ExamplePrompt if it's not already there?

    visit(ast, "code", (node: Code, index: number, parent: Parent) => {
      if (node.lang !== "kitsch") {
        return;
      }

      parent.children[index] = {
        type: "jsx",
        value: `<ExamplePrompt>\n${
          runKitschPrompt(node.value) + "▌"
        }</ExamplePrompt>`,
      } as any;
    });
  };

  return transformer;
}
