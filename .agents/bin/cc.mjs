#!/usr/bin/env node
import { createRequire as __ccCreateRequire } from 'node:module'; const require = __ccCreateRequire(import.meta.url);
var __create = Object.create;
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getProtoOf = Object.getPrototypeOf;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __require = /* @__PURE__ */ ((x) => typeof require !== "undefined" ? require : typeof Proxy !== "undefined" ? new Proxy(x, {
  get: (a, b) => (typeof require !== "undefined" ? require : a)[b]
}) : x)(function(x) {
  if (typeof require !== "undefined") return require.apply(this, arguments);
  throw Error('Dynamic require of "' + x + '" is not supported');
});
var __esm = (fn, res, err) => function __init() {
  if (err) throw err[0];
  try {
    return fn && (res = (0, fn[__getOwnPropNames(fn)[0]])(fn = 0)), res;
  } catch (e) {
    throw err = [e], e;
  }
};
var __commonJS = (cb, mod) => function __require2() {
  try {
    return mod || (0, cb[__getOwnPropNames(cb)[0]])((mod = { exports: {} }).exports, mod), mod.exports;
  } catch (e) {
    throw mod = 0, e;
  }
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(
  // If the importer is in node compatibility mode or this is not an ESM
  // file that has been converted to a CommonJS file using a Babel-
  // compatible transform (i.e. "__esModule" has not been set), then set
  // "default" to the CommonJS "module.exports" for node compatibility.
  isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", { value: mod, enumerable: true }) : target,
  mod
));

// <define:__CC_TEMPLATE_INVENTORY__>
var define_CC_TEMPLATE_INVENTORY_default;
var init_define_CC_TEMPLATE_INVENTORY = __esm({
  "<define:__CC_TEMPLATE_INVENTORY__>"() {
    define_CC_TEMPLATE_INVENTORY_default = [".agents/bin/cc.mjs", ".agents/contracts/activity-lifecycle-record.schema.json", ".agents/contracts/closeout-record.schema.json", ".agents/contracts/context-sync-record.schema.json", ".agents/contracts/context-sync-request.schema.json", ".agents/contracts/fake-activity-source.schema.json", ".agents/contracts/merge-confirmation-record.schema.json", ".agents/contracts/plan-draft-request.schema.json", ".agents/contracts/plan-index.schema.json", ".agents/contracts/plan-publication-discovery.schema.json", ".agents/contracts/plan-publication-record.schema.json", ".agents/contracts/plan-work-breakdown.schema.json", ".agents/contracts/product-knowledge-candidate.schema.json", ".agents/contracts/product-knowledge-domain.schema.json", ".agents/contracts/product-knowledge-project.schema.json", ".agents/contracts/product-knowledge-role.schema.json", ".agents/contracts/product-knowledge-workflow.schema.json", ".agents/contracts/review-preparation.schema.json", ".agents/contracts/review-publication-record.schema.json", ".agents/contracts/run-task-request.schema.json", ".agents/contracts/runtime-manifest.schema.json", ".agents/contracts/task-brief.schema.json", ".agents/contracts/verifier-result.schema.json", ".agents/contracts/whats-next-result.schema.json", ".agents/contracts/work-candidate.schema.json", ".agents/contracts/worker-result.schema.json", ".agents/contracts/workspace-bootstrap-request.schema.json", ".agents/contracts/workspace-configure-request.schema.json", ".agents/contracts/workspace.schema.json", ".agents/skills/configure-workspace/SKILL.md", ".agents/skills/configure-workspace/agents/openai.yaml", ".agents/skills/create-plan/SKILL.md", ".agents/skills/create-plan/agents/openai.yaml", ".agents/skills/finish-work/SKILL.md", ".agents/skills/finish-work/agents/openai.yaml", ".agents/skills/gather-context/SKILL.md", ".agents/skills/gather-context/agents/openai.yaml", ".agents/skills/initialize-workspace/SKILL.md", ".agents/skills/initialize-workspace/agents/openai.yaml", ".agents/skills/publish-plan-tasks/SKILL.md", ".agents/skills/publish-plan-tasks/agents/openai.yaml", ".agents/skills/run-task/SKILL.md", ".agents/skills/run-task/agents/openai.yaml", ".agents/skills/sync-context/SKILL.md", ".agents/skills/sync-context/agents/openai.yaml", ".agents/skills/whats-next/SKILL.md", ".agents/skills/whats-next/agents/openai.yaml", ".agents/templates/product-knowledge/GLOSSARY.md", ".agents/templates/product-knowledge/PROJECT.md", ".agents/templates/product-knowledge/domains/domain/README.md", ".agents/templates/product-knowledge/domains/domain/workflows/workflow.md", ".agents/templates/product-knowledge/roles/README.md", ".agents/templates/product-knowledge/roles/role.md", ".claude/commands/configure-workspace.md", ".claude/commands/create-plan.md", ".claude/commands/finish-work.md", ".claude/commands/gather-context.md", ".claude/commands/initialize-workspace.md", ".claude/commands/publish-plan-tasks.md", ".claude/commands/run-task.md", ".claude/commands/sync-context.md", ".claude/commands/whats-next.md", ".codex/skills/configure-workspace/SKILL.md", ".codex/skills/create-plan/SKILL.md", ".codex/skills/finish-work/SKILL.md", ".codex/skills/gather-context/SKILL.md", ".codex/skills/initialize-workspace/SKILL.md", ".codex/skills/publish-plan-tasks/SKILL.md", ".codex/skills/run-task/SKILL.md", ".codex/skills/sync-context/SKILL.md", ".codex/skills/whats-next/SKILL.md", ".gitignore", "AGENTS.md", "CLAUDE.md", "README.md", "WORKFLOW.md", "agents/coordinator.md", "agents/repository-worker.md", "agents/verifier.md", "context/ARCHITECTURE.md", "context/CONVENTIONS.md", "context/DECISIONS.md", "context/PROJECT.md", "context/SOURCES.md", "context/plans/.gitkeep", "context/plans/product-knowledge/0001-overview.md", "context/plans/product-knowledge/0010-requirements.md", "context/plans/product-knowledge/0020-solution.md", "context/plans/product-knowledge/0040-delivery.md", "context/plans/product-knowledge/0050-verification.md", "context/plans/product-knowledge/0070-risks.md", "context/plans/product-knowledge/0080-work-breakdown.md", "context/plans/product-knowledge/README.md", "contributions/general/.gitkeep", "docs/command-reference.md", "docs/configuration.md", "docs/getting-started.md", "docs/using-the-wrapper.md", "template-manifest.json", "workspace.yaml"];
  }
});

// node_modules/ajv/dist/compile/codegen/code.js
var require_code = __commonJS({
  "node_modules/ajv/dist/compile/codegen/code.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.regexpCode = exports.getEsmExportName = exports.getProperty = exports.safeStringify = exports.stringify = exports.strConcat = exports.addCodeArg = exports.str = exports._ = exports.nil = exports._Code = exports.Name = exports.IDENTIFIER = exports._CodeOrName = void 0;
    var _CodeOrName = class {
    };
    exports._CodeOrName = _CodeOrName;
    exports.IDENTIFIER = /^[a-z$_][a-z$_0-9]*$/i;
    var Name = class extends _CodeOrName {
      constructor(s) {
        super();
        if (!exports.IDENTIFIER.test(s))
          throw new Error("CodeGen: name must be a valid identifier");
        this.str = s;
      }
      toString() {
        return this.str;
      }
      emptyStr() {
        return false;
      }
      get names() {
        return { [this.str]: 1 };
      }
    };
    exports.Name = Name;
    var _Code = class extends _CodeOrName {
      constructor(code) {
        super();
        this._items = typeof code === "string" ? [code] : code;
      }
      toString() {
        return this.str;
      }
      emptyStr() {
        if (this._items.length > 1)
          return false;
        const item = this._items[0];
        return item === "" || item === '""';
      }
      get str() {
        var _a;
        return (_a = this._str) !== null && _a !== void 0 ? _a : this._str = this._items.reduce((s, c) => `${s}${c}`, "");
      }
      get names() {
        var _a;
        return (_a = this._names) !== null && _a !== void 0 ? _a : this._names = this._items.reduce((names, c) => {
          if (c instanceof Name)
            names[c.str] = (names[c.str] || 0) + 1;
          return names;
        }, {});
      }
    };
    exports._Code = _Code;
    exports.nil = new _Code("");
    function _(strs, ...args) {
      const code = [strs[0]];
      let i = 0;
      while (i < args.length) {
        addCodeArg(code, args[i]);
        code.push(strs[++i]);
      }
      return new _Code(code);
    }
    exports._ = _;
    var plus = new _Code("+");
    function str(strs, ...args) {
      const expr = [safeStringify(strs[0])];
      let i = 0;
      while (i < args.length) {
        expr.push(plus);
        addCodeArg(expr, args[i]);
        expr.push(plus, safeStringify(strs[++i]));
      }
      optimize(expr);
      return new _Code(expr);
    }
    exports.str = str;
    function addCodeArg(code, arg) {
      if (arg instanceof _Code)
        code.push(...arg._items);
      else if (arg instanceof Name)
        code.push(arg);
      else
        code.push(interpolate(arg));
    }
    exports.addCodeArg = addCodeArg;
    function optimize(expr) {
      let i = 1;
      while (i < expr.length - 1) {
        if (expr[i] === plus) {
          const res = mergeExprItems(expr[i - 1], expr[i + 1]);
          if (res !== void 0) {
            expr.splice(i - 1, 3, res);
            continue;
          }
          expr[i++] = "+";
        }
        i++;
      }
    }
    function mergeExprItems(a, b) {
      if (b === '""')
        return a;
      if (a === '""')
        return b;
      if (typeof a == "string") {
        if (b instanceof Name || a[a.length - 1] !== '"')
          return;
        if (typeof b != "string")
          return `${a.slice(0, -1)}${b}"`;
        if (b[0] === '"')
          return a.slice(0, -1) + b.slice(1);
        return;
      }
      if (typeof b == "string" && b[0] === '"' && !(a instanceof Name))
        return `"${a}${b.slice(1)}`;
      return;
    }
    function strConcat(c1, c2) {
      return c2.emptyStr() ? c1 : c1.emptyStr() ? c2 : str`${c1}${c2}`;
    }
    exports.strConcat = strConcat;
    function interpolate(x) {
      return typeof x == "number" || typeof x == "boolean" || x === null ? x : safeStringify(Array.isArray(x) ? x.join(",") : x);
    }
    function stringify(x) {
      return new _Code(safeStringify(x));
    }
    exports.stringify = stringify;
    function safeStringify(x) {
      return JSON.stringify(x).replace(/\u2028/g, "\\u2028").replace(/\u2029/g, "\\u2029");
    }
    exports.safeStringify = safeStringify;
    function getProperty(key) {
      return typeof key == "string" && exports.IDENTIFIER.test(key) ? new _Code(`.${key}`) : _`[${key}]`;
    }
    exports.getProperty = getProperty;
    function getEsmExportName(key) {
      if (typeof key == "string" && exports.IDENTIFIER.test(key)) {
        return new _Code(`${key}`);
      }
      throw new Error(`CodeGen: invalid export name: ${key}, use explicit $id name mapping`);
    }
    exports.getEsmExportName = getEsmExportName;
    function regexpCode(rx) {
      return new _Code(rx.toString());
    }
    exports.regexpCode = regexpCode;
  }
});

// node_modules/ajv/dist/compile/codegen/scope.js
var require_scope = __commonJS({
  "node_modules/ajv/dist/compile/codegen/scope.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.ValueScope = exports.ValueScopeName = exports.Scope = exports.varKinds = exports.UsedValueState = void 0;
    var code_1 = require_code();
    var ValueError = class extends Error {
      constructor(name) {
        super(`CodeGen: "code" for ${name} not defined`);
        this.value = name.value;
      }
    };
    var UsedValueState;
    (function(UsedValueState2) {
      UsedValueState2[UsedValueState2["Started"] = 0] = "Started";
      UsedValueState2[UsedValueState2["Completed"] = 1] = "Completed";
    })(UsedValueState || (exports.UsedValueState = UsedValueState = {}));
    exports.varKinds = {
      const: new code_1.Name("const"),
      let: new code_1.Name("let"),
      var: new code_1.Name("var")
    };
    var Scope = class {
      constructor({ prefixes, parent } = {}) {
        this._names = {};
        this._prefixes = prefixes;
        this._parent = parent;
      }
      toName(nameOrPrefix) {
        return nameOrPrefix instanceof code_1.Name ? nameOrPrefix : this.name(nameOrPrefix);
      }
      name(prefix) {
        return new code_1.Name(this._newName(prefix));
      }
      _newName(prefix) {
        const ng = this._names[prefix] || this._nameGroup(prefix);
        return `${prefix}${ng.index++}`;
      }
      _nameGroup(prefix) {
        var _a, _b;
        if (((_b = (_a = this._parent) === null || _a === void 0 ? void 0 : _a._prefixes) === null || _b === void 0 ? void 0 : _b.has(prefix)) || this._prefixes && !this._prefixes.has(prefix)) {
          throw new Error(`CodeGen: prefix "${prefix}" is not allowed in this scope`);
        }
        return this._names[prefix] = { prefix, index: 0 };
      }
    };
    exports.Scope = Scope;
    var ValueScopeName = class extends code_1.Name {
      constructor(prefix, nameStr) {
        super(nameStr);
        this.prefix = prefix;
      }
      setValue(value2, { property, itemIndex }) {
        this.value = value2;
        this.scopePath = (0, code_1._)`.${new code_1.Name(property)}[${itemIndex}]`;
      }
    };
    exports.ValueScopeName = ValueScopeName;
    var line = (0, code_1._)`\n`;
    var ValueScope = class extends Scope {
      constructor(opts) {
        super(opts);
        this._values = {};
        this._scope = opts.scope;
        this.opts = { ...opts, _n: opts.lines ? line : code_1.nil };
      }
      get() {
        return this._scope;
      }
      name(prefix) {
        return new ValueScopeName(prefix, this._newName(prefix));
      }
      value(nameOrPrefix, value2) {
        var _a;
        if (value2.ref === void 0)
          throw new Error("CodeGen: ref must be passed in value");
        const name = this.toName(nameOrPrefix);
        const { prefix } = name;
        const valueKey = (_a = value2.key) !== null && _a !== void 0 ? _a : value2.ref;
        let vs = this._values[prefix];
        if (vs) {
          const _name = vs.get(valueKey);
          if (_name)
            return _name;
        } else {
          vs = this._values[prefix] = /* @__PURE__ */ new Map();
        }
        vs.set(valueKey, name);
        const s = this._scope[prefix] || (this._scope[prefix] = []);
        const itemIndex = s.length;
        s[itemIndex] = value2.ref;
        name.setValue(value2, { property: prefix, itemIndex });
        return name;
      }
      getValue(prefix, keyOrRef) {
        const vs = this._values[prefix];
        if (!vs)
          return;
        return vs.get(keyOrRef);
      }
      scopeRefs(scopeName, values20 = this._values) {
        return this._reduceValues(values20, (name) => {
          if (name.scopePath === void 0)
            throw new Error(`CodeGen: name "${name}" has no value`);
          return (0, code_1._)`${scopeName}${name.scopePath}`;
        });
      }
      scopeCode(values20 = this._values, usedValues, getCode) {
        return this._reduceValues(values20, (name) => {
          if (name.value === void 0)
            throw new Error(`CodeGen: name "${name}" has no value`);
          return name.value.code;
        }, usedValues, getCode);
      }
      _reduceValues(values20, valueCode, usedValues = {}, getCode) {
        let code = code_1.nil;
        for (const prefix in values20) {
          const vs = values20[prefix];
          if (!vs)
            continue;
          const nameSet = usedValues[prefix] = usedValues[prefix] || /* @__PURE__ */ new Map();
          vs.forEach((name) => {
            if (nameSet.has(name))
              return;
            nameSet.set(name, UsedValueState.Started);
            let c = valueCode(name);
            if (c) {
              const def = this.opts.es5 ? exports.varKinds.var : exports.varKinds.const;
              code = (0, code_1._)`${code}${def} ${name} = ${c};${this.opts._n}`;
            } else if (c = getCode === null || getCode === void 0 ? void 0 : getCode(name)) {
              code = (0, code_1._)`${code}${c}${this.opts._n}`;
            } else {
              throw new ValueError(name);
            }
            nameSet.set(name, UsedValueState.Completed);
          });
        }
        return code;
      }
    };
    exports.ValueScope = ValueScope;
  }
});

// node_modules/ajv/dist/compile/codegen/index.js
var require_codegen = __commonJS({
  "node_modules/ajv/dist/compile/codegen/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.or = exports.and = exports.not = exports.CodeGen = exports.operators = exports.varKinds = exports.ValueScopeName = exports.ValueScope = exports.Scope = exports.Name = exports.regexpCode = exports.stringify = exports.getProperty = exports.nil = exports.strConcat = exports.str = exports._ = void 0;
    var code_1 = require_code();
    var scope_1 = require_scope();
    var code_2 = require_code();
    Object.defineProperty(exports, "_", { enumerable: true, get: function() {
      return code_2._;
    } });
    Object.defineProperty(exports, "str", { enumerable: true, get: function() {
      return code_2.str;
    } });
    Object.defineProperty(exports, "strConcat", { enumerable: true, get: function() {
      return code_2.strConcat;
    } });
    Object.defineProperty(exports, "nil", { enumerable: true, get: function() {
      return code_2.nil;
    } });
    Object.defineProperty(exports, "getProperty", { enumerable: true, get: function() {
      return code_2.getProperty;
    } });
    Object.defineProperty(exports, "stringify", { enumerable: true, get: function() {
      return code_2.stringify;
    } });
    Object.defineProperty(exports, "regexpCode", { enumerable: true, get: function() {
      return code_2.regexpCode;
    } });
    Object.defineProperty(exports, "Name", { enumerable: true, get: function() {
      return code_2.Name;
    } });
    var scope_2 = require_scope();
    Object.defineProperty(exports, "Scope", { enumerable: true, get: function() {
      return scope_2.Scope;
    } });
    Object.defineProperty(exports, "ValueScope", { enumerable: true, get: function() {
      return scope_2.ValueScope;
    } });
    Object.defineProperty(exports, "ValueScopeName", { enumerable: true, get: function() {
      return scope_2.ValueScopeName;
    } });
    Object.defineProperty(exports, "varKinds", { enumerable: true, get: function() {
      return scope_2.varKinds;
    } });
    exports.operators = {
      GT: new code_1._Code(">"),
      GTE: new code_1._Code(">="),
      LT: new code_1._Code("<"),
      LTE: new code_1._Code("<="),
      EQ: new code_1._Code("==="),
      NEQ: new code_1._Code("!=="),
      NOT: new code_1._Code("!"),
      OR: new code_1._Code("||"),
      AND: new code_1._Code("&&"),
      ADD: new code_1._Code("+")
    };
    var Node = class {
      optimizeNodes() {
        return this;
      }
      optimizeNames(_names, _constants) {
        return this;
      }
    };
    var Def = class extends Node {
      constructor(varKind, name, rhs) {
        super();
        this.varKind = varKind;
        this.name = name;
        this.rhs = rhs;
      }
      render({ es5, _n }) {
        const varKind = es5 ? scope_1.varKinds.var : this.varKind;
        const rhs = this.rhs === void 0 ? "" : ` = ${this.rhs}`;
        return `${varKind} ${this.name}${rhs};` + _n;
      }
      optimizeNames(names, constants) {
        if (!names[this.name.str])
          return;
        if (this.rhs)
          this.rhs = optimizeExpr(this.rhs, names, constants);
        return this;
      }
      get names() {
        return this.rhs instanceof code_1._CodeOrName ? this.rhs.names : {};
      }
    };
    var Assign = class extends Node {
      constructor(lhs, rhs, sideEffects) {
        super();
        this.lhs = lhs;
        this.rhs = rhs;
        this.sideEffects = sideEffects;
      }
      render({ _n }) {
        return `${this.lhs} = ${this.rhs};` + _n;
      }
      optimizeNames(names, constants) {
        if (this.lhs instanceof code_1.Name && !names[this.lhs.str] && !this.sideEffects)
          return;
        this.rhs = optimizeExpr(this.rhs, names, constants);
        return this;
      }
      get names() {
        const names = this.lhs instanceof code_1.Name ? {} : { ...this.lhs.names };
        return addExprNames(names, this.rhs);
      }
    };
    var AssignOp = class extends Assign {
      constructor(lhs, op, rhs, sideEffects) {
        super(lhs, rhs, sideEffects);
        this.op = op;
      }
      render({ _n }) {
        return `${this.lhs} ${this.op}= ${this.rhs};` + _n;
      }
    };
    var Label = class extends Node {
      constructor(label) {
        super();
        this.label = label;
        this.names = {};
      }
      render({ _n }) {
        return `${this.label}:` + _n;
      }
    };
    var Break = class extends Node {
      constructor(label) {
        super();
        this.label = label;
        this.names = {};
      }
      render({ _n }) {
        const label = this.label ? ` ${this.label}` : "";
        return `break${label};` + _n;
      }
    };
    var Throw = class extends Node {
      constructor(error) {
        super();
        this.error = error;
      }
      render({ _n }) {
        return `throw ${this.error};` + _n;
      }
      get names() {
        return this.error.names;
      }
    };
    var AnyCode = class extends Node {
      constructor(code) {
        super();
        this.code = code;
      }
      render({ _n }) {
        return `${this.code};` + _n;
      }
      optimizeNodes() {
        return `${this.code}` ? this : void 0;
      }
      optimizeNames(names, constants) {
        this.code = optimizeExpr(this.code, names, constants);
        return this;
      }
      get names() {
        return this.code instanceof code_1._CodeOrName ? this.code.names : {};
      }
    };
    var ParentNode = class extends Node {
      constructor(nodes = []) {
        super();
        this.nodes = nodes;
      }
      render(opts) {
        return this.nodes.reduce((code, n) => code + n.render(opts), "");
      }
      optimizeNodes() {
        const { nodes } = this;
        let i = nodes.length;
        while (i--) {
          const n = nodes[i].optimizeNodes();
          if (Array.isArray(n))
            nodes.splice(i, 1, ...n);
          else if (n)
            nodes[i] = n;
          else
            nodes.splice(i, 1);
        }
        return nodes.length > 0 ? this : void 0;
      }
      optimizeNames(names, constants) {
        const { nodes } = this;
        let i = nodes.length;
        while (i--) {
          const n = nodes[i];
          if (n.optimizeNames(names, constants))
            continue;
          subtractNames(names, n.names);
          nodes.splice(i, 1);
        }
        return nodes.length > 0 ? this : void 0;
      }
      get names() {
        return this.nodes.reduce((names, n) => addNames(names, n.names), {});
      }
    };
    var BlockNode = class extends ParentNode {
      render(opts) {
        return "{" + opts._n + super.render(opts) + "}" + opts._n;
      }
    };
    var Root = class extends ParentNode {
    };
    var Else = class extends BlockNode {
    };
    Else.kind = "else";
    var If = class _If extends BlockNode {
      constructor(condition, nodes) {
        super(nodes);
        this.condition = condition;
      }
      render(opts) {
        let code = `if(${this.condition})` + super.render(opts);
        if (this.else)
          code += "else " + this.else.render(opts);
        return code;
      }
      optimizeNodes() {
        super.optimizeNodes();
        const cond = this.condition;
        if (cond === true)
          return this.nodes;
        let e = this.else;
        if (e) {
          const ns = e.optimizeNodes();
          e = this.else = Array.isArray(ns) ? new Else(ns) : ns;
        }
        if (e) {
          if (cond === false)
            return e instanceof _If ? e : e.nodes;
          if (this.nodes.length)
            return this;
          return new _If(not(cond), e instanceof _If ? [e] : e.nodes);
        }
        if (cond === false || !this.nodes.length)
          return void 0;
        return this;
      }
      optimizeNames(names, constants) {
        var _a;
        this.else = (_a = this.else) === null || _a === void 0 ? void 0 : _a.optimizeNames(names, constants);
        if (!(super.optimizeNames(names, constants) || this.else))
          return;
        this.condition = optimizeExpr(this.condition, names, constants);
        return this;
      }
      get names() {
        const names = super.names;
        addExprNames(names, this.condition);
        if (this.else)
          addNames(names, this.else.names);
        return names;
      }
    };
    If.kind = "if";
    var For = class extends BlockNode {
    };
    For.kind = "for";
    var ForLoop = class extends For {
      constructor(iteration) {
        super();
        this.iteration = iteration;
      }
      render(opts) {
        return `for(${this.iteration})` + super.render(opts);
      }
      optimizeNames(names, constants) {
        if (!super.optimizeNames(names, constants))
          return;
        this.iteration = optimizeExpr(this.iteration, names, constants);
        return this;
      }
      get names() {
        return addNames(super.names, this.iteration.names);
      }
    };
    var ForRange = class extends For {
      constructor(varKind, name, from, to) {
        super();
        this.varKind = varKind;
        this.name = name;
        this.from = from;
        this.to = to;
      }
      render(opts) {
        const varKind = opts.es5 ? scope_1.varKinds.var : this.varKind;
        const { name, from, to } = this;
        return `for(${varKind} ${name}=${from}; ${name}<${to}; ${name}++)` + super.render(opts);
      }
      get names() {
        const names = addExprNames(super.names, this.from);
        return addExprNames(names, this.to);
      }
    };
    var ForIter = class extends For {
      constructor(loop, varKind, name, iterable) {
        super();
        this.loop = loop;
        this.varKind = varKind;
        this.name = name;
        this.iterable = iterable;
      }
      render(opts) {
        return `for(${this.varKind} ${this.name} ${this.loop} ${this.iterable})` + super.render(opts);
      }
      optimizeNames(names, constants) {
        if (!super.optimizeNames(names, constants))
          return;
        this.iterable = optimizeExpr(this.iterable, names, constants);
        return this;
      }
      get names() {
        return addNames(super.names, this.iterable.names);
      }
    };
    var Func = class extends BlockNode {
      constructor(name, args, async) {
        super();
        this.name = name;
        this.args = args;
        this.async = async;
      }
      render(opts) {
        const _async = this.async ? "async " : "";
        return `${_async}function ${this.name}(${this.args})` + super.render(opts);
      }
    };
    Func.kind = "func";
    var Return = class extends ParentNode {
      render(opts) {
        return "return " + super.render(opts);
      }
    };
    Return.kind = "return";
    var Try = class extends BlockNode {
      render(opts) {
        let code = "try" + super.render(opts);
        if (this.catch)
          code += this.catch.render(opts);
        if (this.finally)
          code += this.finally.render(opts);
        return code;
      }
      optimizeNodes() {
        var _a, _b;
        super.optimizeNodes();
        (_a = this.catch) === null || _a === void 0 ? void 0 : _a.optimizeNodes();
        (_b = this.finally) === null || _b === void 0 ? void 0 : _b.optimizeNodes();
        return this;
      }
      optimizeNames(names, constants) {
        var _a, _b;
        super.optimizeNames(names, constants);
        (_a = this.catch) === null || _a === void 0 ? void 0 : _a.optimizeNames(names, constants);
        (_b = this.finally) === null || _b === void 0 ? void 0 : _b.optimizeNames(names, constants);
        return this;
      }
      get names() {
        const names = super.names;
        if (this.catch)
          addNames(names, this.catch.names);
        if (this.finally)
          addNames(names, this.finally.names);
        return names;
      }
    };
    var Catch = class extends BlockNode {
      constructor(error) {
        super();
        this.error = error;
      }
      render(opts) {
        return `catch(${this.error})` + super.render(opts);
      }
    };
    Catch.kind = "catch";
    var Finally = class extends BlockNode {
      render(opts) {
        return "finally" + super.render(opts);
      }
    };
    Finally.kind = "finally";
    var CodeGen = class {
      constructor(extScope, opts = {}) {
        this._values = {};
        this._blockStarts = [];
        this._constants = {};
        this.opts = { ...opts, _n: opts.lines ? "\n" : "" };
        this._extScope = extScope;
        this._scope = new scope_1.Scope({ parent: extScope });
        this._nodes = [new Root()];
      }
      toString() {
        return this._root.render(this.opts);
      }
      // returns unique name in the internal scope
      name(prefix) {
        return this._scope.name(prefix);
      }
      // reserves unique name in the external scope
      scopeName(prefix) {
        return this._extScope.name(prefix);
      }
      // reserves unique name in the external scope and assigns value to it
      scopeValue(prefixOrName, value2) {
        const name = this._extScope.value(prefixOrName, value2);
        const vs = this._values[name.prefix] || (this._values[name.prefix] = /* @__PURE__ */ new Set());
        vs.add(name);
        return name;
      }
      getScopeValue(prefix, keyOrRef) {
        return this._extScope.getValue(prefix, keyOrRef);
      }
      // return code that assigns values in the external scope to the names that are used internally
      // (same names that were returned by gen.scopeName or gen.scopeValue)
      scopeRefs(scopeName) {
        return this._extScope.scopeRefs(scopeName, this._values);
      }
      scopeCode() {
        return this._extScope.scopeCode(this._values);
      }
      _def(varKind, nameOrPrefix, rhs, constant) {
        const name = this._scope.toName(nameOrPrefix);
        if (rhs !== void 0 && constant)
          this._constants[name.str] = rhs;
        this._leafNode(new Def(varKind, name, rhs));
        return name;
      }
      // `const` declaration (`var` in es5 mode)
      const(nameOrPrefix, rhs, _constant) {
        return this._def(scope_1.varKinds.const, nameOrPrefix, rhs, _constant);
      }
      // `let` declaration with optional assignment (`var` in es5 mode)
      let(nameOrPrefix, rhs, _constant) {
        return this._def(scope_1.varKinds.let, nameOrPrefix, rhs, _constant);
      }
      // `var` declaration with optional assignment
      var(nameOrPrefix, rhs, _constant) {
        return this._def(scope_1.varKinds.var, nameOrPrefix, rhs, _constant);
      }
      // assignment code
      assign(lhs, rhs, sideEffects) {
        return this._leafNode(new Assign(lhs, rhs, sideEffects));
      }
      // `+=` code
      add(lhs, rhs) {
        return this._leafNode(new AssignOp(lhs, exports.operators.ADD, rhs));
      }
      // appends passed SafeExpr to code or executes Block
      code(c) {
        if (typeof c == "function")
          c();
        else if (c !== code_1.nil)
          this._leafNode(new AnyCode(c));
        return this;
      }
      // returns code for object literal for the passed argument list of key-value pairs
      object(...keyValues) {
        const code = ["{"];
        for (const [key, value2] of keyValues) {
          if (code.length > 1)
            code.push(",");
          code.push(key);
          if (key !== value2 || this.opts.es5) {
            code.push(":");
            (0, code_1.addCodeArg)(code, value2);
          }
        }
        code.push("}");
        return new code_1._Code(code);
      }
      // `if` clause (or statement if `thenBody` and, optionally, `elseBody` are passed)
      if(condition, thenBody, elseBody) {
        this._blockNode(new If(condition));
        if (thenBody && elseBody) {
          this.code(thenBody).else().code(elseBody).endIf();
        } else if (thenBody) {
          this.code(thenBody).endIf();
        } else if (elseBody) {
          throw new Error('CodeGen: "else" body without "then" body');
        }
        return this;
      }
      // `else if` clause - invalid without `if` or after `else` clauses
      elseIf(condition) {
        return this._elseNode(new If(condition));
      }
      // `else` clause - only valid after `if` or `else if` clauses
      else() {
        return this._elseNode(new Else());
      }
      // end `if` statement (needed if gen.if was used only with condition)
      endIf() {
        return this._endBlockNode(If, Else);
      }
      _for(node, forBody) {
        this._blockNode(node);
        if (forBody)
          this.code(forBody).endFor();
        return this;
      }
      // a generic `for` clause (or statement if `forBody` is passed)
      for(iteration, forBody) {
        return this._for(new ForLoop(iteration), forBody);
      }
      // `for` statement for a range of values
      forRange(nameOrPrefix, from, to, forBody, varKind = this.opts.es5 ? scope_1.varKinds.var : scope_1.varKinds.let) {
        const name = this._scope.toName(nameOrPrefix);
        return this._for(new ForRange(varKind, name, from, to), () => forBody(name));
      }
      // `for-of` statement (in es5 mode replace with a normal for loop)
      forOf(nameOrPrefix, iterable, forBody, varKind = scope_1.varKinds.const) {
        const name = this._scope.toName(nameOrPrefix);
        if (this.opts.es5) {
          const arr = iterable instanceof code_1.Name ? iterable : this.var("_arr", iterable);
          return this.forRange("_i", 0, (0, code_1._)`${arr}.length`, (i) => {
            this.var(name, (0, code_1._)`${arr}[${i}]`);
            forBody(name);
          });
        }
        return this._for(new ForIter("of", varKind, name, iterable), () => forBody(name));
      }
      // `for-in` statement.
      // With option `ownProperties` replaced with a `for-of` loop for object keys
      forIn(nameOrPrefix, obj, forBody, varKind = this.opts.es5 ? scope_1.varKinds.var : scope_1.varKinds.const) {
        if (this.opts.ownProperties) {
          return this.forOf(nameOrPrefix, (0, code_1._)`Object.keys(${obj})`, forBody);
        }
        const name = this._scope.toName(nameOrPrefix);
        return this._for(new ForIter("in", varKind, name, obj), () => forBody(name));
      }
      // end `for` loop
      endFor() {
        return this._endBlockNode(For);
      }
      // `label` statement
      label(label) {
        return this._leafNode(new Label(label));
      }
      // `break` statement
      break(label) {
        return this._leafNode(new Break(label));
      }
      // `return` statement
      return(value2) {
        const node = new Return();
        this._blockNode(node);
        this.code(value2);
        if (node.nodes.length !== 1)
          throw new Error('CodeGen: "return" should have one node');
        return this._endBlockNode(Return);
      }
      // `try` statement
      try(tryBody, catchCode, finallyCode) {
        if (!catchCode && !finallyCode)
          throw new Error('CodeGen: "try" without "catch" and "finally"');
        const node = new Try();
        this._blockNode(node);
        this.code(tryBody);
        if (catchCode) {
          const error = this.name("e");
          this._currNode = node.catch = new Catch(error);
          catchCode(error);
        }
        if (finallyCode) {
          this._currNode = node.finally = new Finally();
          this.code(finallyCode);
        }
        return this._endBlockNode(Catch, Finally);
      }
      // `throw` statement
      throw(error) {
        return this._leafNode(new Throw(error));
      }
      // start self-balancing block
      block(body, nodeCount) {
        this._blockStarts.push(this._nodes.length);
        if (body)
          this.code(body).endBlock(nodeCount);
        return this;
      }
      // end the current self-balancing block
      endBlock(nodeCount) {
        const len = this._blockStarts.pop();
        if (len === void 0)
          throw new Error("CodeGen: not in self-balancing block");
        const toClose = this._nodes.length - len;
        if (toClose < 0 || nodeCount !== void 0 && toClose !== nodeCount) {
          throw new Error(`CodeGen: wrong number of nodes: ${toClose} vs ${nodeCount} expected`);
        }
        this._nodes.length = len;
        return this;
      }
      // `function` heading (or definition if funcBody is passed)
      func(name, args = code_1.nil, async, funcBody) {
        this._blockNode(new Func(name, args, async));
        if (funcBody)
          this.code(funcBody).endFunc();
        return this;
      }
      // end function definition
      endFunc() {
        return this._endBlockNode(Func);
      }
      optimize(n = 1) {
        while (n-- > 0) {
          this._root.optimizeNodes();
          this._root.optimizeNames(this._root.names, this._constants);
        }
      }
      _leafNode(node) {
        this._currNode.nodes.push(node);
        return this;
      }
      _blockNode(node) {
        this._currNode.nodes.push(node);
        this._nodes.push(node);
      }
      _endBlockNode(N1, N2) {
        const n = this._currNode;
        if (n instanceof N1 || N2 && n instanceof N2) {
          this._nodes.pop();
          return this;
        }
        throw new Error(`CodeGen: not in block "${N2 ? `${N1.kind}/${N2.kind}` : N1.kind}"`);
      }
      _elseNode(node) {
        const n = this._currNode;
        if (!(n instanceof If)) {
          throw new Error('CodeGen: "else" without "if"');
        }
        this._currNode = n.else = node;
        return this;
      }
      get _root() {
        return this._nodes[0];
      }
      get _currNode() {
        const ns = this._nodes;
        return ns[ns.length - 1];
      }
      set _currNode(node) {
        const ns = this._nodes;
        ns[ns.length - 1] = node;
      }
    };
    exports.CodeGen = CodeGen;
    function addNames(names, from) {
      for (const n in from)
        names[n] = (names[n] || 0) + (from[n] || 0);
      return names;
    }
    function addExprNames(names, from) {
      return from instanceof code_1._CodeOrName ? addNames(names, from.names) : names;
    }
    function optimizeExpr(expr, names, constants) {
      if (expr instanceof code_1.Name)
        return replaceName(expr);
      if (!canOptimize(expr))
        return expr;
      return new code_1._Code(expr._items.reduce((items, c) => {
        if (c instanceof code_1.Name)
          c = replaceName(c);
        if (c instanceof code_1._Code)
          items.push(...c._items);
        else
          items.push(c);
        return items;
      }, []));
      function replaceName(n) {
        const c = constants[n.str];
        if (c === void 0 || names[n.str] !== 1)
          return n;
        delete names[n.str];
        return c;
      }
      function canOptimize(e) {
        return e instanceof code_1._Code && e._items.some((c) => c instanceof code_1.Name && names[c.str] === 1 && constants[c.str] !== void 0);
      }
    }
    function subtractNames(names, from) {
      for (const n in from)
        names[n] = (names[n] || 0) - (from[n] || 0);
    }
    function not(x) {
      return typeof x == "boolean" || typeof x == "number" || x === null ? !x : (0, code_1._)`!${par(x)}`;
    }
    exports.not = not;
    var andCode = mappend(exports.operators.AND);
    function and(...args) {
      return args.reduce(andCode);
    }
    exports.and = and;
    var orCode = mappend(exports.operators.OR);
    function or(...args) {
      return args.reduce(orCode);
    }
    exports.or = or;
    function mappend(op) {
      return (x, y) => x === code_1.nil ? y : y === code_1.nil ? x : (0, code_1._)`${par(x)} ${op} ${par(y)}`;
    }
    function par(x) {
      return x instanceof code_1.Name ? x : (0, code_1._)`(${x})`;
    }
  }
});

// node_modules/ajv/dist/compile/util.js
var require_util = __commonJS({
  "node_modules/ajv/dist/compile/util.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.checkStrictMode = exports.getErrorPath = exports.Type = exports.useFunc = exports.setEvaluated = exports.evaluatedPropsToName = exports.mergeEvaluated = exports.eachItem = exports.unescapeJsonPointer = exports.escapeJsonPointer = exports.escapeFragment = exports.unescapeFragment = exports.schemaRefOrVal = exports.schemaHasRulesButRef = exports.schemaHasRules = exports.checkUnknownRules = exports.alwaysValidSchema = exports.toHash = void 0;
    var codegen_1 = require_codegen();
    var code_1 = require_code();
    function toHash(arr) {
      const hash = {};
      for (const item of arr)
        hash[item] = true;
      return hash;
    }
    exports.toHash = toHash;
    function alwaysValidSchema(it, schema2) {
      if (typeof schema2 == "boolean")
        return schema2;
      if (Object.keys(schema2).length === 0)
        return true;
      checkUnknownRules(it, schema2);
      return !schemaHasRules(schema2, it.self.RULES.all);
    }
    exports.alwaysValidSchema = alwaysValidSchema;
    function checkUnknownRules(it, schema2 = it.schema) {
      const { opts, self } = it;
      if (!opts.strictSchema)
        return;
      if (typeof schema2 === "boolean")
        return;
      const rules = self.RULES.keywords;
      for (const key in schema2) {
        if (!rules[key])
          checkStrictMode(it, `unknown keyword: "${key}"`);
      }
    }
    exports.checkUnknownRules = checkUnknownRules;
    function schemaHasRules(schema2, rules) {
      if (typeof schema2 == "boolean")
        return !schema2;
      for (const key in schema2)
        if (rules[key])
          return true;
      return false;
    }
    exports.schemaHasRules = schemaHasRules;
    function schemaHasRulesButRef(schema2, RULES) {
      if (typeof schema2 == "boolean")
        return !schema2;
      for (const key in schema2)
        if (key !== "$ref" && RULES.all[key])
          return true;
      return false;
    }
    exports.schemaHasRulesButRef = schemaHasRulesButRef;
    function schemaRefOrVal({ topSchemaRef, schemaPath }, schema2, keyword, $data) {
      if (!$data) {
        if (typeof schema2 == "number" || typeof schema2 == "boolean")
          return schema2;
        if (typeof schema2 == "string")
          return (0, codegen_1._)`${schema2}`;
      }
      return (0, codegen_1._)`${topSchemaRef}${schemaPath}${(0, codegen_1.getProperty)(keyword)}`;
    }
    exports.schemaRefOrVal = schemaRefOrVal;
    function unescapeFragment(str) {
      return unescapeJsonPointer(decodeURIComponent(str));
    }
    exports.unescapeFragment = unescapeFragment;
    function escapeFragment(str) {
      return encodeURIComponent(escapeJsonPointer(str));
    }
    exports.escapeFragment = escapeFragment;
    function escapeJsonPointer(str) {
      if (typeof str == "number")
        return `${str}`;
      return str.replace(/~/g, "~0").replace(/\//g, "~1");
    }
    exports.escapeJsonPointer = escapeJsonPointer;
    function unescapeJsonPointer(str) {
      return str.replace(/~1/g, "/").replace(/~0/g, "~");
    }
    exports.unescapeJsonPointer = unescapeJsonPointer;
    function eachItem(xs, f) {
      if (Array.isArray(xs)) {
        for (const x of xs)
          f(x);
      } else {
        f(xs);
      }
    }
    exports.eachItem = eachItem;
    function makeMergeEvaluated({ mergeNames, mergeToName, mergeValues, resultToName }) {
      return (gen, from, to, toName) => {
        const res = to === void 0 ? from : to instanceof codegen_1.Name ? (from instanceof codegen_1.Name ? mergeNames(gen, from, to) : mergeToName(gen, from, to), to) : from instanceof codegen_1.Name ? (mergeToName(gen, to, from), from) : mergeValues(from, to);
        return toName === codegen_1.Name && !(res instanceof codegen_1.Name) ? resultToName(gen, res) : res;
      };
    }
    exports.mergeEvaluated = {
      props: makeMergeEvaluated({
        mergeNames: (gen, from, to) => gen.if((0, codegen_1._)`${to} !== true && ${from} !== undefined`, () => {
          gen.if((0, codegen_1._)`${from} === true`, () => gen.assign(to, true), () => gen.assign(to, (0, codegen_1._)`${to} || {}`).code((0, codegen_1._)`Object.assign(${to}, ${from})`));
        }),
        mergeToName: (gen, from, to) => gen.if((0, codegen_1._)`${to} !== true`, () => {
          if (from === true) {
            gen.assign(to, true);
          } else {
            gen.assign(to, (0, codegen_1._)`${to} || {}`);
            setEvaluated(gen, to, from);
          }
        }),
        mergeValues: (from, to) => from === true ? true : { ...from, ...to },
        resultToName: evaluatedPropsToName
      }),
      items: makeMergeEvaluated({
        mergeNames: (gen, from, to) => gen.if((0, codegen_1._)`${to} !== true && ${from} !== undefined`, () => gen.assign(to, (0, codegen_1._)`${from} === true ? true : ${to} > ${from} ? ${to} : ${from}`)),
        mergeToName: (gen, from, to) => gen.if((0, codegen_1._)`${to} !== true`, () => gen.assign(to, from === true ? true : (0, codegen_1._)`${to} > ${from} ? ${to} : ${from}`)),
        mergeValues: (from, to) => from === true ? true : Math.max(from, to),
        resultToName: (gen, items) => gen.var("items", items)
      })
    };
    function evaluatedPropsToName(gen, ps) {
      if (ps === true)
        return gen.var("props", true);
      const props = gen.var("props", (0, codegen_1._)`{}`);
      if (ps !== void 0)
        setEvaluated(gen, props, ps);
      return props;
    }
    exports.evaluatedPropsToName = evaluatedPropsToName;
    function setEvaluated(gen, props, ps) {
      Object.keys(ps).forEach((p) => gen.assign((0, codegen_1._)`${props}${(0, codegen_1.getProperty)(p)}`, true));
    }
    exports.setEvaluated = setEvaluated;
    var snippets = {};
    function useFunc(gen, f) {
      return gen.scopeValue("func", {
        ref: f,
        code: snippets[f.code] || (snippets[f.code] = new code_1._Code(f.code))
      });
    }
    exports.useFunc = useFunc;
    var Type;
    (function(Type2) {
      Type2[Type2["Num"] = 0] = "Num";
      Type2[Type2["Str"] = 1] = "Str";
    })(Type || (exports.Type = Type = {}));
    function getErrorPath(dataProp, dataPropType, jsPropertySyntax) {
      if (dataProp instanceof codegen_1.Name) {
        const isNumber = dataPropType === Type.Num;
        return jsPropertySyntax ? isNumber ? (0, codegen_1._)`"[" + ${dataProp} + "]"` : (0, codegen_1._)`"['" + ${dataProp} + "']"` : isNumber ? (0, codegen_1._)`"/" + ${dataProp}` : (0, codegen_1._)`"/" + ${dataProp}.replace(/~/g, "~0").replace(/\\//g, "~1")`;
      }
      return jsPropertySyntax ? (0, codegen_1.getProperty)(dataProp).toString() : "/" + escapeJsonPointer(dataProp);
    }
    exports.getErrorPath = getErrorPath;
    function checkStrictMode(it, msg, mode = it.opts.strictSchema) {
      if (!mode)
        return;
      msg = `strict mode: ${msg}`;
      if (mode === true)
        throw new Error(msg);
      it.self.logger.warn(msg);
    }
    exports.checkStrictMode = checkStrictMode;
  }
});

// node_modules/ajv/dist/compile/names.js
var require_names = __commonJS({
  "node_modules/ajv/dist/compile/names.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var names = {
      // validation function arguments
      data: new codegen_1.Name("data"),
      // data passed to validation function
      // args passed from referencing schema
      valCxt: new codegen_1.Name("valCxt"),
      // validation/data context - should not be used directly, it is destructured to the names below
      instancePath: new codegen_1.Name("instancePath"),
      parentData: new codegen_1.Name("parentData"),
      parentDataProperty: new codegen_1.Name("parentDataProperty"),
      rootData: new codegen_1.Name("rootData"),
      // root data - same as the data passed to the first/top validation function
      dynamicAnchors: new codegen_1.Name("dynamicAnchors"),
      // used to support recursiveRef and dynamicRef
      // function scoped variables
      vErrors: new codegen_1.Name("vErrors"),
      // null or array of validation errors
      errors: new codegen_1.Name("errors"),
      // counter of validation errors
      this: new codegen_1.Name("this"),
      // "globals"
      self: new codegen_1.Name("self"),
      scope: new codegen_1.Name("scope"),
      // JTD serialize/parse name for JSON string and position
      json: new codegen_1.Name("json"),
      jsonPos: new codegen_1.Name("jsonPos"),
      jsonLen: new codegen_1.Name("jsonLen"),
      jsonPart: new codegen_1.Name("jsonPart")
    };
    exports.default = names;
  }
});

// node_modules/ajv/dist/compile/errors.js
var require_errors = __commonJS({
  "node_modules/ajv/dist/compile/errors.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.extendErrors = exports.resetErrorsCount = exports.reportExtraError = exports.reportError = exports.keyword$DataError = exports.keywordError = void 0;
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var names_1 = require_names();
    exports.keywordError = {
      message: ({ keyword }) => (0, codegen_1.str)`must pass "${keyword}" keyword validation`
    };
    exports.keyword$DataError = {
      message: ({ keyword, schemaType }) => schemaType ? (0, codegen_1.str)`"${keyword}" keyword must be ${schemaType} ($data)` : (0, codegen_1.str)`"${keyword}" keyword is invalid ($data)`
    };
    function reportError(cxt, error = exports.keywordError, errorPaths, overrideAllErrors) {
      const { it } = cxt;
      const { gen, compositeRule, allErrors } = it;
      const errObj = errorObjectCode(cxt, error, errorPaths);
      if (overrideAllErrors !== null && overrideAllErrors !== void 0 ? overrideAllErrors : compositeRule || allErrors) {
        addError(gen, errObj);
      } else {
        returnErrors(it, (0, codegen_1._)`[${errObj}]`);
      }
    }
    exports.reportError = reportError;
    function reportExtraError(cxt, error = exports.keywordError, errorPaths) {
      const { it } = cxt;
      const { gen, compositeRule, allErrors } = it;
      const errObj = errorObjectCode(cxt, error, errorPaths);
      addError(gen, errObj);
      if (!(compositeRule || allErrors)) {
        returnErrors(it, names_1.default.vErrors);
      }
    }
    exports.reportExtraError = reportExtraError;
    function resetErrorsCount(gen, errsCount) {
      gen.assign(names_1.default.errors, errsCount);
      gen.if((0, codegen_1._)`${names_1.default.vErrors} !== null`, () => gen.if(errsCount, () => gen.assign((0, codegen_1._)`${names_1.default.vErrors}.length`, errsCount), () => gen.assign(names_1.default.vErrors, null)));
    }
    exports.resetErrorsCount = resetErrorsCount;
    function extendErrors({ gen, keyword, schemaValue, data, errsCount, it }) {
      if (errsCount === void 0)
        throw new Error("ajv implementation error");
      const err = gen.name("err");
      gen.forRange("i", errsCount, names_1.default.errors, (i) => {
        gen.const(err, (0, codegen_1._)`${names_1.default.vErrors}[${i}]`);
        gen.if((0, codegen_1._)`${err}.instancePath === undefined`, () => gen.assign((0, codegen_1._)`${err}.instancePath`, (0, codegen_1.strConcat)(names_1.default.instancePath, it.errorPath)));
        gen.assign((0, codegen_1._)`${err}.schemaPath`, (0, codegen_1.str)`${it.errSchemaPath}/${keyword}`);
        if (it.opts.verbose) {
          gen.assign((0, codegen_1._)`${err}.schema`, schemaValue);
          gen.assign((0, codegen_1._)`${err}.data`, data);
        }
      });
    }
    exports.extendErrors = extendErrors;
    function addError(gen, errObj) {
      const err = gen.const("err", errObj);
      gen.if((0, codegen_1._)`${names_1.default.vErrors} === null`, () => gen.assign(names_1.default.vErrors, (0, codegen_1._)`[${err}]`), (0, codegen_1._)`${names_1.default.vErrors}.push(${err})`);
      gen.code((0, codegen_1._)`${names_1.default.errors}++`);
    }
    function returnErrors(it, errs) {
      const { gen, validateName, schemaEnv } = it;
      if (schemaEnv.$async) {
        gen.throw((0, codegen_1._)`new ${it.ValidationError}(${errs})`);
      } else {
        gen.assign((0, codegen_1._)`${validateName}.errors`, errs);
        gen.return(false);
      }
    }
    var E = {
      keyword: new codegen_1.Name("keyword"),
      schemaPath: new codegen_1.Name("schemaPath"),
      // also used in JTD errors
      params: new codegen_1.Name("params"),
      propertyName: new codegen_1.Name("propertyName"),
      message: new codegen_1.Name("message"),
      schema: new codegen_1.Name("schema"),
      parentSchema: new codegen_1.Name("parentSchema")
    };
    function errorObjectCode(cxt, error, errorPaths) {
      const { createErrors } = cxt.it;
      if (createErrors === false)
        return (0, codegen_1._)`{}`;
      return errorObject(cxt, error, errorPaths);
    }
    function errorObject(cxt, error, errorPaths = {}) {
      const { gen, it } = cxt;
      const keyValues = [
        errorInstancePath(it, errorPaths),
        errorSchemaPath(cxt, errorPaths)
      ];
      extraErrorProps(cxt, error, keyValues);
      return gen.object(...keyValues);
    }
    function errorInstancePath({ errorPath }, { instancePath }) {
      const instPath = instancePath ? (0, codegen_1.str)`${errorPath}${(0, util_1.getErrorPath)(instancePath, util_1.Type.Str)}` : errorPath;
      return [names_1.default.instancePath, (0, codegen_1.strConcat)(names_1.default.instancePath, instPath)];
    }
    function errorSchemaPath({ keyword, it: { errSchemaPath } }, { schemaPath, parentSchema }) {
      let schPath = parentSchema ? errSchemaPath : (0, codegen_1.str)`${errSchemaPath}/${keyword}`;
      if (schemaPath) {
        schPath = (0, codegen_1.str)`${schPath}${(0, util_1.getErrorPath)(schemaPath, util_1.Type.Str)}`;
      }
      return [E.schemaPath, schPath];
    }
    function extraErrorProps(cxt, { params, message }, keyValues) {
      const { keyword, data, schemaValue, it } = cxt;
      const { opts, propertyName, topSchemaRef, schemaPath } = it;
      keyValues.push([E.keyword, keyword], [E.params, typeof params == "function" ? params(cxt) : params || (0, codegen_1._)`{}`]);
      if (opts.messages) {
        keyValues.push([E.message, typeof message == "function" ? message(cxt) : message]);
      }
      if (opts.verbose) {
        keyValues.push([E.schema, schemaValue], [E.parentSchema, (0, codegen_1._)`${topSchemaRef}${schemaPath}`], [names_1.default.data, data]);
      }
      if (propertyName)
        keyValues.push([E.propertyName, propertyName]);
    }
  }
});

// node_modules/ajv/dist/compile/validate/boolSchema.js
var require_boolSchema = __commonJS({
  "node_modules/ajv/dist/compile/validate/boolSchema.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.boolOrEmptySchema = exports.topBoolOrEmptySchema = void 0;
    var errors_1 = require_errors();
    var codegen_1 = require_codegen();
    var names_1 = require_names();
    var boolError = {
      message: "boolean schema is false"
    };
    function topBoolOrEmptySchema(it) {
      const { gen, schema: schema2, validateName } = it;
      if (schema2 === false) {
        falseSchemaError(it, false);
      } else if (typeof schema2 == "object" && schema2.$async === true) {
        gen.return(names_1.default.data);
      } else {
        gen.assign((0, codegen_1._)`${validateName}.errors`, null);
        gen.return(true);
      }
    }
    exports.topBoolOrEmptySchema = topBoolOrEmptySchema;
    function boolOrEmptySchema(it, valid) {
      const { gen, schema: schema2 } = it;
      if (schema2 === false) {
        gen.var(valid, false);
        falseSchemaError(it);
      } else {
        gen.var(valid, true);
      }
    }
    exports.boolOrEmptySchema = boolOrEmptySchema;
    function falseSchemaError(it, overrideAllErrors) {
      const { gen, data } = it;
      const cxt = {
        gen,
        keyword: "false schema",
        data,
        schema: false,
        schemaCode: false,
        schemaValue: false,
        params: {},
        it
      };
      (0, errors_1.reportError)(cxt, boolError, void 0, overrideAllErrors);
    }
  }
});

// node_modules/ajv/dist/compile/rules.js
var require_rules = __commonJS({
  "node_modules/ajv/dist/compile/rules.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.getRules = exports.isJSONType = void 0;
    var _jsonTypes = ["string", "number", "integer", "boolean", "null", "object", "array"];
    var jsonTypes = new Set(_jsonTypes);
    function isJSONType(x) {
      return typeof x == "string" && jsonTypes.has(x);
    }
    exports.isJSONType = isJSONType;
    function getRules() {
      const groups = {
        number: { type: "number", rules: [] },
        string: { type: "string", rules: [] },
        array: { type: "array", rules: [] },
        object: { type: "object", rules: [] }
      };
      return {
        types: { ...groups, integer: true, boolean: true, null: true },
        rules: [{ rules: [] }, groups.number, groups.string, groups.array, groups.object],
        post: { rules: [] },
        all: {},
        keywords: {}
      };
    }
    exports.getRules = getRules;
  }
});

// node_modules/ajv/dist/compile/validate/applicability.js
var require_applicability = __commonJS({
  "node_modules/ajv/dist/compile/validate/applicability.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.shouldUseRule = exports.shouldUseGroup = exports.schemaHasRulesForType = void 0;
    function schemaHasRulesForType({ schema: schema2, self }, type) {
      const group = self.RULES.types[type];
      return group && group !== true && shouldUseGroup(schema2, group);
    }
    exports.schemaHasRulesForType = schemaHasRulesForType;
    function shouldUseGroup(schema2, group) {
      return group.rules.some((rule) => shouldUseRule(schema2, rule));
    }
    exports.shouldUseGroup = shouldUseGroup;
    function shouldUseRule(schema2, rule) {
      var _a;
      return schema2[rule.keyword] !== void 0 || ((_a = rule.definition.implements) === null || _a === void 0 ? void 0 : _a.some((kwd) => schema2[kwd] !== void 0));
    }
    exports.shouldUseRule = shouldUseRule;
  }
});

// node_modules/ajv/dist/compile/validate/dataType.js
var require_dataType = __commonJS({
  "node_modules/ajv/dist/compile/validate/dataType.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.reportTypeError = exports.checkDataTypes = exports.checkDataType = exports.coerceAndCheckDataType = exports.getJSONTypes = exports.getSchemaTypes = exports.DataType = void 0;
    var rules_1 = require_rules();
    var applicability_1 = require_applicability();
    var errors_1 = require_errors();
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var DataType;
    (function(DataType2) {
      DataType2[DataType2["Correct"] = 0] = "Correct";
      DataType2[DataType2["Wrong"] = 1] = "Wrong";
    })(DataType || (exports.DataType = DataType = {}));
    function getSchemaTypes(schema2) {
      const types = getJSONTypes(schema2.type);
      const hasNull = types.includes("null");
      if (hasNull) {
        if (schema2.nullable === false)
          throw new Error("type: null contradicts nullable: false");
      } else {
        if (!types.length && schema2.nullable !== void 0) {
          throw new Error('"nullable" cannot be used without "type"');
        }
        if (schema2.nullable === true)
          types.push("null");
      }
      return types;
    }
    exports.getSchemaTypes = getSchemaTypes;
    function getJSONTypes(ts) {
      const types = Array.isArray(ts) ? ts : ts ? [ts] : [];
      if (types.every(rules_1.isJSONType))
        return types;
      throw new Error("type must be JSONType or JSONType[]: " + types.join(","));
    }
    exports.getJSONTypes = getJSONTypes;
    function coerceAndCheckDataType(it, types) {
      const { gen, data, opts } = it;
      const coerceTo = coerceToTypes(types, opts.coerceTypes);
      const checkTypes = types.length > 0 && !(coerceTo.length === 0 && types.length === 1 && (0, applicability_1.schemaHasRulesForType)(it, types[0]));
      if (checkTypes) {
        const wrongType = checkDataTypes(types, data, opts.strictNumbers, DataType.Wrong);
        gen.if(wrongType, () => {
          if (coerceTo.length)
            coerceData(it, types, coerceTo);
          else
            reportTypeError(it);
        });
      }
      return checkTypes;
    }
    exports.coerceAndCheckDataType = coerceAndCheckDataType;
    var COERCIBLE = /* @__PURE__ */ new Set(["string", "number", "integer", "boolean", "null"]);
    function coerceToTypes(types, coerceTypes) {
      return coerceTypes ? types.filter((t) => COERCIBLE.has(t) || coerceTypes === "array" && t === "array") : [];
    }
    function coerceData(it, types, coerceTo) {
      const { gen, data, opts } = it;
      const dataType = gen.let("dataType", (0, codegen_1._)`typeof ${data}`);
      const coerced = gen.let("coerced", (0, codegen_1._)`undefined`);
      if (opts.coerceTypes === "array") {
        gen.if((0, codegen_1._)`${dataType} == 'object' && Array.isArray(${data}) && ${data}.length == 1`, () => gen.assign(data, (0, codegen_1._)`${data}[0]`).assign(dataType, (0, codegen_1._)`typeof ${data}`).if(checkDataTypes(types, data, opts.strictNumbers), () => gen.assign(coerced, data)));
      }
      gen.if((0, codegen_1._)`${coerced} !== undefined`);
      for (const t of coerceTo) {
        if (COERCIBLE.has(t) || t === "array" && opts.coerceTypes === "array") {
          coerceSpecificType(t);
        }
      }
      gen.else();
      reportTypeError(it);
      gen.endIf();
      gen.if((0, codegen_1._)`${coerced} !== undefined`, () => {
        gen.assign(data, coerced);
        assignParentData(it, coerced);
      });
      function coerceSpecificType(t) {
        switch (t) {
          case "string":
            gen.elseIf((0, codegen_1._)`${dataType} == "number" || ${dataType} == "boolean"`).assign(coerced, (0, codegen_1._)`"" + ${data}`).elseIf((0, codegen_1._)`${data} === null`).assign(coerced, (0, codegen_1._)`""`);
            return;
          case "number":
            gen.elseIf((0, codegen_1._)`${dataType} == "boolean" || ${data} === null
              || (${dataType} == "string" && ${data} && ${data} == +${data})`).assign(coerced, (0, codegen_1._)`+${data}`);
            return;
          case "integer":
            gen.elseIf((0, codegen_1._)`${dataType} === "boolean" || ${data} === null
              || (${dataType} === "string" && ${data} && ${data} == +${data} && !(${data} % 1))`).assign(coerced, (0, codegen_1._)`+${data}`);
            return;
          case "boolean":
            gen.elseIf((0, codegen_1._)`${data} === "false" || ${data} === 0 || ${data} === null`).assign(coerced, false).elseIf((0, codegen_1._)`${data} === "true" || ${data} === 1`).assign(coerced, true);
            return;
          case "null":
            gen.elseIf((0, codegen_1._)`${data} === "" || ${data} === 0 || ${data} === false`);
            gen.assign(coerced, null);
            return;
          case "array":
            gen.elseIf((0, codegen_1._)`${dataType} === "string" || ${dataType} === "number"
              || ${dataType} === "boolean" || ${data} === null`).assign(coerced, (0, codegen_1._)`[${data}]`);
        }
      }
    }
    function assignParentData({ gen, parentData, parentDataProperty }, expr) {
      gen.if((0, codegen_1._)`${parentData} !== undefined`, () => gen.assign((0, codegen_1._)`${parentData}[${parentDataProperty}]`, expr));
    }
    function checkDataType(dataType, data, strictNums, correct = DataType.Correct) {
      const EQ = correct === DataType.Correct ? codegen_1.operators.EQ : codegen_1.operators.NEQ;
      let cond;
      switch (dataType) {
        case "null":
          return (0, codegen_1._)`${data} ${EQ} null`;
        case "array":
          cond = (0, codegen_1._)`Array.isArray(${data})`;
          break;
        case "object":
          cond = (0, codegen_1._)`${data} && typeof ${data} == "object" && !Array.isArray(${data})`;
          break;
        case "integer":
          cond = numCond((0, codegen_1._)`!(${data} % 1) && !isNaN(${data})`);
          break;
        case "number":
          cond = numCond();
          break;
        default:
          return (0, codegen_1._)`typeof ${data} ${EQ} ${dataType}`;
      }
      return correct === DataType.Correct ? cond : (0, codegen_1.not)(cond);
      function numCond(_cond = codegen_1.nil) {
        return (0, codegen_1.and)((0, codegen_1._)`typeof ${data} == "number"`, _cond, strictNums ? (0, codegen_1._)`isFinite(${data})` : codegen_1.nil);
      }
    }
    exports.checkDataType = checkDataType;
    function checkDataTypes(dataTypes, data, strictNums, correct) {
      if (dataTypes.length === 1) {
        return checkDataType(dataTypes[0], data, strictNums, correct);
      }
      let cond;
      const types = (0, util_1.toHash)(dataTypes);
      if (types.array && types.object) {
        const notObj = (0, codegen_1._)`typeof ${data} != "object"`;
        cond = types.null ? notObj : (0, codegen_1._)`!${data} || ${notObj}`;
        delete types.null;
        delete types.array;
        delete types.object;
      } else {
        cond = codegen_1.nil;
      }
      if (types.number)
        delete types.integer;
      for (const t in types)
        cond = (0, codegen_1.and)(cond, checkDataType(t, data, strictNums, correct));
      return cond;
    }
    exports.checkDataTypes = checkDataTypes;
    var typeError = {
      message: ({ schema: schema2 }) => `must be ${schema2}`,
      params: ({ schema: schema2, schemaValue }) => typeof schema2 == "string" ? (0, codegen_1._)`{type: ${schema2}}` : (0, codegen_1._)`{type: ${schemaValue}}`
    };
    function reportTypeError(it) {
      const cxt = getTypeErrorContext(it);
      (0, errors_1.reportError)(cxt, typeError);
    }
    exports.reportTypeError = reportTypeError;
    function getTypeErrorContext(it) {
      const { gen, data, schema: schema2 } = it;
      const schemaCode = (0, util_1.schemaRefOrVal)(it, schema2, "type");
      return {
        gen,
        keyword: "type",
        data,
        schema: schema2.type,
        schemaCode,
        schemaValue: schemaCode,
        parentSchema: schema2,
        params: {},
        it
      };
    }
  }
});

// node_modules/ajv/dist/compile/validate/defaults.js
var require_defaults = __commonJS({
  "node_modules/ajv/dist/compile/validate/defaults.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.assignDefaults = void 0;
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    function assignDefaults(it, ty) {
      const { properties, items } = it.schema;
      if (ty === "object" && properties) {
        for (const key in properties) {
          assignDefault(it, key, properties[key].default);
        }
      } else if (ty === "array" && Array.isArray(items)) {
        items.forEach((sch, i) => assignDefault(it, i, sch.default));
      }
    }
    exports.assignDefaults = assignDefaults;
    function assignDefault(it, prop, defaultValue) {
      const { gen, compositeRule, data, opts } = it;
      if (defaultValue === void 0)
        return;
      const childData = (0, codegen_1._)`${data}${(0, codegen_1.getProperty)(prop)}`;
      if (compositeRule) {
        (0, util_1.checkStrictMode)(it, `default is ignored for: ${childData}`);
        return;
      }
      let condition = (0, codegen_1._)`${childData} === undefined`;
      if (opts.useDefaults === "empty") {
        condition = (0, codegen_1._)`${condition} || ${childData} === null || ${childData} === ""`;
      }
      gen.if(condition, (0, codegen_1._)`${childData} = ${(0, codegen_1.stringify)(defaultValue)}`);
    }
  }
});

// node_modules/ajv/dist/vocabularies/code.js
var require_code2 = __commonJS({
  "node_modules/ajv/dist/vocabularies/code.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.validateUnion = exports.validateArray = exports.usePattern = exports.callValidateCode = exports.schemaProperties = exports.allSchemaProperties = exports.noPropertyInData = exports.propertyInData = exports.isOwnProperty = exports.hasPropFunc = exports.reportMissingProp = exports.checkMissingProp = exports.checkReportMissingProp = void 0;
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var names_1 = require_names();
    var util_2 = require_util();
    function checkReportMissingProp(cxt, prop) {
      const { gen, data, it } = cxt;
      gen.if(noPropertyInData(gen, data, prop, it.opts.ownProperties), () => {
        cxt.setParams({ missingProperty: (0, codegen_1._)`${prop}` }, true);
        cxt.error();
      });
    }
    exports.checkReportMissingProp = checkReportMissingProp;
    function checkMissingProp({ gen, data, it: { opts } }, properties, missing) {
      return (0, codegen_1.or)(...properties.map((prop) => (0, codegen_1.and)(noPropertyInData(gen, data, prop, opts.ownProperties), (0, codegen_1._)`${missing} = ${prop}`)));
    }
    exports.checkMissingProp = checkMissingProp;
    function reportMissingProp(cxt, missing) {
      cxt.setParams({ missingProperty: missing }, true);
      cxt.error();
    }
    exports.reportMissingProp = reportMissingProp;
    function hasPropFunc(gen) {
      return gen.scopeValue("func", {
        // eslint-disable-next-line @typescript-eslint/unbound-method
        ref: Object.prototype.hasOwnProperty,
        code: (0, codegen_1._)`Object.prototype.hasOwnProperty`
      });
    }
    exports.hasPropFunc = hasPropFunc;
    function isOwnProperty(gen, data, property) {
      return (0, codegen_1._)`${hasPropFunc(gen)}.call(${data}, ${property})`;
    }
    exports.isOwnProperty = isOwnProperty;
    function propertyInData(gen, data, property, ownProperties) {
      const cond = (0, codegen_1._)`${data}${(0, codegen_1.getProperty)(property)} !== undefined`;
      return ownProperties ? (0, codegen_1._)`${cond} && ${isOwnProperty(gen, data, property)}` : cond;
    }
    exports.propertyInData = propertyInData;
    function noPropertyInData(gen, data, property, ownProperties) {
      const cond = (0, codegen_1._)`${data}${(0, codegen_1.getProperty)(property)} === undefined`;
      return ownProperties ? (0, codegen_1.or)(cond, (0, codegen_1.not)(isOwnProperty(gen, data, property))) : cond;
    }
    exports.noPropertyInData = noPropertyInData;
    function allSchemaProperties(schemaMap) {
      return schemaMap ? Object.keys(schemaMap).filter((p) => p !== "__proto__") : [];
    }
    exports.allSchemaProperties = allSchemaProperties;
    function schemaProperties(it, schemaMap) {
      return allSchemaProperties(schemaMap).filter((p) => !(0, util_1.alwaysValidSchema)(it, schemaMap[p]));
    }
    exports.schemaProperties = schemaProperties;
    function callValidateCode({ schemaCode, data, it: { gen, topSchemaRef, schemaPath, errorPath }, it }, func, context, passSchema) {
      const dataAndSchema = passSchema ? (0, codegen_1._)`${schemaCode}, ${data}, ${topSchemaRef}${schemaPath}` : data;
      const valCxt = [
        [names_1.default.instancePath, (0, codegen_1.strConcat)(names_1.default.instancePath, errorPath)],
        [names_1.default.parentData, it.parentData],
        [names_1.default.parentDataProperty, it.parentDataProperty],
        [names_1.default.rootData, names_1.default.rootData]
      ];
      if (it.opts.dynamicRef)
        valCxt.push([names_1.default.dynamicAnchors, names_1.default.dynamicAnchors]);
      const args = (0, codegen_1._)`${dataAndSchema}, ${gen.object(...valCxt)}`;
      return context !== codegen_1.nil ? (0, codegen_1._)`${func}.call(${context}, ${args})` : (0, codegen_1._)`${func}(${args})`;
    }
    exports.callValidateCode = callValidateCode;
    var newRegExp = (0, codegen_1._)`new RegExp`;
    function usePattern({ gen, it: { opts } }, pattern) {
      const u = opts.unicodeRegExp ? "u" : "";
      const { regExp } = opts.code;
      const rx = regExp(pattern, u);
      return gen.scopeValue("pattern", {
        key: rx.toString(),
        ref: rx,
        code: (0, codegen_1._)`${regExp.code === "new RegExp" ? newRegExp : (0, util_2.useFunc)(gen, regExp)}(${pattern}, ${u})`
      });
    }
    exports.usePattern = usePattern;
    function validateArray(cxt) {
      const { gen, data, keyword, it } = cxt;
      const valid = gen.name("valid");
      if (it.allErrors) {
        const validArr = gen.let("valid", true);
        validateItems(() => gen.assign(validArr, false));
        return validArr;
      }
      gen.var(valid, true);
      validateItems(() => gen.break());
      return valid;
      function validateItems(notValid) {
        const len = gen.const("len", (0, codegen_1._)`${data}.length`);
        gen.forRange("i", 0, len, (i) => {
          cxt.subschema({
            keyword,
            dataProp: i,
            dataPropType: util_1.Type.Num
          }, valid);
          gen.if((0, codegen_1.not)(valid), notValid);
        });
      }
    }
    exports.validateArray = validateArray;
    function validateUnion(cxt) {
      const { gen, schema: schema2, keyword, it } = cxt;
      if (!Array.isArray(schema2))
        throw new Error("ajv implementation error");
      const alwaysValid = schema2.some((sch) => (0, util_1.alwaysValidSchema)(it, sch));
      if (alwaysValid && !it.opts.unevaluated)
        return;
      const valid = gen.let("valid", false);
      const schValid = gen.name("_valid");
      gen.block(() => schema2.forEach((_sch, i) => {
        const schCxt = cxt.subschema({
          keyword,
          schemaProp: i,
          compositeRule: true
        }, schValid);
        gen.assign(valid, (0, codegen_1._)`${valid} || ${schValid}`);
        const merged = cxt.mergeValidEvaluated(schCxt, schValid);
        if (!merged)
          gen.if((0, codegen_1.not)(valid));
      }));
      cxt.result(valid, () => cxt.reset(), () => cxt.error(true));
    }
    exports.validateUnion = validateUnion;
  }
});

// node_modules/ajv/dist/compile/validate/keyword.js
var require_keyword = __commonJS({
  "node_modules/ajv/dist/compile/validate/keyword.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.validateKeywordUsage = exports.validSchemaType = exports.funcKeywordCode = exports.macroKeywordCode = void 0;
    var codegen_1 = require_codegen();
    var names_1 = require_names();
    var code_1 = require_code2();
    var errors_1 = require_errors();
    function macroKeywordCode(cxt, def) {
      const { gen, keyword, schema: schema2, parentSchema, it } = cxt;
      const macroSchema = def.macro.call(it.self, schema2, parentSchema, it);
      const schemaRef = useKeyword(gen, keyword, macroSchema);
      if (it.opts.validateSchema !== false)
        it.self.validateSchema(macroSchema, true);
      const valid = gen.name("valid");
      cxt.subschema({
        schema: macroSchema,
        schemaPath: codegen_1.nil,
        errSchemaPath: `${it.errSchemaPath}/${keyword}`,
        topSchemaRef: schemaRef,
        compositeRule: true
      }, valid);
      cxt.pass(valid, () => cxt.error(true));
    }
    exports.macroKeywordCode = macroKeywordCode;
    function funcKeywordCode(cxt, def) {
      var _a;
      const { gen, keyword, schema: schema2, parentSchema, $data, it } = cxt;
      checkAsyncKeyword(it, def);
      const validate = !$data && def.compile ? def.compile.call(it.self, schema2, parentSchema, it) : def.validate;
      const validateRef = useKeyword(gen, keyword, validate);
      const valid = gen.let("valid");
      cxt.block$data(valid, validateKeyword);
      cxt.ok((_a = def.valid) !== null && _a !== void 0 ? _a : valid);
      function validateKeyword() {
        if (def.errors === false) {
          assignValid();
          if (def.modifying)
            modifyData(cxt);
          reportErrs(() => cxt.error());
        } else {
          const ruleErrs = def.async ? validateAsync() : validateSync();
          if (def.modifying)
            modifyData(cxt);
          reportErrs(() => addErrs(cxt, ruleErrs));
        }
      }
      function validateAsync() {
        const ruleErrs = gen.let("ruleErrs", null);
        gen.try(() => assignValid((0, codegen_1._)`await `), (e) => gen.assign(valid, false).if((0, codegen_1._)`${e} instanceof ${it.ValidationError}`, () => gen.assign(ruleErrs, (0, codegen_1._)`${e}.errors`), () => gen.throw(e)));
        return ruleErrs;
      }
      function validateSync() {
        const validateErrs = (0, codegen_1._)`${validateRef}.errors`;
        gen.assign(validateErrs, null);
        assignValid(codegen_1.nil);
        return validateErrs;
      }
      function assignValid(_await = def.async ? (0, codegen_1._)`await ` : codegen_1.nil) {
        const passCxt = it.opts.passContext ? names_1.default.this : names_1.default.self;
        const passSchema = !("compile" in def && !$data || def.schema === false);
        gen.assign(valid, (0, codegen_1._)`${_await}${(0, code_1.callValidateCode)(cxt, validateRef, passCxt, passSchema)}`, def.modifying);
      }
      function reportErrs(errors2) {
        var _a2;
        gen.if((0, codegen_1.not)((_a2 = def.valid) !== null && _a2 !== void 0 ? _a2 : valid), errors2);
      }
    }
    exports.funcKeywordCode = funcKeywordCode;
    function modifyData(cxt) {
      const { gen, data, it } = cxt;
      gen.if(it.parentData, () => gen.assign(data, (0, codegen_1._)`${it.parentData}[${it.parentDataProperty}]`));
    }
    function addErrs(cxt, errs) {
      const { gen } = cxt;
      gen.if((0, codegen_1._)`Array.isArray(${errs})`, () => {
        gen.assign(names_1.default.vErrors, (0, codegen_1._)`${names_1.default.vErrors} === null ? ${errs} : ${names_1.default.vErrors}.concat(${errs})`).assign(names_1.default.errors, (0, codegen_1._)`${names_1.default.vErrors}.length`);
        (0, errors_1.extendErrors)(cxt);
      }, () => cxt.error());
    }
    function checkAsyncKeyword({ schemaEnv }, def) {
      if (def.async && !schemaEnv.$async)
        throw new Error("async keyword in sync schema");
    }
    function useKeyword(gen, keyword, result3) {
      if (result3 === void 0)
        throw new Error(`keyword "${keyword}" failed to compile`);
      return gen.scopeValue("keyword", typeof result3 == "function" ? { ref: result3 } : { ref: result3, code: (0, codegen_1.stringify)(result3) });
    }
    function validSchemaType(schema2, schemaType, allowUndefined = false) {
      return !schemaType.length || schemaType.some((st) => st === "array" ? Array.isArray(schema2) : st === "object" ? schema2 && typeof schema2 == "object" && !Array.isArray(schema2) : typeof schema2 == st || allowUndefined && typeof schema2 == "undefined");
    }
    exports.validSchemaType = validSchemaType;
    function validateKeywordUsage({ schema: schema2, opts, self, errSchemaPath }, def, keyword) {
      if (Array.isArray(def.keyword) ? !def.keyword.includes(keyword) : def.keyword !== keyword) {
        throw new Error("ajv implementation error");
      }
      const deps = def.dependencies;
      if (deps === null || deps === void 0 ? void 0 : deps.some((kwd) => !Object.prototype.hasOwnProperty.call(schema2, kwd))) {
        throw new Error(`parent schema must have dependencies of ${keyword}: ${deps.join(",")}`);
      }
      if (def.validateSchema) {
        const valid = def.validateSchema(schema2[keyword]);
        if (!valid) {
          const msg = `keyword "${keyword}" value is invalid at path "${errSchemaPath}": ` + self.errorsText(def.validateSchema.errors);
          if (opts.validateSchema === "log")
            self.logger.error(msg);
          else
            throw new Error(msg);
        }
      }
    }
    exports.validateKeywordUsage = validateKeywordUsage;
  }
});

// node_modules/ajv/dist/compile/validate/subschema.js
var require_subschema = __commonJS({
  "node_modules/ajv/dist/compile/validate/subschema.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.extendSubschemaMode = exports.extendSubschemaData = exports.getSubschema = void 0;
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    function getSubschema(it, { keyword, schemaProp, schema: schema2, schemaPath, errSchemaPath, topSchemaRef }) {
      if (keyword !== void 0 && schema2 !== void 0) {
        throw new Error('both "keyword" and "schema" passed, only one allowed');
      }
      if (keyword !== void 0) {
        const sch = it.schema[keyword];
        return schemaProp === void 0 ? {
          schema: sch,
          schemaPath: (0, codegen_1._)`${it.schemaPath}${(0, codegen_1.getProperty)(keyword)}`,
          errSchemaPath: `${it.errSchemaPath}/${keyword}`
        } : {
          schema: sch[schemaProp],
          schemaPath: (0, codegen_1._)`${it.schemaPath}${(0, codegen_1.getProperty)(keyword)}${(0, codegen_1.getProperty)(schemaProp)}`,
          errSchemaPath: `${it.errSchemaPath}/${keyword}/${(0, util_1.escapeFragment)(schemaProp)}`
        };
      }
      if (schema2 !== void 0) {
        if (schemaPath === void 0 || errSchemaPath === void 0 || topSchemaRef === void 0) {
          throw new Error('"schemaPath", "errSchemaPath" and "topSchemaRef" are required with "schema"');
        }
        return {
          schema: schema2,
          schemaPath,
          topSchemaRef,
          errSchemaPath
        };
      }
      throw new Error('either "keyword" or "schema" must be passed');
    }
    exports.getSubschema = getSubschema;
    function extendSubschemaData(subschema, it, { dataProp, dataPropType: dpType, data, dataTypes, propertyName }) {
      if (data !== void 0 && dataProp !== void 0) {
        throw new Error('both "data" and "dataProp" passed, only one allowed');
      }
      const { gen } = it;
      if (dataProp !== void 0) {
        const { errorPath, dataPathArr, opts } = it;
        const nextData = gen.let("data", (0, codegen_1._)`${it.data}${(0, codegen_1.getProperty)(dataProp)}`, true);
        dataContextProps(nextData);
        subschema.errorPath = (0, codegen_1.str)`${errorPath}${(0, util_1.getErrorPath)(dataProp, dpType, opts.jsPropertySyntax)}`;
        subschema.parentDataProperty = (0, codegen_1._)`${dataProp}`;
        subschema.dataPathArr = [...dataPathArr, subschema.parentDataProperty];
      }
      if (data !== void 0) {
        const nextData = data instanceof codegen_1.Name ? data : gen.let("data", data, true);
        dataContextProps(nextData);
        if (propertyName !== void 0)
          subschema.propertyName = propertyName;
      }
      if (dataTypes)
        subschema.dataTypes = dataTypes;
      function dataContextProps(_nextData) {
        subschema.data = _nextData;
        subschema.dataLevel = it.dataLevel + 1;
        subschema.dataTypes = [];
        it.definedProperties = /* @__PURE__ */ new Set();
        subschema.parentData = it.data;
        subschema.dataNames = [...it.dataNames, _nextData];
      }
    }
    exports.extendSubschemaData = extendSubschemaData;
    function extendSubschemaMode(subschema, { jtdDiscriminator, jtdMetadata, compositeRule, createErrors, allErrors }) {
      if (compositeRule !== void 0)
        subschema.compositeRule = compositeRule;
      if (createErrors !== void 0)
        subschema.createErrors = createErrors;
      if (allErrors !== void 0)
        subschema.allErrors = allErrors;
      subschema.jtdDiscriminator = jtdDiscriminator;
      subschema.jtdMetadata = jtdMetadata;
    }
    exports.extendSubschemaMode = extendSubschemaMode;
  }
});

// node_modules/fast-deep-equal/index.js
var require_fast_deep_equal = __commonJS({
  "node_modules/fast-deep-equal/index.js"(exports, module) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    module.exports = function equal(a, b) {
      if (a === b) return true;
      if (a && b && typeof a == "object" && typeof b == "object") {
        if (a.constructor !== b.constructor) return false;
        var length, i, keys;
        if (Array.isArray(a)) {
          length = a.length;
          if (length != b.length) return false;
          for (i = length; i-- !== 0; )
            if (!equal(a[i], b[i])) return false;
          return true;
        }
        if (a.constructor === RegExp) return a.source === b.source && a.flags === b.flags;
        if (a.valueOf !== Object.prototype.valueOf) return a.valueOf() === b.valueOf();
        if (a.toString !== Object.prototype.toString) return a.toString() === b.toString();
        keys = Object.keys(a);
        length = keys.length;
        if (length !== Object.keys(b).length) return false;
        for (i = length; i-- !== 0; )
          if (!Object.prototype.hasOwnProperty.call(b, keys[i])) return false;
        for (i = length; i-- !== 0; ) {
          var key = keys[i];
          if (!equal(a[key], b[key])) return false;
        }
        return true;
      }
      return a !== a && b !== b;
    };
  }
});

// node_modules/json-schema-traverse/index.js
var require_json_schema_traverse = __commonJS({
  "node_modules/json-schema-traverse/index.js"(exports, module) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var traverse = module.exports = function(schema2, opts, cb) {
      if (typeof opts == "function") {
        cb = opts;
        opts = {};
      }
      cb = opts.cb || cb;
      var pre = typeof cb == "function" ? cb : cb.pre || function() {
      };
      var post = cb.post || function() {
      };
      _traverse(opts, pre, post, schema2, "", schema2);
    };
    traverse.keywords = {
      additionalItems: true,
      items: true,
      contains: true,
      additionalProperties: true,
      propertyNames: true,
      not: true,
      if: true,
      then: true,
      else: true
    };
    traverse.arrayKeywords = {
      items: true,
      allOf: true,
      anyOf: true,
      oneOf: true
    };
    traverse.propsKeywords = {
      $defs: true,
      definitions: true,
      properties: true,
      patternProperties: true,
      dependencies: true
    };
    traverse.skipKeywords = {
      default: true,
      enum: true,
      const: true,
      required: true,
      maximum: true,
      minimum: true,
      exclusiveMaximum: true,
      exclusiveMinimum: true,
      multipleOf: true,
      maxLength: true,
      minLength: true,
      pattern: true,
      format: true,
      maxItems: true,
      minItems: true,
      uniqueItems: true,
      maxProperties: true,
      minProperties: true
    };
    function _traverse(opts, pre, post, schema2, jsonPtr, rootSchema, parentJsonPtr, parentKeyword, parentSchema, keyIndex) {
      if (schema2 && typeof schema2 == "object" && !Array.isArray(schema2)) {
        pre(schema2, jsonPtr, rootSchema, parentJsonPtr, parentKeyword, parentSchema, keyIndex);
        for (var key in schema2) {
          var sch = schema2[key];
          if (Array.isArray(sch)) {
            if (key in traverse.arrayKeywords) {
              for (var i = 0; i < sch.length; i++)
                _traverse(opts, pre, post, sch[i], jsonPtr + "/" + key + "/" + i, rootSchema, jsonPtr, key, schema2, i);
            }
          } else if (key in traverse.propsKeywords) {
            if (sch && typeof sch == "object") {
              for (var prop in sch)
                _traverse(opts, pre, post, sch[prop], jsonPtr + "/" + key + "/" + escapeJsonPtr(prop), rootSchema, jsonPtr, key, schema2, prop);
            }
          } else if (key in traverse.keywords || opts.allKeys && !(key in traverse.skipKeywords)) {
            _traverse(opts, pre, post, sch, jsonPtr + "/" + key, rootSchema, jsonPtr, key, schema2);
          }
        }
        post(schema2, jsonPtr, rootSchema, parentJsonPtr, parentKeyword, parentSchema, keyIndex);
      }
    }
    function escapeJsonPtr(str) {
      return str.replace(/~/g, "~0").replace(/\//g, "~1");
    }
  }
});

// node_modules/ajv/dist/compile/resolve.js
var require_resolve = __commonJS({
  "node_modules/ajv/dist/compile/resolve.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.getSchemaRefs = exports.resolveUrl = exports.normalizeId = exports._getFullPath = exports.getFullPath = exports.inlineRef = void 0;
    var util_1 = require_util();
    var equal = require_fast_deep_equal();
    var traverse = require_json_schema_traverse();
    var SIMPLE_INLINED = /* @__PURE__ */ new Set([
      "type",
      "format",
      "pattern",
      "maxLength",
      "minLength",
      "maxProperties",
      "minProperties",
      "maxItems",
      "minItems",
      "maximum",
      "minimum",
      "uniqueItems",
      "multipleOf",
      "required",
      "enum",
      "const"
    ]);
    function inlineRef(schema2, limit = true) {
      if (typeof schema2 == "boolean")
        return true;
      if (limit === true)
        return !hasRef(schema2);
      if (!limit)
        return false;
      return countKeys(schema2) <= limit;
    }
    exports.inlineRef = inlineRef;
    var REF_KEYWORDS = /* @__PURE__ */ new Set([
      "$ref",
      "$recursiveRef",
      "$recursiveAnchor",
      "$dynamicRef",
      "$dynamicAnchor"
    ]);
    function hasRef(schema2) {
      for (const key in schema2) {
        if (REF_KEYWORDS.has(key))
          return true;
        const sch = schema2[key];
        if (Array.isArray(sch) && sch.some(hasRef))
          return true;
        if (typeof sch == "object" && hasRef(sch))
          return true;
      }
      return false;
    }
    function countKeys(schema2) {
      let count = 0;
      for (const key in schema2) {
        if (key === "$ref")
          return Infinity;
        count++;
        if (SIMPLE_INLINED.has(key))
          continue;
        if (typeof schema2[key] == "object") {
          (0, util_1.eachItem)(schema2[key], (sch) => count += countKeys(sch));
        }
        if (count === Infinity)
          return Infinity;
      }
      return count;
    }
    function getFullPath(resolver, id = "", normalize) {
      if (normalize !== false)
        id = normalizeId(id);
      const p = resolver.parse(id);
      return _getFullPath(resolver, p);
    }
    exports.getFullPath = getFullPath;
    function _getFullPath(resolver, p) {
      const serialized = resolver.serialize(p);
      return serialized.split("#")[0] + "#";
    }
    exports._getFullPath = _getFullPath;
    var TRAILING_SLASH_HASH = /#\/?$/;
    function normalizeId(id) {
      return id ? id.replace(TRAILING_SLASH_HASH, "") : "";
    }
    exports.normalizeId = normalizeId;
    function resolveUrl(resolver, baseId, id) {
      id = normalizeId(id);
      return resolver.resolve(baseId, id);
    }
    exports.resolveUrl = resolveUrl;
    var ANCHOR = /^[a-z_][-a-z0-9._]*$/i;
    function getSchemaRefs(schema2, baseId) {
      if (typeof schema2 == "boolean")
        return {};
      const { schemaId, uriResolver } = this.opts;
      const schId = normalizeId(schema2[schemaId] || baseId);
      const baseIds = { "": schId };
      const pathPrefix = getFullPath(uriResolver, schId, false);
      const localRefs = {};
      const schemaRefs = /* @__PURE__ */ new Set();
      traverse(schema2, { allKeys: true }, (sch, jsonPtr, _, parentJsonPtr) => {
        if (parentJsonPtr === void 0)
          return;
        const fullPath = pathPrefix + jsonPtr;
        let innerBaseId = baseIds[parentJsonPtr];
        if (typeof sch[schemaId] == "string")
          innerBaseId = addRef.call(this, sch[schemaId]);
        addAnchor.call(this, sch.$anchor);
        addAnchor.call(this, sch.$dynamicAnchor);
        baseIds[jsonPtr] = innerBaseId;
        function addRef(ref) {
          const _resolve = this.opts.uriResolver.resolve;
          ref = normalizeId(innerBaseId ? _resolve(innerBaseId, ref) : ref);
          if (schemaRefs.has(ref))
            throw ambiguos(ref);
          schemaRefs.add(ref);
          let schOrRef = this.refs[ref];
          if (typeof schOrRef == "string")
            schOrRef = this.refs[schOrRef];
          if (typeof schOrRef == "object") {
            checkAmbiguosRef(sch, schOrRef.schema, ref);
          } else if (ref !== normalizeId(fullPath)) {
            if (ref[0] === "#") {
              checkAmbiguosRef(sch, localRefs[ref], ref);
              localRefs[ref] = sch;
            } else {
              this.refs[ref] = fullPath;
            }
          }
          return ref;
        }
        function addAnchor(anchor) {
          if (typeof anchor == "string") {
            if (!ANCHOR.test(anchor))
              throw new Error(`invalid anchor "${anchor}"`);
            addRef.call(this, `#${anchor}`);
          }
        }
      });
      return localRefs;
      function checkAmbiguosRef(sch1, sch2, ref) {
        if (sch2 !== void 0 && !equal(sch1, sch2))
          throw ambiguos(ref);
      }
      function ambiguos(ref) {
        return new Error(`reference "${ref}" resolves to more than one schema`);
      }
    }
    exports.getSchemaRefs = getSchemaRefs;
  }
});

// node_modules/ajv/dist/compile/validate/index.js
var require_validate = __commonJS({
  "node_modules/ajv/dist/compile/validate/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.getData = exports.KeywordCxt = exports.validateFunctionCode = void 0;
    var boolSchema_1 = require_boolSchema();
    var dataType_1 = require_dataType();
    var applicability_1 = require_applicability();
    var dataType_2 = require_dataType();
    var defaults_1 = require_defaults();
    var keyword_1 = require_keyword();
    var subschema_1 = require_subschema();
    var codegen_1 = require_codegen();
    var names_1 = require_names();
    var resolve_1 = require_resolve();
    var util_1 = require_util();
    var errors_1 = require_errors();
    function validateFunctionCode(it) {
      if (isSchemaObj(it)) {
        checkKeywords(it);
        if (schemaCxtHasRules(it)) {
          topSchemaObjCode(it);
          return;
        }
      }
      validateFunction(it, () => (0, boolSchema_1.topBoolOrEmptySchema)(it));
    }
    exports.validateFunctionCode = validateFunctionCode;
    function validateFunction({ gen, validateName, schema: schema2, schemaEnv, opts }, body) {
      if (opts.code.es5) {
        gen.func(validateName, (0, codegen_1._)`${names_1.default.data}, ${names_1.default.valCxt}`, schemaEnv.$async, () => {
          gen.code((0, codegen_1._)`"use strict"; ${funcSourceUrl(schema2, opts)}`);
          destructureValCxtES5(gen, opts);
          gen.code(body);
        });
      } else {
        gen.func(validateName, (0, codegen_1._)`${names_1.default.data}, ${destructureValCxt(opts)}`, schemaEnv.$async, () => gen.code(funcSourceUrl(schema2, opts)).code(body));
      }
    }
    function destructureValCxt(opts) {
      return (0, codegen_1._)`{${names_1.default.instancePath}="", ${names_1.default.parentData}, ${names_1.default.parentDataProperty}, ${names_1.default.rootData}=${names_1.default.data}${opts.dynamicRef ? (0, codegen_1._)`, ${names_1.default.dynamicAnchors}={}` : codegen_1.nil}}={}`;
    }
    function destructureValCxtES5(gen, opts) {
      gen.if(names_1.default.valCxt, () => {
        gen.var(names_1.default.instancePath, (0, codegen_1._)`${names_1.default.valCxt}.${names_1.default.instancePath}`);
        gen.var(names_1.default.parentData, (0, codegen_1._)`${names_1.default.valCxt}.${names_1.default.parentData}`);
        gen.var(names_1.default.parentDataProperty, (0, codegen_1._)`${names_1.default.valCxt}.${names_1.default.parentDataProperty}`);
        gen.var(names_1.default.rootData, (0, codegen_1._)`${names_1.default.valCxt}.${names_1.default.rootData}`);
        if (opts.dynamicRef)
          gen.var(names_1.default.dynamicAnchors, (0, codegen_1._)`${names_1.default.valCxt}.${names_1.default.dynamicAnchors}`);
      }, () => {
        gen.var(names_1.default.instancePath, (0, codegen_1._)`""`);
        gen.var(names_1.default.parentData, (0, codegen_1._)`undefined`);
        gen.var(names_1.default.parentDataProperty, (0, codegen_1._)`undefined`);
        gen.var(names_1.default.rootData, names_1.default.data);
        if (opts.dynamicRef)
          gen.var(names_1.default.dynamicAnchors, (0, codegen_1._)`{}`);
      });
    }
    function topSchemaObjCode(it) {
      const { schema: schema2, opts, gen } = it;
      validateFunction(it, () => {
        if (opts.$comment && schema2.$comment)
          commentKeyword(it);
        checkNoDefault(it);
        gen.let(names_1.default.vErrors, null);
        gen.let(names_1.default.errors, 0);
        if (opts.unevaluated)
          resetEvaluated(it);
        typeAndKeywords(it);
        returnResults(it);
      });
      return;
    }
    function resetEvaluated(it) {
      const { gen, validateName } = it;
      it.evaluated = gen.const("evaluated", (0, codegen_1._)`${validateName}.evaluated`);
      gen.if((0, codegen_1._)`${it.evaluated}.dynamicProps`, () => gen.assign((0, codegen_1._)`${it.evaluated}.props`, (0, codegen_1._)`undefined`));
      gen.if((0, codegen_1._)`${it.evaluated}.dynamicItems`, () => gen.assign((0, codegen_1._)`${it.evaluated}.items`, (0, codegen_1._)`undefined`));
    }
    function funcSourceUrl(schema2, opts) {
      const schId = typeof schema2 == "object" && schema2[opts.schemaId];
      return schId && (opts.code.source || opts.code.process) ? (0, codegen_1._)`/*# sourceURL=${schId} */` : codegen_1.nil;
    }
    function subschemaCode(it, valid) {
      if (isSchemaObj(it)) {
        checkKeywords(it);
        if (schemaCxtHasRules(it)) {
          subSchemaObjCode(it, valid);
          return;
        }
      }
      (0, boolSchema_1.boolOrEmptySchema)(it, valid);
    }
    function schemaCxtHasRules({ schema: schema2, self }) {
      if (typeof schema2 == "boolean")
        return !schema2;
      for (const key in schema2)
        if (self.RULES.all[key])
          return true;
      return false;
    }
    function isSchemaObj(it) {
      return typeof it.schema != "boolean";
    }
    function subSchemaObjCode(it, valid) {
      const { schema: schema2, gen, opts } = it;
      if (opts.$comment && schema2.$comment)
        commentKeyword(it);
      updateContext(it);
      checkAsyncSchema(it);
      const errsCount = gen.const("_errs", names_1.default.errors);
      typeAndKeywords(it, errsCount);
      gen.var(valid, (0, codegen_1._)`${errsCount} === ${names_1.default.errors}`);
    }
    function checkKeywords(it) {
      (0, util_1.checkUnknownRules)(it);
      checkRefsAndKeywords(it);
    }
    function typeAndKeywords(it, errsCount) {
      if (it.opts.jtd)
        return schemaKeywords(it, [], false, errsCount);
      const types = (0, dataType_1.getSchemaTypes)(it.schema);
      const checkedTypes = (0, dataType_1.coerceAndCheckDataType)(it, types);
      schemaKeywords(it, types, !checkedTypes, errsCount);
    }
    function checkRefsAndKeywords(it) {
      const { schema: schema2, errSchemaPath, opts, self } = it;
      if (schema2.$ref && opts.ignoreKeywordsWithRef && (0, util_1.schemaHasRulesButRef)(schema2, self.RULES)) {
        self.logger.warn(`$ref: keywords ignored in schema at path "${errSchemaPath}"`);
      }
    }
    function checkNoDefault(it) {
      const { schema: schema2, opts } = it;
      if (schema2.default !== void 0 && opts.useDefaults && opts.strictSchema) {
        (0, util_1.checkStrictMode)(it, "default is ignored in the schema root");
      }
    }
    function updateContext(it) {
      const schId = it.schema[it.opts.schemaId];
      if (schId)
        it.baseId = (0, resolve_1.resolveUrl)(it.opts.uriResolver, it.baseId, schId);
    }
    function checkAsyncSchema(it) {
      if (it.schema.$async && !it.schemaEnv.$async)
        throw new Error("async schema in sync schema");
    }
    function commentKeyword({ gen, schemaEnv, schema: schema2, errSchemaPath, opts }) {
      const msg = schema2.$comment;
      if (opts.$comment === true) {
        gen.code((0, codegen_1._)`${names_1.default.self}.logger.log(${msg})`);
      } else if (typeof opts.$comment == "function") {
        const schemaPath = (0, codegen_1.str)`${errSchemaPath}/$comment`;
        const rootName = gen.scopeValue("root", { ref: schemaEnv.root });
        gen.code((0, codegen_1._)`${names_1.default.self}.opts.$comment(${msg}, ${schemaPath}, ${rootName}.schema)`);
      }
    }
    function returnResults(it) {
      const { gen, schemaEnv, validateName, ValidationError, opts } = it;
      if (schemaEnv.$async) {
        gen.if((0, codegen_1._)`${names_1.default.errors} === 0`, () => gen.return(names_1.default.data), () => gen.throw((0, codegen_1._)`new ${ValidationError}(${names_1.default.vErrors})`));
      } else {
        gen.assign((0, codegen_1._)`${validateName}.errors`, names_1.default.vErrors);
        if (opts.unevaluated)
          assignEvaluated(it);
        gen.return((0, codegen_1._)`${names_1.default.errors} === 0`);
      }
    }
    function assignEvaluated({ gen, evaluated, props, items }) {
      if (props instanceof codegen_1.Name)
        gen.assign((0, codegen_1._)`${evaluated}.props`, props);
      if (items instanceof codegen_1.Name)
        gen.assign((0, codegen_1._)`${evaluated}.items`, items);
    }
    function schemaKeywords(it, types, typeErrors, errsCount) {
      const { gen, schema: schema2, data, allErrors, opts, self } = it;
      const { RULES } = self;
      if (schema2.$ref && (opts.ignoreKeywordsWithRef || !(0, util_1.schemaHasRulesButRef)(schema2, RULES))) {
        gen.block(() => keywordCode(it, "$ref", RULES.all.$ref.definition));
        return;
      }
      if (!opts.jtd)
        checkStrictTypes(it, types);
      gen.block(() => {
        for (const group of RULES.rules)
          groupKeywords(group);
        groupKeywords(RULES.post);
      });
      function groupKeywords(group) {
        if (!(0, applicability_1.shouldUseGroup)(schema2, group))
          return;
        if (group.type) {
          gen.if((0, dataType_2.checkDataType)(group.type, data, opts.strictNumbers));
          iterateKeywords(it, group);
          if (types.length === 1 && types[0] === group.type && typeErrors) {
            gen.else();
            (0, dataType_2.reportTypeError)(it);
          }
          gen.endIf();
        } else {
          iterateKeywords(it, group);
        }
        if (!allErrors)
          gen.if((0, codegen_1._)`${names_1.default.errors} === ${errsCount || 0}`);
      }
    }
    function iterateKeywords(it, group) {
      const { gen, schema: schema2, opts: { useDefaults } } = it;
      if (useDefaults)
        (0, defaults_1.assignDefaults)(it, group.type);
      gen.block(() => {
        for (const rule of group.rules) {
          if ((0, applicability_1.shouldUseRule)(schema2, rule)) {
            keywordCode(it, rule.keyword, rule.definition, group.type);
          }
        }
      });
    }
    function checkStrictTypes(it, types) {
      if (it.schemaEnv.meta || !it.opts.strictTypes)
        return;
      checkContextTypes(it, types);
      if (!it.opts.allowUnionTypes)
        checkMultipleTypes(it, types);
      checkKeywordTypes(it, it.dataTypes);
    }
    function checkContextTypes(it, types) {
      if (!types.length)
        return;
      if (!it.dataTypes.length) {
        it.dataTypes = types;
        return;
      }
      types.forEach((t) => {
        if (!includesType(it.dataTypes, t)) {
          strictTypesError(it, `type "${t}" not allowed by context "${it.dataTypes.join(",")}"`);
        }
      });
      narrowSchemaTypes(it, types);
    }
    function checkMultipleTypes(it, ts) {
      if (ts.length > 1 && !(ts.length === 2 && ts.includes("null"))) {
        strictTypesError(it, "use allowUnionTypes to allow union type keyword");
      }
    }
    function checkKeywordTypes(it, ts) {
      const rules = it.self.RULES.all;
      for (const keyword in rules) {
        const rule = rules[keyword];
        if (typeof rule == "object" && (0, applicability_1.shouldUseRule)(it.schema, rule)) {
          const { type } = rule.definition;
          if (type.length && !type.some((t) => hasApplicableType(ts, t))) {
            strictTypesError(it, `missing type "${type.join(",")}" for keyword "${keyword}"`);
          }
        }
      }
    }
    function hasApplicableType(schTs, kwdT) {
      return schTs.includes(kwdT) || kwdT === "number" && schTs.includes("integer");
    }
    function includesType(ts, t) {
      return ts.includes(t) || t === "integer" && ts.includes("number");
    }
    function narrowSchemaTypes(it, withTypes) {
      const ts = [];
      for (const t of it.dataTypes) {
        if (includesType(withTypes, t))
          ts.push(t);
        else if (withTypes.includes("integer") && t === "number")
          ts.push("integer");
      }
      it.dataTypes = ts;
    }
    function strictTypesError(it, msg) {
      const schemaPath = it.schemaEnv.baseId + it.errSchemaPath;
      msg += ` at "${schemaPath}" (strictTypes)`;
      (0, util_1.checkStrictMode)(it, msg, it.opts.strictTypes);
    }
    var KeywordCxt = class {
      constructor(it, def, keyword) {
        (0, keyword_1.validateKeywordUsage)(it, def, keyword);
        this.gen = it.gen;
        this.allErrors = it.allErrors;
        this.keyword = keyword;
        this.data = it.data;
        this.schema = it.schema[keyword];
        this.$data = def.$data && it.opts.$data && this.schema && this.schema.$data;
        this.schemaValue = (0, util_1.schemaRefOrVal)(it, this.schema, keyword, this.$data);
        this.schemaType = def.schemaType;
        this.parentSchema = it.schema;
        this.params = {};
        this.it = it;
        this.def = def;
        if (this.$data) {
          this.schemaCode = it.gen.const("vSchema", getData(this.$data, it));
        } else {
          this.schemaCode = this.schemaValue;
          if (!(0, keyword_1.validSchemaType)(this.schema, def.schemaType, def.allowUndefined)) {
            throw new Error(`${keyword} value must be ${JSON.stringify(def.schemaType)}`);
          }
        }
        if ("code" in def ? def.trackErrors : def.errors !== false) {
          this.errsCount = it.gen.const("_errs", names_1.default.errors);
        }
      }
      result(condition, successAction, failAction) {
        this.failResult((0, codegen_1.not)(condition), successAction, failAction);
      }
      failResult(condition, successAction, failAction) {
        this.gen.if(condition);
        if (failAction)
          failAction();
        else
          this.error();
        if (successAction) {
          this.gen.else();
          successAction();
          if (this.allErrors)
            this.gen.endIf();
        } else {
          if (this.allErrors)
            this.gen.endIf();
          else
            this.gen.else();
        }
      }
      pass(condition, failAction) {
        this.failResult((0, codegen_1.not)(condition), void 0, failAction);
      }
      fail(condition) {
        if (condition === void 0) {
          this.error();
          if (!this.allErrors)
            this.gen.if(false);
          return;
        }
        this.gen.if(condition);
        this.error();
        if (this.allErrors)
          this.gen.endIf();
        else
          this.gen.else();
      }
      fail$data(condition) {
        if (!this.$data)
          return this.fail(condition);
        const { schemaCode } = this;
        this.fail((0, codegen_1._)`${schemaCode} !== undefined && (${(0, codegen_1.or)(this.invalid$data(), condition)})`);
      }
      error(append, errorParams, errorPaths) {
        if (errorParams) {
          this.setParams(errorParams);
          this._error(append, errorPaths);
          this.setParams({});
          return;
        }
        this._error(append, errorPaths);
      }
      _error(append, errorPaths) {
        ;
        (append ? errors_1.reportExtraError : errors_1.reportError)(this, this.def.error, errorPaths);
      }
      $dataError() {
        (0, errors_1.reportError)(this, this.def.$dataError || errors_1.keyword$DataError);
      }
      reset() {
        if (this.errsCount === void 0)
          throw new Error('add "trackErrors" to keyword definition');
        (0, errors_1.resetErrorsCount)(this.gen, this.errsCount);
      }
      ok(cond) {
        if (!this.allErrors)
          this.gen.if(cond);
      }
      setParams(obj, assign) {
        if (assign)
          Object.assign(this.params, obj);
        else
          this.params = obj;
      }
      block$data(valid, codeBlock, $dataValid = codegen_1.nil) {
        this.gen.block(() => {
          this.check$data(valid, $dataValid);
          codeBlock();
        });
      }
      check$data(valid = codegen_1.nil, $dataValid = codegen_1.nil) {
        if (!this.$data)
          return;
        const { gen, schemaCode, schemaType, def } = this;
        gen.if((0, codegen_1.or)((0, codegen_1._)`${schemaCode} === undefined`, $dataValid));
        if (valid !== codegen_1.nil)
          gen.assign(valid, true);
        if (schemaType.length || def.validateSchema) {
          gen.elseIf(this.invalid$data());
          this.$dataError();
          if (valid !== codegen_1.nil)
            gen.assign(valid, false);
        }
        gen.else();
      }
      invalid$data() {
        const { gen, schemaCode, schemaType, def, it } = this;
        return (0, codegen_1.or)(wrong$DataType(), invalid$DataSchema());
        function wrong$DataType() {
          if (schemaType.length) {
            if (!(schemaCode instanceof codegen_1.Name))
              throw new Error("ajv implementation error");
            const st = Array.isArray(schemaType) ? schemaType : [schemaType];
            return (0, codegen_1._)`${(0, dataType_2.checkDataTypes)(st, schemaCode, it.opts.strictNumbers, dataType_2.DataType.Wrong)}`;
          }
          return codegen_1.nil;
        }
        function invalid$DataSchema() {
          if (def.validateSchema) {
            const validateSchemaRef = gen.scopeValue("validate$data", { ref: def.validateSchema });
            return (0, codegen_1._)`!${validateSchemaRef}(${schemaCode})`;
          }
          return codegen_1.nil;
        }
      }
      subschema(appl, valid) {
        const subschema = (0, subschema_1.getSubschema)(this.it, appl);
        (0, subschema_1.extendSubschemaData)(subschema, this.it, appl);
        (0, subschema_1.extendSubschemaMode)(subschema, appl);
        const nextContext = { ...this.it, ...subschema, items: void 0, props: void 0 };
        subschemaCode(nextContext, valid);
        return nextContext;
      }
      mergeEvaluated(schemaCxt, toName) {
        const { it, gen } = this;
        if (!it.opts.unevaluated)
          return;
        if (it.props !== true && schemaCxt.props !== void 0) {
          it.props = util_1.mergeEvaluated.props(gen, schemaCxt.props, it.props, toName);
        }
        if (it.items !== true && schemaCxt.items !== void 0) {
          it.items = util_1.mergeEvaluated.items(gen, schemaCxt.items, it.items, toName);
        }
      }
      mergeValidEvaluated(schemaCxt, valid) {
        const { it, gen } = this;
        if (it.opts.unevaluated && (it.props !== true || it.items !== true)) {
          gen.if(valid, () => this.mergeEvaluated(schemaCxt, codegen_1.Name));
          return true;
        }
      }
    };
    exports.KeywordCxt = KeywordCxt;
    function keywordCode(it, keyword, def, ruleType) {
      const cxt = new KeywordCxt(it, def, keyword);
      if ("code" in def) {
        def.code(cxt, ruleType);
      } else if (cxt.$data && def.validate) {
        (0, keyword_1.funcKeywordCode)(cxt, def);
      } else if ("macro" in def) {
        (0, keyword_1.macroKeywordCode)(cxt, def);
      } else if (def.compile || def.validate) {
        (0, keyword_1.funcKeywordCode)(cxt, def);
      }
    }
    var JSON_POINTER = /^\/(?:[^~]|~0|~1)*$/;
    var RELATIVE_JSON_POINTER = /^([0-9]+)(#|\/(?:[^~]|~0|~1)*)?$/;
    function getData($data, { dataLevel, dataNames, dataPathArr }) {
      let jsonPointer;
      let data;
      if ($data === "")
        return names_1.default.rootData;
      if ($data[0] === "/") {
        if (!JSON_POINTER.test($data))
          throw new Error(`Invalid JSON-pointer: ${$data}`);
        jsonPointer = $data;
        data = names_1.default.rootData;
      } else {
        const matches = RELATIVE_JSON_POINTER.exec($data);
        if (!matches)
          throw new Error(`Invalid JSON-pointer: ${$data}`);
        const up = +matches[1];
        jsonPointer = matches[2];
        if (jsonPointer === "#") {
          if (up >= dataLevel)
            throw new Error(errorMsg("property/index", up));
          return dataPathArr[dataLevel - up];
        }
        if (up > dataLevel)
          throw new Error(errorMsg("data", up));
        data = dataNames[dataLevel - up];
        if (!jsonPointer)
          return data;
      }
      let expr = data;
      const segments = jsonPointer.split("/");
      for (const segment of segments) {
        if (segment) {
          data = (0, codegen_1._)`${data}${(0, codegen_1.getProperty)((0, util_1.unescapeJsonPointer)(segment))}`;
          expr = (0, codegen_1._)`${expr} && ${data}`;
        }
      }
      return expr;
      function errorMsg(pointerType, up) {
        return `Cannot access ${pointerType} ${up} levels up, current level is ${dataLevel}`;
      }
    }
    exports.getData = getData;
  }
});

// node_modules/ajv/dist/runtime/validation_error.js
var require_validation_error = __commonJS({
  "node_modules/ajv/dist/runtime/validation_error.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var ValidationError = class extends Error {
      constructor(errors2) {
        super("validation failed");
        this.errors = errors2;
        this.ajv = this.validation = true;
      }
    };
    exports.default = ValidationError;
  }
});

// node_modules/ajv/dist/compile/ref_error.js
var require_ref_error = __commonJS({
  "node_modules/ajv/dist/compile/ref_error.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var resolve_1 = require_resolve();
    var MissingRefError = class extends Error {
      constructor(resolver, baseId, ref, msg) {
        super(msg || `can't resolve reference ${ref} from id ${baseId}`);
        this.missingRef = (0, resolve_1.resolveUrl)(resolver, baseId, ref);
        this.missingSchema = (0, resolve_1.normalizeId)((0, resolve_1.getFullPath)(resolver, this.missingRef));
      }
    };
    exports.default = MissingRefError;
  }
});

// node_modules/ajv/dist/compile/index.js
var require_compile = __commonJS({
  "node_modules/ajv/dist/compile/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.resolveSchema = exports.getCompilingSchema = exports.resolveRef = exports.compileSchema = exports.SchemaEnv = void 0;
    var codegen_1 = require_codegen();
    var validation_error_1 = require_validation_error();
    var names_1 = require_names();
    var resolve_1 = require_resolve();
    var util_1 = require_util();
    var validate_1 = require_validate();
    var SchemaEnv = class {
      constructor(env) {
        var _a;
        this.refs = {};
        this.dynamicAnchors = {};
        let schema2;
        if (typeof env.schema == "object")
          schema2 = env.schema;
        this.schema = env.schema;
        this.schemaId = env.schemaId;
        this.root = env.root || this;
        this.baseId = (_a = env.baseId) !== null && _a !== void 0 ? _a : (0, resolve_1.normalizeId)(schema2 === null || schema2 === void 0 ? void 0 : schema2[env.schemaId || "$id"]);
        this.schemaPath = env.schemaPath;
        this.localRefs = env.localRefs;
        this.meta = env.meta;
        this.$async = schema2 === null || schema2 === void 0 ? void 0 : schema2.$async;
        this.refs = {};
      }
    };
    exports.SchemaEnv = SchemaEnv;
    function compileSchema(sch) {
      const _sch = getCompilingSchema.call(this, sch);
      if (_sch)
        return _sch;
      const rootId = (0, resolve_1.getFullPath)(this.opts.uriResolver, sch.root.baseId);
      const { es5, lines } = this.opts.code;
      const { ownProperties } = this.opts;
      const gen = new codegen_1.CodeGen(this.scope, { es5, lines, ownProperties });
      let _ValidationError;
      if (sch.$async) {
        _ValidationError = gen.scopeValue("Error", {
          ref: validation_error_1.default,
          code: (0, codegen_1._)`require("ajv/dist/runtime/validation_error").default`
        });
      }
      const validateName = gen.scopeName("validate");
      sch.validateName = validateName;
      const schemaCxt = {
        gen,
        allErrors: this.opts.allErrors,
        data: names_1.default.data,
        parentData: names_1.default.parentData,
        parentDataProperty: names_1.default.parentDataProperty,
        dataNames: [names_1.default.data],
        dataPathArr: [codegen_1.nil],
        // TODO can its length be used as dataLevel if nil is removed?
        dataLevel: 0,
        dataTypes: [],
        definedProperties: /* @__PURE__ */ new Set(),
        topSchemaRef: gen.scopeValue("schema", this.opts.code.source === true ? { ref: sch.schema, code: (0, codegen_1.stringify)(sch.schema) } : { ref: sch.schema }),
        validateName,
        ValidationError: _ValidationError,
        schema: sch.schema,
        schemaEnv: sch,
        rootId,
        baseId: sch.baseId || rootId,
        schemaPath: codegen_1.nil,
        errSchemaPath: sch.schemaPath || (this.opts.jtd ? "" : "#"),
        errorPath: (0, codegen_1._)`""`,
        opts: this.opts,
        self: this
      };
      let sourceCode;
      try {
        this._compilations.add(sch);
        (0, validate_1.validateFunctionCode)(schemaCxt);
        gen.optimize(this.opts.code.optimize);
        const validateCode = gen.toString();
        sourceCode = `${gen.scopeRefs(names_1.default.scope)}return ${validateCode}`;
        if (this.opts.code.process)
          sourceCode = this.opts.code.process(sourceCode, sch);
        const makeValidate = new Function(`${names_1.default.self}`, `${names_1.default.scope}`, sourceCode);
        const validate = makeValidate(this, this.scope.get());
        this.scope.value(validateName, { ref: validate });
        validate.errors = null;
        validate.schema = sch.schema;
        validate.schemaEnv = sch;
        if (sch.$async)
          validate.$async = true;
        if (this.opts.code.source === true) {
          validate.source = { validateName, validateCode, scopeValues: gen._values };
        }
        if (this.opts.unevaluated) {
          const { props, items } = schemaCxt;
          validate.evaluated = {
            props: props instanceof codegen_1.Name ? void 0 : props,
            items: items instanceof codegen_1.Name ? void 0 : items,
            dynamicProps: props instanceof codegen_1.Name,
            dynamicItems: items instanceof codegen_1.Name
          };
          if (validate.source)
            validate.source.evaluated = (0, codegen_1.stringify)(validate.evaluated);
        }
        sch.validate = validate;
        return sch;
      } catch (e) {
        delete sch.validate;
        delete sch.validateName;
        if (sourceCode)
          this.logger.error("Error compiling schema, function code:", sourceCode);
        throw e;
      } finally {
        this._compilations.delete(sch);
      }
    }
    exports.compileSchema = compileSchema;
    function resolveRef(root, baseId, ref) {
      var _a;
      ref = (0, resolve_1.resolveUrl)(this.opts.uriResolver, baseId, ref);
      const schOrFunc = root.refs[ref];
      if (schOrFunc)
        return schOrFunc;
      let _sch = resolve35.call(this, root, ref);
      if (_sch === void 0) {
        const schema2 = (_a = root.localRefs) === null || _a === void 0 ? void 0 : _a[ref];
        const { schemaId } = this.opts;
        if (schema2)
          _sch = new SchemaEnv({ schema: schema2, schemaId, root, baseId });
      }
      if (_sch === void 0)
        return;
      return root.refs[ref] = inlineOrCompile.call(this, _sch);
    }
    exports.resolveRef = resolveRef;
    function inlineOrCompile(sch) {
      if ((0, resolve_1.inlineRef)(sch.schema, this.opts.inlineRefs))
        return sch.schema;
      return sch.validate ? sch : compileSchema.call(this, sch);
    }
    function getCompilingSchema(schEnv) {
      for (const sch of this._compilations) {
        if (sameSchemaEnv(sch, schEnv))
          return sch;
      }
    }
    exports.getCompilingSchema = getCompilingSchema;
    function sameSchemaEnv(s1, s2) {
      return s1.schema === s2.schema && s1.root === s2.root && s1.baseId === s2.baseId;
    }
    function resolve35(root, ref) {
      let sch;
      while (typeof (sch = this.refs[ref]) == "string")
        ref = sch;
      return sch || this.schemas[ref] || resolveSchema.call(this, root, ref);
    }
    function resolveSchema(root, ref) {
      const p = this.opts.uriResolver.parse(ref);
      const refPath = (0, resolve_1._getFullPath)(this.opts.uriResolver, p);
      let baseId = (0, resolve_1.getFullPath)(this.opts.uriResolver, root.baseId, void 0);
      if (Object.keys(root.schema).length > 0 && refPath === baseId) {
        return getJsonPointer.call(this, p, root);
      }
      const id = (0, resolve_1.normalizeId)(refPath);
      const schOrRef = this.refs[id] || this.schemas[id];
      if (typeof schOrRef == "string") {
        const sch = resolveSchema.call(this, root, schOrRef);
        if (typeof (sch === null || sch === void 0 ? void 0 : sch.schema) !== "object")
          return;
        return getJsonPointer.call(this, p, sch);
      }
      if (typeof (schOrRef === null || schOrRef === void 0 ? void 0 : schOrRef.schema) !== "object")
        return;
      if (!schOrRef.validate)
        compileSchema.call(this, schOrRef);
      if (id === (0, resolve_1.normalizeId)(ref)) {
        const { schema: schema2 } = schOrRef;
        const { schemaId } = this.opts;
        const schId = schema2[schemaId];
        if (schId)
          baseId = (0, resolve_1.resolveUrl)(this.opts.uriResolver, baseId, schId);
        return new SchemaEnv({ schema: schema2, schemaId, root, baseId });
      }
      return getJsonPointer.call(this, p, schOrRef);
    }
    exports.resolveSchema = resolveSchema;
    var PREVENT_SCOPE_CHANGE = /* @__PURE__ */ new Set([
      "properties",
      "patternProperties",
      "enum",
      "dependencies",
      "definitions"
    ]);
    function getJsonPointer(parsedRef, { baseId, schema: schema2, root }) {
      var _a;
      if (((_a = parsedRef.fragment) === null || _a === void 0 ? void 0 : _a[0]) !== "/")
        return;
      for (const part of parsedRef.fragment.slice(1).split("/")) {
        if (typeof schema2 === "boolean")
          return;
        const partSchema = schema2[(0, util_1.unescapeFragment)(part)];
        if (partSchema === void 0)
          return;
        schema2 = partSchema;
        const schId = typeof schema2 === "object" && schema2[this.opts.schemaId];
        if (!PREVENT_SCOPE_CHANGE.has(part) && schId) {
          baseId = (0, resolve_1.resolveUrl)(this.opts.uriResolver, baseId, schId);
        }
      }
      let env;
      if (typeof schema2 != "boolean" && schema2.$ref && !(0, util_1.schemaHasRulesButRef)(schema2, this.RULES)) {
        const $ref = (0, resolve_1.resolveUrl)(this.opts.uriResolver, baseId, schema2.$ref);
        env = resolveSchema.call(this, root, $ref);
      }
      const { schemaId } = this.opts;
      env = env || new SchemaEnv({ schema: schema2, schemaId, root, baseId });
      if (env.schema !== env.root.schema)
        return env;
      return void 0;
    }
  }
});

// node_modules/ajv/dist/refs/data.json
var require_data = __commonJS({
  "node_modules/ajv/dist/refs/data.json"(exports, module) {
    module.exports = {
      $id: "https://raw.githubusercontent.com/ajv-validator/ajv/master/lib/refs/data.json#",
      description: "Meta-schema for $data reference (JSON AnySchema extension proposal)",
      type: "object",
      required: ["$data"],
      properties: {
        $data: {
          type: "string",
          anyOf: [{ format: "relative-json-pointer" }, { format: "json-pointer" }]
        }
      },
      additionalProperties: false
    };
  }
});

// node_modules/fast-uri/lib/utils.js
var require_utils = __commonJS({
  "node_modules/fast-uri/lib/utils.js"(exports, module) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var isUUID = RegExp.prototype.test.bind(/^[\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12}$/iu);
    var isIPv4 = RegExp.prototype.test.bind(/^(?:(?:25[0-5]|2[0-4]\d|1\d{2}|[1-9]\d|\d)\.){3}(?:25[0-5]|2[0-4]\d|1\d{2}|[1-9]\d|\d)$/u);
    var isHexPair = RegExp.prototype.test.bind(/^[\da-f]{2}$/iu);
    var isUnreserved = RegExp.prototype.test.bind(/^[\da-z\-._~]$/iu);
    var isPathCharacter = RegExp.prototype.test.bind(/^[\da-z\-._~!$&'()*+,;=:@/]$/iu);
    function stringArrayToHexStripped(input) {
      let acc = "";
      let code = 0;
      let i = 0;
      for (i = 0; i < input.length; i++) {
        code = input[i].charCodeAt(0);
        if (code === 48) {
          continue;
        }
        if (!(code >= 48 && code <= 57 || code >= 65 && code <= 70 || code >= 97 && code <= 102)) {
          return "";
        }
        acc += input[i];
        break;
      }
      for (i += 1; i < input.length; i++) {
        code = input[i].charCodeAt(0);
        if (!(code >= 48 && code <= 57 || code >= 65 && code <= 70 || code >= 97 && code <= 102)) {
          return "";
        }
        acc += input[i];
      }
      return acc;
    }
    var nonSimpleDomain = RegExp.prototype.test.bind(/[^!"$&'()*+,\-.;=_`a-z{}~]/u);
    function consumeIsZone(buffer) {
      buffer.length = 0;
      return true;
    }
    function consumeHextets(buffer, address, output) {
      if (buffer.length) {
        const hex = stringArrayToHexStripped(buffer);
        if (hex !== "") {
          address.push(hex);
        } else {
          output.error = true;
          return false;
        }
        buffer.length = 0;
      }
      return true;
    }
    function getIPV6(input) {
      let tokenCount = 0;
      const output = { error: false, address: "", zone: "" };
      const address = [];
      const buffer = [];
      let endipv6Encountered = false;
      let endIpv6 = false;
      let consume = consumeHextets;
      for (let i = 0; i < input.length; i++) {
        const cursor = input[i];
        if (cursor === "[" || cursor === "]") {
          continue;
        }
        if (cursor === ":") {
          if (endipv6Encountered === true) {
            endIpv6 = true;
          }
          if (!consume(buffer, address, output)) {
            break;
          }
          if (++tokenCount > 7) {
            output.error = true;
            break;
          }
          if (i > 0 && input[i - 1] === ":") {
            endipv6Encountered = true;
          }
          address.push(":");
          continue;
        } else if (cursor === "%") {
          if (!consume(buffer, address, output)) {
            break;
          }
          consume = consumeIsZone;
        } else {
          buffer.push(cursor);
          continue;
        }
      }
      if (buffer.length) {
        if (consume === consumeIsZone) {
          output.zone = buffer.join("");
        } else if (endIpv6) {
          address.push(buffer.join(""));
        } else {
          address.push(stringArrayToHexStripped(buffer));
        }
      }
      output.address = address.join("");
      return output;
    }
    function normalizeIPv6(host) {
      if (findToken(host, ":") < 2) {
        return { host, isIPV6: false };
      }
      const ipv6 = getIPV6(host);
      if (!ipv6.error) {
        let newHost = ipv6.address;
        let escapedHost = ipv6.address;
        if (ipv6.zone) {
          newHost += "%" + ipv6.zone;
          escapedHost += "%25" + ipv6.zone;
        }
        return { host: newHost, isIPV6: true, escapedHost };
      } else {
        return { host, isIPV6: false };
      }
    }
    function findToken(str, token) {
      let ind = 0;
      for (let i = 0; i < str.length; i++) {
        if (str[i] === token) ind++;
      }
      return ind;
    }
    function removeDotSegments(path2) {
      let input = path2;
      const output = [];
      let nextSlash = -1;
      let len = 0;
      while (len = input.length) {
        if (len === 1) {
          if (input === ".") {
            break;
          } else if (input === "/") {
            output.push("/");
            break;
          } else {
            output.push(input);
            break;
          }
        } else if (len === 2) {
          if (input[0] === ".") {
            if (input[1] === ".") {
              break;
            } else if (input[1] === "/") {
              input = input.slice(2);
              continue;
            }
          } else if (input[0] === "/") {
            if (input[1] === "." || input[1] === "/") {
              output.push("/");
              break;
            }
          }
        } else if (len === 3) {
          if (input === "/..") {
            if (output.length !== 0) {
              output.pop();
            }
            output.push("/");
            break;
          }
        }
        if (input[0] === ".") {
          if (input[1] === ".") {
            if (input[2] === "/") {
              input = input.slice(3);
              continue;
            }
          } else if (input[1] === "/") {
            input = input.slice(2);
            continue;
          }
        } else if (input[0] === "/") {
          if (input[1] === ".") {
            if (input[2] === "/") {
              input = input.slice(2);
              continue;
            } else if (input[2] === ".") {
              if (input[3] === "/") {
                input = input.slice(3);
                if (output.length !== 0) {
                  output.pop();
                }
                continue;
              }
            }
          }
        }
        if ((nextSlash = input.indexOf("/", 1)) === -1) {
          output.push(input);
          break;
        } else {
          output.push(input.slice(0, nextSlash));
          input = input.slice(nextSlash);
        }
      }
      return output.join("");
    }
    var HOST_DELIMS = { "@": "%40", "/": "%2F", "?": "%3F", "#": "%23", ":": "%3A" };
    var HOST_DELIM_RE = /[@/?#:]/g;
    var HOST_DELIM_NO_COLON_RE = /[@/?#]/g;
    function reescapeHostDelimiters(host, isIP) {
      const re = isIP ? HOST_DELIM_NO_COLON_RE : HOST_DELIM_RE;
      re.lastIndex = 0;
      return host.replace(re, (ch) => HOST_DELIMS[ch]);
    }
    function normalizePercentEncoding(input, decodeUnreserved = false) {
      if (input.indexOf("%") === -1) {
        return input;
      }
      let output = "";
      for (let i = 0; i < input.length; i++) {
        if (input[i] === "%" && i + 2 < input.length) {
          const hex = input.slice(i + 1, i + 3);
          if (isHexPair(hex)) {
            const normalizedHex = hex.toUpperCase();
            const decoded2 = String.fromCharCode(parseInt(normalizedHex, 16));
            if (decodeUnreserved && isUnreserved(decoded2)) {
              output += decoded2;
            } else {
              output += "%" + normalizedHex;
            }
            i += 2;
            continue;
          }
        }
        output += input[i];
      }
      return output;
    }
    function normalizePathEncoding(input) {
      let output = "";
      for (let i = 0; i < input.length; i++) {
        if (input[i] === "%" && i + 2 < input.length) {
          const hex = input.slice(i + 1, i + 3);
          if (isHexPair(hex)) {
            const normalizedHex = hex.toUpperCase();
            const decoded2 = String.fromCharCode(parseInt(normalizedHex, 16));
            if (decoded2 !== "." && isUnreserved(decoded2)) {
              output += decoded2;
            } else {
              output += "%" + normalizedHex;
            }
            i += 2;
            continue;
          }
        }
        if (isPathCharacter(input[i])) {
          output += input[i];
        } else {
          output += escape(input[i]);
        }
      }
      return output;
    }
    function escapePreservingEscapes(input) {
      let output = "";
      for (let i = 0; i < input.length; i++) {
        if (input[i] === "%" && i + 2 < input.length) {
          const hex = input.slice(i + 1, i + 3);
          if (isHexPair(hex)) {
            output += "%" + hex.toUpperCase();
            i += 2;
            continue;
          }
        }
        output += escape(input[i]);
      }
      return output;
    }
    function recomposeAuthority(component) {
      const uriTokens = [];
      if (component.userinfo !== void 0) {
        uriTokens.push(component.userinfo);
        uriTokens.push("@");
      }
      if (component.host !== void 0) {
        let host = unescape(component.host);
        if (!isIPv4(host)) {
          const ipV6res = normalizeIPv6(host);
          if (ipV6res.isIPV6 === true) {
            host = `[${ipV6res.escapedHost}]`;
          } else {
            host = reescapeHostDelimiters(host, false);
          }
        }
        uriTokens.push(host);
      }
      if (typeof component.port === "number" || typeof component.port === "string") {
        uriTokens.push(":");
        uriTokens.push(String(component.port));
      }
      return uriTokens.length ? uriTokens.join("") : void 0;
    }
    module.exports = {
      nonSimpleDomain,
      recomposeAuthority,
      reescapeHostDelimiters,
      normalizePercentEncoding,
      normalizePathEncoding,
      escapePreservingEscapes,
      removeDotSegments,
      isIPv4,
      isUUID,
      normalizeIPv6,
      stringArrayToHexStripped
    };
  }
});

// node_modules/fast-uri/lib/schemes.js
var require_schemes = __commonJS({
  "node_modules/fast-uri/lib/schemes.js"(exports, module) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var { isUUID } = require_utils();
    var URN_REG = /([\da-z][\d\-a-z]{0,31}):((?:[\w!$'()*+,\-.:;=@]|%[\da-f]{2})+)/iu;
    var supportedSchemeNames = (
      /** @type {const} */
      [
        "http",
        "https",
        "ws",
        "wss",
        "urn",
        "urn:uuid"
      ]
    );
    function isValidSchemeName(name) {
      return supportedSchemeNames.indexOf(
        /** @type {*} */
        name
      ) !== -1;
    }
    function wsIsSecure(wsComponent) {
      if (wsComponent.secure === true) {
        return true;
      } else if (wsComponent.secure === false) {
        return false;
      } else if (wsComponent.scheme) {
        return wsComponent.scheme.length === 3 && (wsComponent.scheme[0] === "w" || wsComponent.scheme[0] === "W") && (wsComponent.scheme[1] === "s" || wsComponent.scheme[1] === "S") && (wsComponent.scheme[2] === "s" || wsComponent.scheme[2] === "S");
      } else {
        return false;
      }
    }
    function httpParse(component) {
      if (!component.host) {
        component.error = component.error || "HTTP URIs must have a host.";
      }
      return component;
    }
    function httpSerialize(component) {
      const secure = String(component.scheme).toLowerCase() === "https";
      if (component.port === (secure ? 443 : 80) || component.port === "") {
        component.port = void 0;
      }
      if (!component.path) {
        component.path = "/";
      }
      return component;
    }
    function wsParse(wsComponent) {
      wsComponent.secure = wsIsSecure(wsComponent);
      wsComponent.resourceName = (wsComponent.path || "/") + (wsComponent.query ? "?" + wsComponent.query : "");
      wsComponent.path = void 0;
      wsComponent.query = void 0;
      return wsComponent;
    }
    function wsSerialize(wsComponent) {
      if (wsComponent.port === (wsIsSecure(wsComponent) ? 443 : 80) || wsComponent.port === "") {
        wsComponent.port = void 0;
      }
      if (typeof wsComponent.secure === "boolean") {
        wsComponent.scheme = wsComponent.secure ? "wss" : "ws";
        wsComponent.secure = void 0;
      }
      if (wsComponent.resourceName) {
        const [path2, query] = wsComponent.resourceName.split("?");
        wsComponent.path = path2 && path2 !== "/" ? path2 : void 0;
        wsComponent.query = query;
        wsComponent.resourceName = void 0;
      }
      wsComponent.fragment = void 0;
      return wsComponent;
    }
    function urnParse(urnComponent, options) {
      if (!urnComponent.path) {
        urnComponent.error = "URN can not be parsed";
        return urnComponent;
      }
      const matches = urnComponent.path.match(URN_REG);
      if (matches) {
        const scheme = options.scheme || urnComponent.scheme || "urn";
        urnComponent.nid = matches[1].toLowerCase();
        urnComponent.nss = matches[2];
        const urnScheme = `${scheme}:${options.nid || urnComponent.nid}`;
        const schemeHandler = getSchemeHandler(urnScheme);
        urnComponent.path = void 0;
        if (schemeHandler) {
          urnComponent = schemeHandler.parse(urnComponent, options);
        }
      } else {
        urnComponent.error = urnComponent.error || "URN can not be parsed.";
      }
      return urnComponent;
    }
    function urnSerialize(urnComponent, options) {
      if (urnComponent.nid === void 0) {
        throw new Error("URN without nid cannot be serialized");
      }
      const scheme = options.scheme || urnComponent.scheme || "urn";
      const nid = urnComponent.nid.toLowerCase();
      const urnScheme = `${scheme}:${options.nid || nid}`;
      const schemeHandler = getSchemeHandler(urnScheme);
      if (schemeHandler) {
        urnComponent = schemeHandler.serialize(urnComponent, options);
      }
      const uriComponent = urnComponent;
      const nss = urnComponent.nss;
      uriComponent.path = `${nid || options.nid}:${nss}`;
      options.skipEscape = true;
      return uriComponent;
    }
    function urnuuidParse(urnComponent, options) {
      const uuidComponent = urnComponent;
      uuidComponent.uuid = uuidComponent.nss;
      uuidComponent.nss = void 0;
      if (!options.tolerant && (!uuidComponent.uuid || !isUUID(uuidComponent.uuid))) {
        uuidComponent.error = uuidComponent.error || "UUID is not valid.";
      }
      return uuidComponent;
    }
    function urnuuidSerialize(uuidComponent) {
      const urnComponent = uuidComponent;
      urnComponent.nss = (uuidComponent.uuid || "").toLowerCase();
      return urnComponent;
    }
    var http = (
      /** @type {SchemeHandler} */
      {
        scheme: "http",
        domainHost: true,
        parse: httpParse,
        serialize: httpSerialize
      }
    );
    var https = (
      /** @type {SchemeHandler} */
      {
        scheme: "https",
        domainHost: http.domainHost,
        parse: httpParse,
        serialize: httpSerialize
      }
    );
    var ws = (
      /** @type {SchemeHandler} */
      {
        scheme: "ws",
        domainHost: true,
        parse: wsParse,
        serialize: wsSerialize
      }
    );
    var wss = (
      /** @type {SchemeHandler} */
      {
        scheme: "wss",
        domainHost: ws.domainHost,
        parse: ws.parse,
        serialize: ws.serialize
      }
    );
    var urn = (
      /** @type {SchemeHandler} */
      {
        scheme: "urn",
        parse: urnParse,
        serialize: urnSerialize,
        skipNormalize: true
      }
    );
    var urnuuid = (
      /** @type {SchemeHandler} */
      {
        scheme: "urn:uuid",
        parse: urnuuidParse,
        serialize: urnuuidSerialize,
        skipNormalize: true
      }
    );
    var SCHEMES = (
      /** @type {Record<SchemeName, SchemeHandler>} */
      {
        http,
        https,
        ws,
        wss,
        urn,
        "urn:uuid": urnuuid
      }
    );
    Object.setPrototypeOf(SCHEMES, null);
    function getSchemeHandler(scheme) {
      return scheme && (SCHEMES[
        /** @type {SchemeName} */
        scheme
      ] || SCHEMES[
        /** @type {SchemeName} */
        scheme.toLowerCase()
      ]) || void 0;
    }
    module.exports = {
      wsIsSecure,
      SCHEMES,
      isValidSchemeName,
      getSchemeHandler
    };
  }
});

// node_modules/fast-uri/index.js
var require_fast_uri = __commonJS({
  "node_modules/fast-uri/index.js"(exports, module) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var { normalizeIPv6, removeDotSegments, recomposeAuthority, normalizePercentEncoding, normalizePathEncoding, escapePreservingEscapes, reescapeHostDelimiters, isIPv4, nonSimpleDomain } = require_utils();
    var { SCHEMES, getSchemeHandler } = require_schemes();
    function normalize(uri, options) {
      if (typeof uri === "string") {
        uri = /** @type {T} */
        normalizeString(uri, options);
      } else if (typeof uri === "object") {
        uri = /** @type {T} */
        parse(serialize(uri, options), options);
      }
      return uri;
    }
    function resolve35(baseURI, relativeURI, options) {
      const schemelessOptions = options ? Object.assign({ scheme: "null" }, options) : { scheme: "null" };
      const { parsed: baseParsed, malformedAuthorityOrPort: baseMalformed } = parseWithStatus(baseURI, schemelessOptions);
      const { parsed: relativeParsed, malformedAuthorityOrPort: relativeMalformed } = parseWithStatus(relativeURI, schemelessOptions);
      if (baseMalformed || relativeMalformed) {
        throw new Error(baseParsed.error || relativeParsed.error || "URI is malformed.");
      }
      const resolved = resolveComponent(baseParsed, relativeParsed, schemelessOptions, true);
      schemelessOptions.skipEscape = true;
      return serialize(resolved, schemelessOptions);
    }
    function resolveComponent(base, relative8, options, skipNormalization) {
      const target = {};
      if (!skipNormalization) {
        base = parse(serialize(base, options), options);
        relative8 = parse(serialize(relative8, options), options);
      }
      options = options || {};
      if (!options.tolerant && relative8.scheme) {
        target.scheme = relative8.scheme;
        target.userinfo = relative8.userinfo;
        target.host = relative8.host;
        target.port = relative8.port;
        target.path = removeDotSegments(relative8.path || "");
        target.query = relative8.query;
      } else {
        if (relative8.userinfo !== void 0 || relative8.host !== void 0 || relative8.port !== void 0) {
          target.userinfo = relative8.userinfo;
          target.host = relative8.host;
          target.port = relative8.port;
          target.path = removeDotSegments(relative8.path || "");
          target.query = relative8.query;
        } else {
          if (!relative8.path) {
            target.path = base.path;
            if (relative8.query !== void 0) {
              target.query = relative8.query;
            } else {
              target.query = base.query;
            }
          } else {
            if (relative8.path[0] === "/") {
              target.path = removeDotSegments(relative8.path);
            } else {
              if ((base.userinfo !== void 0 || base.host !== void 0 || base.port !== void 0) && !base.path) {
                target.path = "/" + relative8.path;
              } else if (!base.path) {
                target.path = relative8.path;
              } else {
                target.path = base.path.slice(0, base.path.lastIndexOf("/") + 1) + relative8.path;
              }
              target.path = removeDotSegments(target.path);
            }
            target.query = relative8.query;
          }
          target.userinfo = base.userinfo;
          target.host = base.host;
          target.port = base.port;
        }
        target.scheme = base.scheme;
      }
      target.fragment = relative8.fragment;
      return target;
    }
    function equal(uriA, uriB, options) {
      const normalizedA = normalizeComparableURI(uriA, options);
      const normalizedB = normalizeComparableURI(uriB, options);
      return normalizedA !== void 0 && normalizedB !== void 0 && normalizedA.toLowerCase() === normalizedB.toLowerCase();
    }
    function serialize(cmpts, opts) {
      const component = {
        host: cmpts.host,
        scheme: cmpts.scheme,
        userinfo: cmpts.userinfo,
        port: cmpts.port,
        path: cmpts.path,
        query: cmpts.query,
        nid: cmpts.nid,
        nss: cmpts.nss,
        uuid: cmpts.uuid,
        fragment: cmpts.fragment,
        reference: cmpts.reference,
        resourceName: cmpts.resourceName,
        secure: cmpts.secure,
        error: ""
      };
      const options = Object.assign({}, opts);
      const uriTokens = [];
      const schemeHandler = getSchemeHandler(options.scheme || component.scheme);
      if (schemeHandler && schemeHandler.serialize) schemeHandler.serialize(component, options);
      if (component.path !== void 0) {
        if (!options.skipEscape) {
          component.path = escapePreservingEscapes(component.path);
          if (component.scheme !== void 0) {
            component.path = component.path.split("%3A").join(":");
          }
        } else {
          component.path = normalizePercentEncoding(component.path);
        }
      }
      if (options.reference !== "suffix" && component.scheme) {
        uriTokens.push(component.scheme, ":");
      }
      const authority = recomposeAuthority(component);
      if (authority !== void 0) {
        if (options.reference !== "suffix") {
          uriTokens.push("//");
        }
        uriTokens.push(authority);
        if (component.path && component.path[0] !== "/") {
          uriTokens.push("/");
        }
      }
      if (component.path !== void 0) {
        let s = component.path;
        if (!options.absolutePath && (!schemeHandler || !schemeHandler.absolutePath)) {
          s = removeDotSegments(s);
        }
        if (authority === void 0 && s[0] === "/" && s[1] === "/") {
          s = "/%2F" + s.slice(2);
        }
        uriTokens.push(s);
      }
      if (component.query !== void 0) {
        uriTokens.push("?", component.query);
      }
      if (component.fragment !== void 0) {
        uriTokens.push("#", component.fragment);
      }
      return uriTokens.join("");
    }
    var URI_PARSE = /^(?:([^#/:?]+):)?(?:\/\/((?:([^#/?@]*)@)?(\[[^#/?\]]+\]|[^#/:?]*)(?::(\d*))?))?([^#?]*)(?:\?([^#]*))?(?:#((?:.|[\n\r])*))?/u;
    var AUTHORITY_PREFIX = /^(?:[^#/:?]+:)?\/\/([^/?#]*)/;
    var AUTHORITY_INTRODUCER_REGION = /^(?:[^#/:?]+:)?([/\\\t\n\r]*)/;
    function getParseError(parsed, matches) {
      if (matches[2] !== void 0 && parsed.path && parsed.path[0] !== "/") {
        return 'URI path must start with "/" when authority is present.';
      }
      if (typeof parsed.port === "number" && (parsed.port < 0 || parsed.port > 65535)) {
        return "URI port is malformed.";
      }
      return void 0;
    }
    function parseWithStatus(uri, opts) {
      const options = Object.assign({}, opts);
      const parsed = {
        scheme: void 0,
        userinfo: void 0,
        host: "",
        port: void 0,
        path: "",
        query: void 0,
        fragment: void 0
      };
      let malformedAuthorityOrPort = false;
      let isIP = false;
      if (options.reference === "suffix") {
        if (options.scheme) {
          uri = options.scheme + ":" + uri;
        } else {
          uri = "//" + uri;
        }
      }
      const authorityMatch = uri.match(AUTHORITY_PREFIX);
      if (authorityMatch !== null && authorityMatch[1].indexOf("\\") !== -1) {
        parsed.error = "URI authority must not contain a literal backslash.";
        malformedAuthorityOrPort = true;
      }
      const introducerMatch = uri.match(AUTHORITY_INTRODUCER_REGION);
      if (introducerMatch !== null) {
        const region = introducerMatch[1];
        const normalizedRegion = region.replace(/[\t\n\r]/g, "");
        if (normalizedRegion.length >= 2) {
          if (normalizedRegion.slice(0, 2) !== "//") {
            parsed.error = parsed.error || "URI authority must not contain a literal backslash.";
            malformedAuthorityOrPort = true;
          } else if (region.length !== normalizedRegion.length) {
            parsed.error = parsed.error || "URI authority introducer must not contain whitespace.";
            malformedAuthorityOrPort = true;
          }
        }
      }
      const matches = uri.match(URI_PARSE);
      if (matches) {
        parsed.scheme = matches[1];
        parsed.userinfo = matches[3];
        parsed.host = matches[4];
        parsed.port = parseInt(matches[5], 10);
        parsed.path = matches[6] || "";
        parsed.query = matches[7];
        parsed.fragment = matches[8];
        if (isNaN(parsed.port)) {
          parsed.port = matches[5];
        }
        const parseError = getParseError(parsed, matches);
        if (parseError !== void 0) {
          parsed.error = parsed.error || parseError;
          malformedAuthorityOrPort = true;
        }
        if (parsed.host) {
          const ipv4result = isIPv4(parsed.host);
          if (ipv4result === false) {
            const ipv6result = normalizeIPv6(parsed.host);
            parsed.host = ipv6result.host.toLowerCase();
            isIP = ipv6result.isIPV6;
          } else {
            isIP = true;
          }
        }
        if (parsed.scheme === void 0 && parsed.userinfo === void 0 && parsed.host === void 0 && parsed.port === void 0 && parsed.query === void 0 && !parsed.path) {
          parsed.reference = "same-document";
        } else if (parsed.scheme === void 0) {
          parsed.reference = "relative";
        } else if (parsed.fragment === void 0) {
          parsed.reference = "absolute";
        } else {
          parsed.reference = "uri";
        }
        if (options.reference && options.reference !== "suffix" && options.reference !== parsed.reference) {
          parsed.error = parsed.error || "URI is not a " + options.reference + " reference.";
        }
        const schemeHandler = getSchemeHandler(options.scheme || parsed.scheme);
        if (!options.unicodeSupport && (!schemeHandler || !schemeHandler.unicodeSupport)) {
          if (parsed.host && (options.domainHost || schemeHandler && schemeHandler.domainHost) && isIP === false && nonSimpleDomain(parsed.host)) {
            try {
              parsed.host = new URL("http://" + parsed.host).hostname;
            } catch (e) {
              parsed.error = parsed.error || "Host's domain name can not be converted to ASCII: " + e;
            }
          }
        }
        if (!schemeHandler || schemeHandler && !schemeHandler.skipNormalize) {
          if (uri.indexOf("%") !== -1) {
            if (parsed.scheme !== void 0) {
              parsed.scheme = unescape(parsed.scheme);
            }
            if (parsed.host !== void 0) {
              parsed.host = reescapeHostDelimiters(unescape(parsed.host), isIP);
            }
          }
          if (parsed.path) {
            parsed.path = normalizePathEncoding(parsed.path);
          }
          if (parsed.fragment) {
            try {
              parsed.fragment = encodeURI(decodeURIComponent(parsed.fragment));
            } catch {
              parsed.error = parsed.error || "URI malformed";
            }
          }
        }
        if (schemeHandler && schemeHandler.parse) {
          schemeHandler.parse(parsed, options);
        }
      } else {
        parsed.error = parsed.error || "URI can not be parsed.";
      }
      return { parsed, malformedAuthorityOrPort };
    }
    function parse(uri, opts) {
      return parseWithStatus(uri, opts).parsed;
    }
    function normalizeString(uri, opts) {
      return normalizeStringWithStatus(uri, opts).normalized;
    }
    function normalizeStringWithStatus(uri, opts) {
      const { parsed, malformedAuthorityOrPort } = parseWithStatus(uri, opts);
      return {
        normalized: malformedAuthorityOrPort ? uri : serialize(parsed, opts),
        malformedAuthorityOrPort
      };
    }
    function normalizeComparableURI(uri, opts) {
      if (typeof uri === "string") {
        const { normalized, malformedAuthorityOrPort } = normalizeStringWithStatus(uri, opts);
        return malformedAuthorityOrPort ? void 0 : normalized;
      }
      if (typeof uri === "object") {
        return serialize(uri, opts);
      }
    }
    var fastUri = {
      SCHEMES,
      normalize,
      resolve: resolve35,
      resolveComponent,
      equal,
      serialize,
      parse
    };
    module.exports = fastUri;
    module.exports.default = fastUri;
    module.exports.fastUri = fastUri;
  }
});

// node_modules/ajv/dist/runtime/uri.js
var require_uri = __commonJS({
  "node_modules/ajv/dist/runtime/uri.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var uri = require_fast_uri();
    uri.code = 'require("ajv/dist/runtime/uri").default';
    exports.default = uri;
  }
});

// node_modules/ajv/dist/core.js
var require_core = __commonJS({
  "node_modules/ajv/dist/core.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.CodeGen = exports.Name = exports.nil = exports.stringify = exports.str = exports._ = exports.KeywordCxt = void 0;
    var validate_1 = require_validate();
    Object.defineProperty(exports, "KeywordCxt", { enumerable: true, get: function() {
      return validate_1.KeywordCxt;
    } });
    var codegen_1 = require_codegen();
    Object.defineProperty(exports, "_", { enumerable: true, get: function() {
      return codegen_1._;
    } });
    Object.defineProperty(exports, "str", { enumerable: true, get: function() {
      return codegen_1.str;
    } });
    Object.defineProperty(exports, "stringify", { enumerable: true, get: function() {
      return codegen_1.stringify;
    } });
    Object.defineProperty(exports, "nil", { enumerable: true, get: function() {
      return codegen_1.nil;
    } });
    Object.defineProperty(exports, "Name", { enumerable: true, get: function() {
      return codegen_1.Name;
    } });
    Object.defineProperty(exports, "CodeGen", { enumerable: true, get: function() {
      return codegen_1.CodeGen;
    } });
    var validation_error_1 = require_validation_error();
    var ref_error_1 = require_ref_error();
    var rules_1 = require_rules();
    var compile_1 = require_compile();
    var codegen_2 = require_codegen();
    var resolve_1 = require_resolve();
    var dataType_1 = require_dataType();
    var util_1 = require_util();
    var $dataRefSchema = require_data();
    var uri_1 = require_uri();
    var defaultRegExp = (str, flags) => new RegExp(str, flags);
    defaultRegExp.code = "new RegExp";
    var META_IGNORE_OPTIONS = ["removeAdditional", "useDefaults", "coerceTypes"];
    var EXT_SCOPE_NAMES = /* @__PURE__ */ new Set([
      "validate",
      "serialize",
      "parse",
      "wrapper",
      "root",
      "schema",
      "keyword",
      "pattern",
      "formats",
      "validate$data",
      "func",
      "obj",
      "Error"
    ]);
    var removedOptions = {
      errorDataPath: "",
      format: "`validateFormats: false` can be used instead.",
      nullable: '"nullable" keyword is supported by default.',
      jsonPointers: "Deprecated jsPropertySyntax can be used instead.",
      extendRefs: "Deprecated ignoreKeywordsWithRef can be used instead.",
      missingRefs: "Pass empty schema with $id that should be ignored to ajv.addSchema.",
      processCode: "Use option `code: {process: (code, schemaEnv: object) => string}`",
      sourceCode: "Use option `code: {source: true}`",
      strictDefaults: "It is default now, see option `strict`.",
      strictKeywords: "It is default now, see option `strict`.",
      uniqueItems: '"uniqueItems" keyword is always validated.',
      unknownFormats: "Disable strict mode or pass `true` to `ajv.addFormat` (or `formats` option).",
      cache: "Map is used as cache, schema object as key.",
      serialize: "Map is used as cache, schema object as key.",
      ajvErrors: "It is default now."
    };
    var deprecatedOptions = {
      ignoreKeywordsWithRef: "",
      jsPropertySyntax: "",
      unicode: '"minLength"/"maxLength" account for unicode characters by default.'
    };
    var MAX_EXPRESSION = 200;
    function requiredOptions(o) {
      var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t, _u, _v, _w, _x, _y, _z, _0;
      const s = o.strict;
      const _optz = (_a = o.code) === null || _a === void 0 ? void 0 : _a.optimize;
      const optimize = _optz === true || _optz === void 0 ? 1 : _optz || 0;
      const regExp = (_c = (_b = o.code) === null || _b === void 0 ? void 0 : _b.regExp) !== null && _c !== void 0 ? _c : defaultRegExp;
      const uriResolver = (_d = o.uriResolver) !== null && _d !== void 0 ? _d : uri_1.default;
      return {
        strictSchema: (_f = (_e = o.strictSchema) !== null && _e !== void 0 ? _e : s) !== null && _f !== void 0 ? _f : true,
        strictNumbers: (_h = (_g = o.strictNumbers) !== null && _g !== void 0 ? _g : s) !== null && _h !== void 0 ? _h : true,
        strictTypes: (_k = (_j = o.strictTypes) !== null && _j !== void 0 ? _j : s) !== null && _k !== void 0 ? _k : "log",
        strictTuples: (_m = (_l = o.strictTuples) !== null && _l !== void 0 ? _l : s) !== null && _m !== void 0 ? _m : "log",
        strictRequired: (_p = (_o = o.strictRequired) !== null && _o !== void 0 ? _o : s) !== null && _p !== void 0 ? _p : false,
        code: o.code ? { ...o.code, optimize, regExp } : { optimize, regExp },
        loopRequired: (_q = o.loopRequired) !== null && _q !== void 0 ? _q : MAX_EXPRESSION,
        loopEnum: (_r = o.loopEnum) !== null && _r !== void 0 ? _r : MAX_EXPRESSION,
        meta: (_s = o.meta) !== null && _s !== void 0 ? _s : true,
        messages: (_t = o.messages) !== null && _t !== void 0 ? _t : true,
        inlineRefs: (_u = o.inlineRefs) !== null && _u !== void 0 ? _u : true,
        schemaId: (_v = o.schemaId) !== null && _v !== void 0 ? _v : "$id",
        addUsedSchema: (_w = o.addUsedSchema) !== null && _w !== void 0 ? _w : true,
        validateSchema: (_x = o.validateSchema) !== null && _x !== void 0 ? _x : true,
        validateFormats: (_y = o.validateFormats) !== null && _y !== void 0 ? _y : true,
        unicodeRegExp: (_z = o.unicodeRegExp) !== null && _z !== void 0 ? _z : true,
        int32range: (_0 = o.int32range) !== null && _0 !== void 0 ? _0 : true,
        uriResolver
      };
    }
    var Ajv = class {
      constructor(opts = {}) {
        this.schemas = {};
        this.refs = {};
        this.formats = /* @__PURE__ */ Object.create(null);
        this._compilations = /* @__PURE__ */ new Set();
        this._loading = {};
        this._cache = /* @__PURE__ */ new Map();
        opts = this.opts = { ...opts, ...requiredOptions(opts) };
        const { es5, lines } = this.opts.code;
        this.scope = new codegen_2.ValueScope({ scope: {}, prefixes: EXT_SCOPE_NAMES, es5, lines });
        this.logger = getLogger(opts.logger);
        const formatOpt = opts.validateFormats;
        opts.validateFormats = false;
        this.RULES = (0, rules_1.getRules)();
        checkOptions.call(this, removedOptions, opts, "NOT SUPPORTED");
        checkOptions.call(this, deprecatedOptions, opts, "DEPRECATED", "warn");
        this._metaOpts = getMetaSchemaOptions.call(this);
        if (opts.formats)
          addInitialFormats.call(this);
        this._addVocabularies();
        this._addDefaultMetaSchema();
        if (opts.keywords)
          addInitialKeywords.call(this, opts.keywords);
        if (typeof opts.meta == "object")
          this.addMetaSchema(opts.meta);
        addInitialSchemas.call(this);
        opts.validateFormats = formatOpt;
      }
      _addVocabularies() {
        this.addKeyword("$async");
      }
      _addDefaultMetaSchema() {
        const { $data, meta, schemaId } = this.opts;
        let _dataRefSchema = $dataRefSchema;
        if (schemaId === "id") {
          _dataRefSchema = { ...$dataRefSchema };
          _dataRefSchema.id = _dataRefSchema.$id;
          delete _dataRefSchema.$id;
        }
        if (meta && $data)
          this.addMetaSchema(_dataRefSchema, _dataRefSchema[schemaId], false);
      }
      defaultMeta() {
        const { meta, schemaId } = this.opts;
        return this.opts.defaultMeta = typeof meta == "object" ? meta[schemaId] || meta : void 0;
      }
      validate(schemaKeyRef, data) {
        let v;
        if (typeof schemaKeyRef == "string") {
          v = this.getSchema(schemaKeyRef);
          if (!v)
            throw new Error(`no schema with key or ref "${schemaKeyRef}"`);
        } else {
          v = this.compile(schemaKeyRef);
        }
        const valid = v(data);
        if (!("$async" in v))
          this.errors = v.errors;
        return valid;
      }
      compile(schema2, _meta) {
        const sch = this._addSchema(schema2, _meta);
        return sch.validate || this._compileSchemaEnv(sch);
      }
      compileAsync(schema2, meta) {
        if (typeof this.opts.loadSchema != "function") {
          throw new Error("options.loadSchema should be a function");
        }
        const { loadSchema } = this.opts;
        return runCompileAsync.call(this, schema2, meta);
        async function runCompileAsync(_schema, _meta) {
          await loadMetaSchema.call(this, _schema.$schema);
          const sch = this._addSchema(_schema, _meta);
          return sch.validate || _compileAsync.call(this, sch);
        }
        async function loadMetaSchema($ref) {
          if ($ref && !this.getSchema($ref)) {
            await runCompileAsync.call(this, { $ref }, true);
          }
        }
        async function _compileAsync(sch) {
          try {
            return this._compileSchemaEnv(sch);
          } catch (e) {
            if (!(e instanceof ref_error_1.default))
              throw e;
            checkLoaded.call(this, e);
            await loadMissingSchema.call(this, e.missingSchema);
            return _compileAsync.call(this, sch);
          }
        }
        function checkLoaded({ missingSchema: ref, missingRef }) {
          if (this.refs[ref]) {
            throw new Error(`AnySchema ${ref} is loaded but ${missingRef} cannot be resolved`);
          }
        }
        async function loadMissingSchema(ref) {
          const _schema = await _loadSchema.call(this, ref);
          if (!this.refs[ref])
            await loadMetaSchema.call(this, _schema.$schema);
          if (!this.refs[ref])
            this.addSchema(_schema, ref, meta);
        }
        async function _loadSchema(ref) {
          const p = this._loading[ref];
          if (p)
            return p;
          try {
            return await (this._loading[ref] = loadSchema(ref));
          } finally {
            delete this._loading[ref];
          }
        }
      }
      // Adds schema to the instance
      addSchema(schema2, key, _meta, _validateSchema = this.opts.validateSchema) {
        if (Array.isArray(schema2)) {
          for (const sch of schema2)
            this.addSchema(sch, void 0, _meta, _validateSchema);
          return this;
        }
        let id;
        if (typeof schema2 === "object") {
          const { schemaId } = this.opts;
          id = schema2[schemaId];
          if (id !== void 0 && typeof id != "string") {
            throw new Error(`schema ${schemaId} must be string`);
          }
        }
        key = (0, resolve_1.normalizeId)(key || id);
        this._checkUnique(key);
        this.schemas[key] = this._addSchema(schema2, _meta, key, _validateSchema, true);
        return this;
      }
      // Add schema that will be used to validate other schemas
      // options in META_IGNORE_OPTIONS are alway set to false
      addMetaSchema(schema2, key, _validateSchema = this.opts.validateSchema) {
        this.addSchema(schema2, key, true, _validateSchema);
        return this;
      }
      //  Validate schema against its meta-schema
      validateSchema(schema2, throwOrLogError) {
        if (typeof schema2 == "boolean")
          return true;
        let $schema;
        $schema = schema2.$schema;
        if ($schema !== void 0 && typeof $schema != "string") {
          throw new Error("$schema must be a string");
        }
        $schema = $schema || this.opts.defaultMeta || this.defaultMeta();
        if (!$schema) {
          this.logger.warn("meta-schema not available");
          this.errors = null;
          return true;
        }
        const valid = this.validate($schema, schema2);
        if (!valid && throwOrLogError) {
          const message = "schema is invalid: " + this.errorsText();
          if (this.opts.validateSchema === "log")
            this.logger.error(message);
          else
            throw new Error(message);
        }
        return valid;
      }
      // Get compiled schema by `key` or `ref`.
      // (`key` that was passed to `addSchema` or full schema reference - `schema.$id` or resolved id)
      getSchema(keyRef) {
        let sch;
        while (typeof (sch = getSchEnv.call(this, keyRef)) == "string")
          keyRef = sch;
        if (sch === void 0) {
          const { schemaId } = this.opts;
          const root = new compile_1.SchemaEnv({ schema: {}, schemaId });
          sch = compile_1.resolveSchema.call(this, root, keyRef);
          if (!sch)
            return;
          this.refs[keyRef] = sch;
        }
        return sch.validate || this._compileSchemaEnv(sch);
      }
      // Remove cached schema(s).
      // If no parameter is passed all schemas but meta-schemas are removed.
      // If RegExp is passed all schemas with key/id matching pattern but meta-schemas are removed.
      // Even if schema is referenced by other schemas it still can be removed as other schemas have local references.
      removeSchema(schemaKeyRef) {
        if (schemaKeyRef instanceof RegExp) {
          this._removeAllSchemas(this.schemas, schemaKeyRef);
          this._removeAllSchemas(this.refs, schemaKeyRef);
          return this;
        }
        switch (typeof schemaKeyRef) {
          case "undefined":
            this._removeAllSchemas(this.schemas);
            this._removeAllSchemas(this.refs);
            this._cache.clear();
            return this;
          case "string": {
            const sch = getSchEnv.call(this, schemaKeyRef);
            if (typeof sch == "object")
              this._cache.delete(sch.schema);
            delete this.schemas[schemaKeyRef];
            delete this.refs[schemaKeyRef];
            return this;
          }
          case "object": {
            const cacheKey = schemaKeyRef;
            this._cache.delete(cacheKey);
            let id = schemaKeyRef[this.opts.schemaId];
            if (id) {
              id = (0, resolve_1.normalizeId)(id);
              delete this.schemas[id];
              delete this.refs[id];
            }
            return this;
          }
          default:
            throw new Error("ajv.removeSchema: invalid parameter");
        }
      }
      // add "vocabulary" - a collection of keywords
      addVocabulary(definitions) {
        for (const def of definitions)
          this.addKeyword(def);
        return this;
      }
      addKeyword(kwdOrDef, def) {
        let keyword;
        if (typeof kwdOrDef == "string") {
          keyword = kwdOrDef;
          if (typeof def == "object") {
            this.logger.warn("these parameters are deprecated, see docs for addKeyword");
            def.keyword = keyword;
          }
        } else if (typeof kwdOrDef == "object" && def === void 0) {
          def = kwdOrDef;
          keyword = def.keyword;
          if (Array.isArray(keyword) && !keyword.length) {
            throw new Error("addKeywords: keyword must be string or non-empty array");
          }
        } else {
          throw new Error("invalid addKeywords parameters");
        }
        checkKeyword.call(this, keyword, def);
        if (!def) {
          (0, util_1.eachItem)(keyword, (kwd) => addRule.call(this, kwd));
          return this;
        }
        keywordMetaschema.call(this, def);
        const definition = {
          ...def,
          type: (0, dataType_1.getJSONTypes)(def.type),
          schemaType: (0, dataType_1.getJSONTypes)(def.schemaType)
        };
        (0, util_1.eachItem)(keyword, definition.type.length === 0 ? (k) => addRule.call(this, k, definition) : (k) => definition.type.forEach((t) => addRule.call(this, k, definition, t)));
        return this;
      }
      getKeyword(keyword) {
        const rule = this.RULES.all[keyword];
        return typeof rule == "object" ? rule.definition : !!rule;
      }
      // Remove keyword
      removeKeyword(keyword) {
        const { RULES } = this;
        delete RULES.keywords[keyword];
        delete RULES.all[keyword];
        for (const group of RULES.rules) {
          const i = group.rules.findIndex((rule) => rule.keyword === keyword);
          if (i >= 0)
            group.rules.splice(i, 1);
        }
        return this;
      }
      // Add format
      addFormat(name, format) {
        if (typeof format == "string")
          format = new RegExp(format);
        this.formats[name] = format;
        return this;
      }
      errorsText(errors2 = this.errors, { separator = ", ", dataVar = "data" } = {}) {
        if (!errors2 || errors2.length === 0)
          return "No errors";
        return errors2.map((e) => `${dataVar}${e.instancePath} ${e.message}`).reduce((text, msg) => text + separator + msg);
      }
      $dataMetaSchema(metaSchema, keywordsJsonPointers) {
        const rules = this.RULES.all;
        metaSchema = JSON.parse(JSON.stringify(metaSchema));
        for (const jsonPointer of keywordsJsonPointers) {
          const segments = jsonPointer.split("/").slice(1);
          let keywords = metaSchema;
          for (const seg of segments)
            keywords = keywords[seg];
          for (const key in rules) {
            const rule = rules[key];
            if (typeof rule != "object")
              continue;
            const { $data } = rule.definition;
            const schema2 = keywords[key];
            if ($data && schema2)
              keywords[key] = schemaOrData(schema2);
          }
        }
        return metaSchema;
      }
      _removeAllSchemas(schemas, regex) {
        for (const keyRef in schemas) {
          const sch = schemas[keyRef];
          if (!regex || regex.test(keyRef)) {
            if (typeof sch == "string") {
              delete schemas[keyRef];
            } else if (sch && !sch.meta) {
              this._cache.delete(sch.schema);
              delete schemas[keyRef];
            }
          }
        }
      }
      _addSchema(schema2, meta, baseId, validateSchema = this.opts.validateSchema, addSchema = this.opts.addUsedSchema) {
        let id;
        const { schemaId } = this.opts;
        if (typeof schema2 == "object") {
          id = schema2[schemaId];
        } else {
          if (this.opts.jtd)
            throw new Error("schema must be object");
          else if (typeof schema2 != "boolean")
            throw new Error("schema must be object or boolean");
        }
        let sch = this._cache.get(schema2);
        if (sch !== void 0)
          return sch;
        baseId = (0, resolve_1.normalizeId)(id || baseId);
        const localRefs = resolve_1.getSchemaRefs.call(this, schema2, baseId);
        sch = new compile_1.SchemaEnv({ schema: schema2, schemaId, meta, baseId, localRefs });
        this._cache.set(sch.schema, sch);
        if (addSchema && !baseId.startsWith("#")) {
          if (baseId)
            this._checkUnique(baseId);
          this.refs[baseId] = sch;
        }
        if (validateSchema)
          this.validateSchema(schema2, true);
        return sch;
      }
      _checkUnique(id) {
        if (this.schemas[id] || this.refs[id]) {
          throw new Error(`schema with key or id "${id}" already exists`);
        }
      }
      _compileSchemaEnv(sch) {
        if (sch.meta)
          this._compileMetaSchema(sch);
        else
          compile_1.compileSchema.call(this, sch);
        if (!sch.validate)
          throw new Error("ajv implementation error");
        return sch.validate;
      }
      _compileMetaSchema(sch) {
        const currentOpts = this.opts;
        this.opts = this._metaOpts;
        try {
          compile_1.compileSchema.call(this, sch);
        } finally {
          this.opts = currentOpts;
        }
      }
    };
    Ajv.ValidationError = validation_error_1.default;
    Ajv.MissingRefError = ref_error_1.default;
    exports.default = Ajv;
    function checkOptions(checkOpts, options, msg, log = "error") {
      for (const key in checkOpts) {
        const opt = key;
        if (opt in options)
          this.logger[log](`${msg}: option ${key}. ${checkOpts[opt]}`);
      }
    }
    function getSchEnv(keyRef) {
      keyRef = (0, resolve_1.normalizeId)(keyRef);
      return this.schemas[keyRef] || this.refs[keyRef];
    }
    function addInitialSchemas() {
      const optsSchemas = this.opts.schemas;
      if (!optsSchemas)
        return;
      if (Array.isArray(optsSchemas))
        this.addSchema(optsSchemas);
      else
        for (const key in optsSchemas)
          this.addSchema(optsSchemas[key], key);
    }
    function addInitialFormats() {
      for (const name in this.opts.formats) {
        const format = this.opts.formats[name];
        if (format)
          this.addFormat(name, format);
      }
    }
    function addInitialKeywords(defs) {
      if (Array.isArray(defs)) {
        this.addVocabulary(defs);
        return;
      }
      this.logger.warn("keywords option as map is deprecated, pass array");
      for (const keyword in defs) {
        const def = defs[keyword];
        if (!def.keyword)
          def.keyword = keyword;
        this.addKeyword(def);
      }
    }
    function getMetaSchemaOptions() {
      const metaOpts = { ...this.opts };
      for (const opt of META_IGNORE_OPTIONS)
        delete metaOpts[opt];
      return metaOpts;
    }
    var noLogs = { log() {
    }, warn() {
    }, error() {
    } };
    function getLogger(logger) {
      if (logger === false)
        return noLogs;
      if (logger === void 0)
        return console;
      if (logger.log && logger.warn && logger.error)
        return logger;
      throw new Error("logger must implement log, warn and error methods");
    }
    var KEYWORD_NAME = /^[a-z_$][a-z0-9_$:-]*$/i;
    function checkKeyword(keyword, def) {
      const { RULES } = this;
      (0, util_1.eachItem)(keyword, (kwd) => {
        if (RULES.keywords[kwd])
          throw new Error(`Keyword ${kwd} is already defined`);
        if (!KEYWORD_NAME.test(kwd))
          throw new Error(`Keyword ${kwd} has invalid name`);
      });
      if (!def)
        return;
      if (def.$data && !("code" in def || "validate" in def)) {
        throw new Error('$data keyword must have "code" or "validate" function');
      }
    }
    function addRule(keyword, definition, dataType) {
      var _a;
      const post = definition === null || definition === void 0 ? void 0 : definition.post;
      if (dataType && post)
        throw new Error('keyword with "post" flag cannot have "type"');
      const { RULES } = this;
      let ruleGroup = post ? RULES.post : RULES.rules.find(({ type: t }) => t === dataType);
      if (!ruleGroup) {
        ruleGroup = { type: dataType, rules: [] };
        RULES.rules.push(ruleGroup);
      }
      RULES.keywords[keyword] = true;
      if (!definition)
        return;
      const rule = {
        keyword,
        definition: {
          ...definition,
          type: (0, dataType_1.getJSONTypes)(definition.type),
          schemaType: (0, dataType_1.getJSONTypes)(definition.schemaType)
        }
      };
      if (definition.before)
        addBeforeRule.call(this, ruleGroup, rule, definition.before);
      else
        ruleGroup.rules.push(rule);
      RULES.all[keyword] = rule;
      (_a = definition.implements) === null || _a === void 0 ? void 0 : _a.forEach((kwd) => this.addKeyword(kwd));
    }
    function addBeforeRule(ruleGroup, rule, before) {
      const i = ruleGroup.rules.findIndex((_rule) => _rule.keyword === before);
      if (i >= 0) {
        ruleGroup.rules.splice(i, 0, rule);
      } else {
        ruleGroup.rules.push(rule);
        this.logger.warn(`rule ${before} is not defined`);
      }
    }
    function keywordMetaschema(def) {
      let { metaSchema } = def;
      if (metaSchema === void 0)
        return;
      if (def.$data && this.opts.$data)
        metaSchema = schemaOrData(metaSchema);
      def.validateSchema = this.compile(metaSchema, true);
    }
    var $dataRef = {
      $ref: "https://raw.githubusercontent.com/ajv-validator/ajv/master/lib/refs/data.json#"
    };
    function schemaOrData(schema2) {
      return { anyOf: [schema2, $dataRef] };
    }
  }
});

// node_modules/ajv/dist/vocabularies/core/id.js
var require_id = __commonJS({
  "node_modules/ajv/dist/vocabularies/core/id.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var def = {
      keyword: "id",
      code() {
        throw new Error('NOT SUPPORTED: keyword "id", use "$id" for schema ID');
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/core/ref.js
var require_ref = __commonJS({
  "node_modules/ajv/dist/vocabularies/core/ref.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.callRef = exports.getValidate = void 0;
    var ref_error_1 = require_ref_error();
    var code_1 = require_code2();
    var codegen_1 = require_codegen();
    var names_1 = require_names();
    var compile_1 = require_compile();
    var util_1 = require_util();
    var def = {
      keyword: "$ref",
      schemaType: "string",
      code(cxt) {
        const { gen, schema: $ref, it } = cxt;
        const { baseId, schemaEnv: env, validateName, opts, self } = it;
        const { root } = env;
        if (($ref === "#" || $ref === "#/") && baseId === root.baseId)
          return callRootRef();
        const schOrEnv = compile_1.resolveRef.call(self, root, baseId, $ref);
        if (schOrEnv === void 0)
          throw new ref_error_1.default(it.opts.uriResolver, baseId, $ref);
        if (schOrEnv instanceof compile_1.SchemaEnv)
          return callValidate(schOrEnv);
        return inlineRefSchema(schOrEnv);
        function callRootRef() {
          if (env === root)
            return callRef(cxt, validateName, env, env.$async);
          const rootName = gen.scopeValue("root", { ref: root });
          return callRef(cxt, (0, codegen_1._)`${rootName}.validate`, root, root.$async);
        }
        function callValidate(sch) {
          const v = getValidate(cxt, sch);
          callRef(cxt, v, sch, sch.$async);
        }
        function inlineRefSchema(sch) {
          const schName = gen.scopeValue("schema", opts.code.source === true ? { ref: sch, code: (0, codegen_1.stringify)(sch) } : { ref: sch });
          const valid = gen.name("valid");
          const schCxt = cxt.subschema({
            schema: sch,
            dataTypes: [],
            schemaPath: codegen_1.nil,
            topSchemaRef: schName,
            errSchemaPath: $ref
          }, valid);
          cxt.mergeEvaluated(schCxt);
          cxt.ok(valid);
        }
      }
    };
    function getValidate(cxt, sch) {
      const { gen } = cxt;
      return sch.validate ? gen.scopeValue("validate", { ref: sch.validate }) : (0, codegen_1._)`${gen.scopeValue("wrapper", { ref: sch })}.validate`;
    }
    exports.getValidate = getValidate;
    function callRef(cxt, v, sch, $async) {
      const { gen, it } = cxt;
      const { allErrors, schemaEnv: env, opts } = it;
      const passCxt = opts.passContext ? names_1.default.this : codegen_1.nil;
      if ($async)
        callAsyncRef();
      else
        callSyncRef();
      function callAsyncRef() {
        if (!env.$async)
          throw new Error("async schema referenced by sync schema");
        const valid = gen.let("valid");
        gen.try(() => {
          gen.code((0, codegen_1._)`await ${(0, code_1.callValidateCode)(cxt, v, passCxt)}`);
          addEvaluatedFrom(v);
          if (!allErrors)
            gen.assign(valid, true);
        }, (e) => {
          gen.if((0, codegen_1._)`!(${e} instanceof ${it.ValidationError})`, () => gen.throw(e));
          addErrorsFrom(e);
          if (!allErrors)
            gen.assign(valid, false);
        });
        cxt.ok(valid);
      }
      function callSyncRef() {
        cxt.result((0, code_1.callValidateCode)(cxt, v, passCxt), () => addEvaluatedFrom(v), () => addErrorsFrom(v));
      }
      function addErrorsFrom(source) {
        const errs = (0, codegen_1._)`${source}.errors`;
        gen.assign(names_1.default.vErrors, (0, codegen_1._)`${names_1.default.vErrors} === null ? ${errs} : ${names_1.default.vErrors}.concat(${errs})`);
        gen.assign(names_1.default.errors, (0, codegen_1._)`${names_1.default.vErrors}.length`);
      }
      function addEvaluatedFrom(source) {
        var _a;
        if (!it.opts.unevaluated)
          return;
        const schEvaluated = (_a = sch === null || sch === void 0 ? void 0 : sch.validate) === null || _a === void 0 ? void 0 : _a.evaluated;
        if (it.props !== true) {
          if (schEvaluated && !schEvaluated.dynamicProps) {
            if (schEvaluated.props !== void 0) {
              it.props = util_1.mergeEvaluated.props(gen, schEvaluated.props, it.props);
            }
          } else {
            const props = gen.var("props", (0, codegen_1._)`${source}.evaluated.props`);
            it.props = util_1.mergeEvaluated.props(gen, props, it.props, codegen_1.Name);
          }
        }
        if (it.items !== true) {
          if (schEvaluated && !schEvaluated.dynamicItems) {
            if (schEvaluated.items !== void 0) {
              it.items = util_1.mergeEvaluated.items(gen, schEvaluated.items, it.items);
            }
          } else {
            const items = gen.var("items", (0, codegen_1._)`${source}.evaluated.items`);
            it.items = util_1.mergeEvaluated.items(gen, items, it.items, codegen_1.Name);
          }
        }
      }
    }
    exports.callRef = callRef;
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/core/index.js
var require_core2 = __commonJS({
  "node_modules/ajv/dist/vocabularies/core/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var id_1 = require_id();
    var ref_1 = require_ref();
    var core = [
      "$schema",
      "$id",
      "$defs",
      "$vocabulary",
      { keyword: "$comment" },
      "definitions",
      id_1.default,
      ref_1.default
    ];
    exports.default = core;
  }
});

// node_modules/ajv/dist/vocabularies/validation/limitNumber.js
var require_limitNumber = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/limitNumber.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var ops = codegen_1.operators;
    var KWDs = {
      maximum: { okStr: "<=", ok: ops.LTE, fail: ops.GT },
      minimum: { okStr: ">=", ok: ops.GTE, fail: ops.LT },
      exclusiveMaximum: { okStr: "<", ok: ops.LT, fail: ops.GTE },
      exclusiveMinimum: { okStr: ">", ok: ops.GT, fail: ops.LTE }
    };
    var error = {
      message: ({ keyword, schemaCode }) => (0, codegen_1.str)`must be ${KWDs[keyword].okStr} ${schemaCode}`,
      params: ({ keyword, schemaCode }) => (0, codegen_1._)`{comparison: ${KWDs[keyword].okStr}, limit: ${schemaCode}}`
    };
    var def = {
      keyword: Object.keys(KWDs),
      type: "number",
      schemaType: "number",
      $data: true,
      error,
      code(cxt) {
        const { keyword, data, schemaCode } = cxt;
        cxt.fail$data((0, codegen_1._)`${data} ${KWDs[keyword].fail} ${schemaCode} || isNaN(${data})`);
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/validation/multipleOf.js
var require_multipleOf = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/multipleOf.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var error = {
      message: ({ schemaCode }) => (0, codegen_1.str)`must be multiple of ${schemaCode}`,
      params: ({ schemaCode }) => (0, codegen_1._)`{multipleOf: ${schemaCode}}`
    };
    var def = {
      keyword: "multipleOf",
      type: "number",
      schemaType: "number",
      $data: true,
      error,
      code(cxt) {
        const { gen, data, schemaCode, it } = cxt;
        const prec = it.opts.multipleOfPrecision;
        const res = gen.let("res");
        const invalid = prec ? (0, codegen_1._)`Math.abs(Math.round(${res}) - ${res}) > 1e-${prec}` : (0, codegen_1._)`${res} !== parseInt(${res})`;
        cxt.fail$data((0, codegen_1._)`(${schemaCode} === 0 || (${res} = ${data}/${schemaCode}, ${invalid}))`);
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/runtime/ucs2length.js
var require_ucs2length = __commonJS({
  "node_modules/ajv/dist/runtime/ucs2length.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    function ucs2length(str) {
      const len = str.length;
      let length = 0;
      let pos = 0;
      let value2;
      while (pos < len) {
        length++;
        value2 = str.charCodeAt(pos++);
        if (value2 >= 55296 && value2 <= 56319 && pos < len) {
          value2 = str.charCodeAt(pos);
          if ((value2 & 64512) === 56320)
            pos++;
        }
      }
      return length;
    }
    exports.default = ucs2length;
    ucs2length.code = 'require("ajv/dist/runtime/ucs2length").default';
  }
});

// node_modules/ajv/dist/vocabularies/validation/limitLength.js
var require_limitLength = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/limitLength.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var ucs2length_1 = require_ucs2length();
    var error = {
      message({ keyword, schemaCode }) {
        const comp = keyword === "maxLength" ? "more" : "fewer";
        return (0, codegen_1.str)`must NOT have ${comp} than ${schemaCode} characters`;
      },
      params: ({ schemaCode }) => (0, codegen_1._)`{limit: ${schemaCode}}`
    };
    var def = {
      keyword: ["maxLength", "minLength"],
      type: "string",
      schemaType: "number",
      $data: true,
      error,
      code(cxt) {
        const { keyword, data, schemaCode, it } = cxt;
        const op = keyword === "maxLength" ? codegen_1.operators.GT : codegen_1.operators.LT;
        const len = it.opts.unicode === false ? (0, codegen_1._)`${data}.length` : (0, codegen_1._)`${(0, util_1.useFunc)(cxt.gen, ucs2length_1.default)}(${data})`;
        cxt.fail$data((0, codegen_1._)`${len} ${op} ${schemaCode}`);
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/validation/pattern.js
var require_pattern = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/pattern.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var code_1 = require_code2();
    var util_1 = require_util();
    var codegen_1 = require_codegen();
    var error = {
      message: ({ schemaCode }) => (0, codegen_1.str)`must match pattern "${schemaCode}"`,
      params: ({ schemaCode }) => (0, codegen_1._)`{pattern: ${schemaCode}}`
    };
    var def = {
      keyword: "pattern",
      type: "string",
      schemaType: "string",
      $data: true,
      error,
      code(cxt) {
        const { gen, data, $data, schema: schema2, schemaCode, it } = cxt;
        const u = it.opts.unicodeRegExp ? "u" : "";
        if ($data) {
          const { regExp } = it.opts.code;
          const regExpCode = regExp.code === "new RegExp" ? (0, codegen_1._)`new RegExp` : (0, util_1.useFunc)(gen, regExp);
          const valid = gen.let("valid");
          gen.try(() => gen.assign(valid, (0, codegen_1._)`${regExpCode}(${schemaCode}, ${u}).test(${data})`), () => gen.assign(valid, false));
          cxt.fail$data((0, codegen_1._)`!${valid}`);
        } else {
          const regExp = (0, code_1.usePattern)(cxt, schema2);
          cxt.fail$data((0, codegen_1._)`!${regExp}.test(${data})`);
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/validation/limitProperties.js
var require_limitProperties = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/limitProperties.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var error = {
      message({ keyword, schemaCode }) {
        const comp = keyword === "maxProperties" ? "more" : "fewer";
        return (0, codegen_1.str)`must NOT have ${comp} than ${schemaCode} properties`;
      },
      params: ({ schemaCode }) => (0, codegen_1._)`{limit: ${schemaCode}}`
    };
    var def = {
      keyword: ["maxProperties", "minProperties"],
      type: "object",
      schemaType: "number",
      $data: true,
      error,
      code(cxt) {
        const { keyword, data, schemaCode } = cxt;
        const op = keyword === "maxProperties" ? codegen_1.operators.GT : codegen_1.operators.LT;
        cxt.fail$data((0, codegen_1._)`Object.keys(${data}).length ${op} ${schemaCode}`);
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/validation/required.js
var require_required = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/required.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var code_1 = require_code2();
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var error = {
      message: ({ params: { missingProperty } }) => (0, codegen_1.str)`must have required property '${missingProperty}'`,
      params: ({ params: { missingProperty } }) => (0, codegen_1._)`{missingProperty: ${missingProperty}}`
    };
    var def = {
      keyword: "required",
      type: "object",
      schemaType: "array",
      $data: true,
      error,
      code(cxt) {
        const { gen, schema: schema2, schemaCode, data, $data, it } = cxt;
        const { opts } = it;
        if (!$data && schema2.length === 0)
          return;
        const useLoop = schema2.length >= opts.loopRequired;
        if (it.allErrors)
          allErrorsMode();
        else
          exitOnErrorMode();
        if (opts.strictRequired) {
          const props = cxt.parentSchema.properties;
          const { definedProperties } = cxt.it;
          for (const requiredKey of schema2) {
            if ((props === null || props === void 0 ? void 0 : props[requiredKey]) === void 0 && !definedProperties.has(requiredKey)) {
              const schemaPath = it.schemaEnv.baseId + it.errSchemaPath;
              const msg = `required property "${requiredKey}" is not defined at "${schemaPath}" (strictRequired)`;
              (0, util_1.checkStrictMode)(it, msg, it.opts.strictRequired);
            }
          }
        }
        function allErrorsMode() {
          if (useLoop || $data) {
            cxt.block$data(codegen_1.nil, loopAllRequired);
          } else {
            for (const prop of schema2) {
              (0, code_1.checkReportMissingProp)(cxt, prop);
            }
          }
        }
        function exitOnErrorMode() {
          const missing = gen.let("missing");
          if (useLoop || $data) {
            const valid = gen.let("valid", true);
            cxt.block$data(valid, () => loopUntilMissing(missing, valid));
            cxt.ok(valid);
          } else {
            gen.if((0, code_1.checkMissingProp)(cxt, schema2, missing));
            (0, code_1.reportMissingProp)(cxt, missing);
            gen.else();
          }
        }
        function loopAllRequired() {
          gen.forOf("prop", schemaCode, (prop) => {
            cxt.setParams({ missingProperty: prop });
            gen.if((0, code_1.noPropertyInData)(gen, data, prop, opts.ownProperties), () => cxt.error());
          });
        }
        function loopUntilMissing(missing, valid) {
          cxt.setParams({ missingProperty: missing });
          gen.forOf(missing, schemaCode, () => {
            gen.assign(valid, (0, code_1.propertyInData)(gen, data, missing, opts.ownProperties));
            gen.if((0, codegen_1.not)(valid), () => {
              cxt.error();
              gen.break();
            });
          }, codegen_1.nil);
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/validation/limitItems.js
var require_limitItems = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/limitItems.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var error = {
      message({ keyword, schemaCode }) {
        const comp = keyword === "maxItems" ? "more" : "fewer";
        return (0, codegen_1.str)`must NOT have ${comp} than ${schemaCode} items`;
      },
      params: ({ schemaCode }) => (0, codegen_1._)`{limit: ${schemaCode}}`
    };
    var def = {
      keyword: ["maxItems", "minItems"],
      type: "array",
      schemaType: "number",
      $data: true,
      error,
      code(cxt) {
        const { keyword, data, schemaCode } = cxt;
        const op = keyword === "maxItems" ? codegen_1.operators.GT : codegen_1.operators.LT;
        cxt.fail$data((0, codegen_1._)`${data}.length ${op} ${schemaCode}`);
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/runtime/equal.js
var require_equal = __commonJS({
  "node_modules/ajv/dist/runtime/equal.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var equal = require_fast_deep_equal();
    equal.code = 'require("ajv/dist/runtime/equal").default';
    exports.default = equal;
  }
});

// node_modules/ajv/dist/vocabularies/validation/uniqueItems.js
var require_uniqueItems = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/uniqueItems.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var dataType_1 = require_dataType();
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var equal_1 = require_equal();
    var error = {
      message: ({ params: { i, j } }) => (0, codegen_1.str)`must NOT have duplicate items (items ## ${j} and ${i} are identical)`,
      params: ({ params: { i, j } }) => (0, codegen_1._)`{i: ${i}, j: ${j}}`
    };
    var def = {
      keyword: "uniqueItems",
      type: "array",
      schemaType: "boolean",
      $data: true,
      error,
      code(cxt) {
        const { gen, data, $data, schema: schema2, parentSchema, schemaCode, it } = cxt;
        if (!$data && !schema2)
          return;
        const valid = gen.let("valid");
        const itemTypes = parentSchema.items ? (0, dataType_1.getSchemaTypes)(parentSchema.items) : [];
        cxt.block$data(valid, validateUniqueItems, (0, codegen_1._)`${schemaCode} === false`);
        cxt.ok(valid);
        function validateUniqueItems() {
          const i = gen.let("i", (0, codegen_1._)`${data}.length`);
          const j = gen.let("j");
          cxt.setParams({ i, j });
          gen.assign(valid, true);
          gen.if((0, codegen_1._)`${i} > 1`, () => (canOptimize() ? loopN : loopN2)(i, j));
        }
        function canOptimize() {
          return itemTypes.length > 0 && !itemTypes.some((t) => t === "object" || t === "array");
        }
        function loopN(i, j) {
          const item = gen.name("item");
          const wrongType = (0, dataType_1.checkDataTypes)(itemTypes, item, it.opts.strictNumbers, dataType_1.DataType.Wrong);
          const indices = gen.const("indices", (0, codegen_1._)`{}`);
          gen.for((0, codegen_1._)`;${i}--;`, () => {
            gen.let(item, (0, codegen_1._)`${data}[${i}]`);
            gen.if(wrongType, (0, codegen_1._)`continue`);
            if (itemTypes.length > 1)
              gen.if((0, codegen_1._)`typeof ${item} == "string"`, (0, codegen_1._)`${item} += "_"`);
            gen.if((0, codegen_1._)`typeof ${indices}[${item}] == "number"`, () => {
              gen.assign(j, (0, codegen_1._)`${indices}[${item}]`);
              cxt.error();
              gen.assign(valid, false).break();
            }).code((0, codegen_1._)`${indices}[${item}] = ${i}`);
          });
        }
        function loopN2(i, j) {
          const eql = (0, util_1.useFunc)(gen, equal_1.default);
          const outer = gen.name("outer");
          gen.label(outer).for((0, codegen_1._)`;${i}--;`, () => gen.for((0, codegen_1._)`${j} = ${i}; ${j}--;`, () => gen.if((0, codegen_1._)`${eql}(${data}[${i}], ${data}[${j}])`, () => {
            cxt.error();
            gen.assign(valid, false).break(outer);
          })));
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/validation/const.js
var require_const = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/const.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var equal_1 = require_equal();
    var error = {
      message: "must be equal to constant",
      params: ({ schemaCode }) => (0, codegen_1._)`{allowedValue: ${schemaCode}}`
    };
    var def = {
      keyword: "const",
      $data: true,
      error,
      code(cxt) {
        const { gen, data, $data, schemaCode, schema: schema2 } = cxt;
        if ($data || schema2 && typeof schema2 == "object") {
          cxt.fail$data((0, codegen_1._)`!${(0, util_1.useFunc)(gen, equal_1.default)}(${data}, ${schemaCode})`);
        } else {
          cxt.fail((0, codegen_1._)`${schema2} !== ${data}`);
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/validation/enum.js
var require_enum = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/enum.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var equal_1 = require_equal();
    var error = {
      message: "must be equal to one of the allowed values",
      params: ({ schemaCode }) => (0, codegen_1._)`{allowedValues: ${schemaCode}}`
    };
    var def = {
      keyword: "enum",
      schemaType: "array",
      $data: true,
      error,
      code(cxt) {
        const { gen, data, $data, schema: schema2, schemaCode, it } = cxt;
        if (!$data && schema2.length === 0)
          throw new Error("enum must have non-empty array");
        const useLoop = schema2.length >= it.opts.loopEnum;
        let eql;
        const getEql = () => eql !== null && eql !== void 0 ? eql : eql = (0, util_1.useFunc)(gen, equal_1.default);
        let valid;
        if (useLoop || $data) {
          valid = gen.let("valid");
          cxt.block$data(valid, loopEnum);
        } else {
          if (!Array.isArray(schema2))
            throw new Error("ajv implementation error");
          const vSchema = gen.const("vSchema", schemaCode);
          valid = (0, codegen_1.or)(...schema2.map((_x, i) => equalCode(vSchema, i)));
        }
        cxt.pass(valid);
        function loopEnum() {
          gen.assign(valid, false);
          gen.forOf("v", schemaCode, (v) => gen.if((0, codegen_1._)`${getEql()}(${data}, ${v})`, () => gen.assign(valid, true).break()));
        }
        function equalCode(vSchema, i) {
          const sch = schema2[i];
          return typeof sch === "object" && sch !== null ? (0, codegen_1._)`${getEql()}(${data}, ${vSchema}[${i}])` : (0, codegen_1._)`${data} === ${sch}`;
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/validation/index.js
var require_validation = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var limitNumber_1 = require_limitNumber();
    var multipleOf_1 = require_multipleOf();
    var limitLength_1 = require_limitLength();
    var pattern_1 = require_pattern();
    var limitProperties_1 = require_limitProperties();
    var required_1 = require_required();
    var limitItems_1 = require_limitItems();
    var uniqueItems_1 = require_uniqueItems();
    var const_1 = require_const();
    var enum_1 = require_enum();
    var validation = [
      // number
      limitNumber_1.default,
      multipleOf_1.default,
      // string
      limitLength_1.default,
      pattern_1.default,
      // object
      limitProperties_1.default,
      required_1.default,
      // array
      limitItems_1.default,
      uniqueItems_1.default,
      // any
      { keyword: "type", schemaType: ["string", "array"] },
      { keyword: "nullable", schemaType: "boolean" },
      const_1.default,
      enum_1.default
    ];
    exports.default = validation;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/additionalItems.js
var require_additionalItems = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/additionalItems.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.validateAdditionalItems = void 0;
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var error = {
      message: ({ params: { len } }) => (0, codegen_1.str)`must NOT have more than ${len} items`,
      params: ({ params: { len } }) => (0, codegen_1._)`{limit: ${len}}`
    };
    var def = {
      keyword: "additionalItems",
      type: "array",
      schemaType: ["boolean", "object"],
      before: "uniqueItems",
      error,
      code(cxt) {
        const { parentSchema, it } = cxt;
        const { items } = parentSchema;
        if (!Array.isArray(items)) {
          (0, util_1.checkStrictMode)(it, '"additionalItems" is ignored when "items" is not an array of schemas');
          return;
        }
        validateAdditionalItems(cxt, items);
      }
    };
    function validateAdditionalItems(cxt, items) {
      const { gen, schema: schema2, data, keyword, it } = cxt;
      it.items = true;
      const len = gen.const("len", (0, codegen_1._)`${data}.length`);
      if (schema2 === false) {
        cxt.setParams({ len: items.length });
        cxt.pass((0, codegen_1._)`${len} <= ${items.length}`);
      } else if (typeof schema2 == "object" && !(0, util_1.alwaysValidSchema)(it, schema2)) {
        const valid = gen.var("valid", (0, codegen_1._)`${len} <= ${items.length}`);
        gen.if((0, codegen_1.not)(valid), () => validateItems(valid));
        cxt.ok(valid);
      }
      function validateItems(valid) {
        gen.forRange("i", items.length, len, (i) => {
          cxt.subschema({ keyword, dataProp: i, dataPropType: util_1.Type.Num }, valid);
          if (!it.allErrors)
            gen.if((0, codegen_1.not)(valid), () => gen.break());
        });
      }
    }
    exports.validateAdditionalItems = validateAdditionalItems;
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/items.js
var require_items = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/items.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.validateTuple = void 0;
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var code_1 = require_code2();
    var def = {
      keyword: "items",
      type: "array",
      schemaType: ["object", "array", "boolean"],
      before: "uniqueItems",
      code(cxt) {
        const { schema: schema2, it } = cxt;
        if (Array.isArray(schema2))
          return validateTuple(cxt, "additionalItems", schema2);
        it.items = true;
        if ((0, util_1.alwaysValidSchema)(it, schema2))
          return;
        cxt.ok((0, code_1.validateArray)(cxt));
      }
    };
    function validateTuple(cxt, extraItems, schArr = cxt.schema) {
      const { gen, parentSchema, data, keyword, it } = cxt;
      checkStrictTuple(parentSchema);
      if (it.opts.unevaluated && schArr.length && it.items !== true) {
        it.items = util_1.mergeEvaluated.items(gen, schArr.length, it.items);
      }
      const valid = gen.name("valid");
      const len = gen.const("len", (0, codegen_1._)`${data}.length`);
      schArr.forEach((sch, i) => {
        if ((0, util_1.alwaysValidSchema)(it, sch))
          return;
        gen.if((0, codegen_1._)`${len} > ${i}`, () => cxt.subschema({
          keyword,
          schemaProp: i,
          dataProp: i
        }, valid));
        cxt.ok(valid);
      });
      function checkStrictTuple(sch) {
        const { opts, errSchemaPath } = it;
        const l = schArr.length;
        const fullTuple = l === sch.minItems && (l === sch.maxItems || sch[extraItems] === false);
        if (opts.strictTuples && !fullTuple) {
          const msg = `"${keyword}" is ${l}-tuple, but minItems or maxItems/${extraItems} are not specified or different at path "${errSchemaPath}"`;
          (0, util_1.checkStrictMode)(it, msg, opts.strictTuples);
        }
      }
    }
    exports.validateTuple = validateTuple;
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/prefixItems.js
var require_prefixItems = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/prefixItems.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var items_1 = require_items();
    var def = {
      keyword: "prefixItems",
      type: "array",
      schemaType: ["array"],
      before: "uniqueItems",
      code: (cxt) => (0, items_1.validateTuple)(cxt, "items")
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/items2020.js
var require_items2020 = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/items2020.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var code_1 = require_code2();
    var additionalItems_1 = require_additionalItems();
    var error = {
      message: ({ params: { len } }) => (0, codegen_1.str)`must NOT have more than ${len} items`,
      params: ({ params: { len } }) => (0, codegen_1._)`{limit: ${len}}`
    };
    var def = {
      keyword: "items",
      type: "array",
      schemaType: ["object", "boolean"],
      before: "uniqueItems",
      error,
      code(cxt) {
        const { schema: schema2, parentSchema, it } = cxt;
        const { prefixItems } = parentSchema;
        it.items = true;
        if ((0, util_1.alwaysValidSchema)(it, schema2))
          return;
        if (prefixItems)
          (0, additionalItems_1.validateAdditionalItems)(cxt, prefixItems);
        else
          cxt.ok((0, code_1.validateArray)(cxt));
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/contains.js
var require_contains = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/contains.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var error = {
      message: ({ params: { min, max } }) => max === void 0 ? (0, codegen_1.str)`must contain at least ${min} valid item(s)` : (0, codegen_1.str)`must contain at least ${min} and no more than ${max} valid item(s)`,
      params: ({ params: { min, max } }) => max === void 0 ? (0, codegen_1._)`{minContains: ${min}}` : (0, codegen_1._)`{minContains: ${min}, maxContains: ${max}}`
    };
    var def = {
      keyword: "contains",
      type: "array",
      schemaType: ["object", "boolean"],
      before: "uniqueItems",
      trackErrors: true,
      error,
      code(cxt) {
        const { gen, schema: schema2, parentSchema, data, it } = cxt;
        let min;
        let max;
        const { minContains, maxContains } = parentSchema;
        if (it.opts.next) {
          min = minContains === void 0 ? 1 : minContains;
          max = maxContains;
        } else {
          min = 1;
        }
        const len = gen.const("len", (0, codegen_1._)`${data}.length`);
        cxt.setParams({ min, max });
        if (max === void 0 && min === 0) {
          (0, util_1.checkStrictMode)(it, `"minContains" == 0 without "maxContains": "contains" keyword ignored`);
          return;
        }
        if (max !== void 0 && min > max) {
          (0, util_1.checkStrictMode)(it, `"minContains" > "maxContains" is always invalid`);
          cxt.fail();
          return;
        }
        if ((0, util_1.alwaysValidSchema)(it, schema2)) {
          let cond = (0, codegen_1._)`${len} >= ${min}`;
          if (max !== void 0)
            cond = (0, codegen_1._)`${cond} && ${len} <= ${max}`;
          cxt.pass(cond);
          return;
        }
        it.items = true;
        const valid = gen.name("valid");
        if (max === void 0 && min === 1) {
          validateItems(valid, () => gen.if(valid, () => gen.break()));
        } else if (min === 0) {
          gen.let(valid, true);
          if (max !== void 0)
            gen.if((0, codegen_1._)`${data}.length > 0`, validateItemsWithCount);
        } else {
          gen.let(valid, false);
          validateItemsWithCount();
        }
        cxt.result(valid, () => cxt.reset());
        function validateItemsWithCount() {
          const schValid = gen.name("_valid");
          const count = gen.let("count", 0);
          validateItems(schValid, () => gen.if(schValid, () => checkLimits(count)));
        }
        function validateItems(_valid, block) {
          gen.forRange("i", 0, len, (i) => {
            cxt.subschema({
              keyword: "contains",
              dataProp: i,
              dataPropType: util_1.Type.Num,
              compositeRule: true
            }, _valid);
            block();
          });
        }
        function checkLimits(count) {
          gen.code((0, codegen_1._)`${count}++`);
          if (max === void 0) {
            gen.if((0, codegen_1._)`${count} >= ${min}`, () => gen.assign(valid, true).break());
          } else {
            gen.if((0, codegen_1._)`${count} > ${max}`, () => gen.assign(valid, false).break());
            if (min === 1)
              gen.assign(valid, true);
            else
              gen.if((0, codegen_1._)`${count} >= ${min}`, () => gen.assign(valid, true));
          }
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/dependencies.js
var require_dependencies = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/dependencies.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.validateSchemaDeps = exports.validatePropertyDeps = exports.error = void 0;
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var code_1 = require_code2();
    exports.error = {
      message: ({ params: { property, depsCount, deps } }) => {
        const property_ies = depsCount === 1 ? "property" : "properties";
        return (0, codegen_1.str)`must have ${property_ies} ${deps} when property ${property} is present`;
      },
      params: ({ params: { property, depsCount, deps, missingProperty } }) => (0, codegen_1._)`{property: ${property},
    missingProperty: ${missingProperty},
    depsCount: ${depsCount},
    deps: ${deps}}`
      // TODO change to reference
    };
    var def = {
      keyword: "dependencies",
      type: "object",
      schemaType: "object",
      error: exports.error,
      code(cxt) {
        const [propDeps, schDeps] = splitDependencies(cxt);
        validatePropertyDeps(cxt, propDeps);
        validateSchemaDeps(cxt, schDeps);
      }
    };
    function splitDependencies({ schema: schema2 }) {
      const propertyDeps = {};
      const schemaDeps = {};
      for (const key in schema2) {
        if (key === "__proto__")
          continue;
        const deps = Array.isArray(schema2[key]) ? propertyDeps : schemaDeps;
        deps[key] = schema2[key];
      }
      return [propertyDeps, schemaDeps];
    }
    function validatePropertyDeps(cxt, propertyDeps = cxt.schema) {
      const { gen, data, it } = cxt;
      if (Object.keys(propertyDeps).length === 0)
        return;
      const missing = gen.let("missing");
      for (const prop in propertyDeps) {
        const deps = propertyDeps[prop];
        if (deps.length === 0)
          continue;
        const hasProperty = (0, code_1.propertyInData)(gen, data, prop, it.opts.ownProperties);
        cxt.setParams({
          property: prop,
          depsCount: deps.length,
          deps: deps.join(", ")
        });
        if (it.allErrors) {
          gen.if(hasProperty, () => {
            for (const depProp of deps) {
              (0, code_1.checkReportMissingProp)(cxt, depProp);
            }
          });
        } else {
          gen.if((0, codegen_1._)`${hasProperty} && (${(0, code_1.checkMissingProp)(cxt, deps, missing)})`);
          (0, code_1.reportMissingProp)(cxt, missing);
          gen.else();
        }
      }
    }
    exports.validatePropertyDeps = validatePropertyDeps;
    function validateSchemaDeps(cxt, schemaDeps = cxt.schema) {
      const { gen, data, keyword, it } = cxt;
      const valid = gen.name("valid");
      for (const prop in schemaDeps) {
        if ((0, util_1.alwaysValidSchema)(it, schemaDeps[prop]))
          continue;
        gen.if(
          (0, code_1.propertyInData)(gen, data, prop, it.opts.ownProperties),
          () => {
            const schCxt = cxt.subschema({ keyword, schemaProp: prop }, valid);
            cxt.mergeValidEvaluated(schCxt, valid);
          },
          () => gen.var(valid, true)
          // TODO var
        );
        cxt.ok(valid);
      }
    }
    exports.validateSchemaDeps = validateSchemaDeps;
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/propertyNames.js
var require_propertyNames = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/propertyNames.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var error = {
      message: "property name must be valid",
      params: ({ params }) => (0, codegen_1._)`{propertyName: ${params.propertyName}}`
    };
    var def = {
      keyword: "propertyNames",
      type: "object",
      schemaType: ["object", "boolean"],
      error,
      code(cxt) {
        const { gen, schema: schema2, data, it } = cxt;
        if ((0, util_1.alwaysValidSchema)(it, schema2))
          return;
        const valid = gen.name("valid");
        gen.forIn("key", data, (key) => {
          cxt.setParams({ propertyName: key });
          cxt.subschema({
            keyword: "propertyNames",
            data: key,
            dataTypes: ["string"],
            propertyName: key,
            compositeRule: true
          }, valid);
          gen.if((0, codegen_1.not)(valid), () => {
            cxt.error(true);
            if (!it.allErrors)
              gen.break();
          });
        });
        cxt.ok(valid);
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/additionalProperties.js
var require_additionalProperties = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/additionalProperties.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var code_1 = require_code2();
    var codegen_1 = require_codegen();
    var names_1 = require_names();
    var util_1 = require_util();
    var error = {
      message: "must NOT have additional properties",
      params: ({ params }) => (0, codegen_1._)`{additionalProperty: ${params.additionalProperty}}`
    };
    var def = {
      keyword: "additionalProperties",
      type: ["object"],
      schemaType: ["boolean", "object"],
      allowUndefined: true,
      trackErrors: true,
      error,
      code(cxt) {
        const { gen, schema: schema2, parentSchema, data, errsCount, it } = cxt;
        if (!errsCount)
          throw new Error("ajv implementation error");
        const { allErrors, opts } = it;
        it.props = true;
        if (opts.removeAdditional !== "all" && (0, util_1.alwaysValidSchema)(it, schema2))
          return;
        const props = (0, code_1.allSchemaProperties)(parentSchema.properties);
        const patProps = (0, code_1.allSchemaProperties)(parentSchema.patternProperties);
        checkAdditionalProperties();
        cxt.ok((0, codegen_1._)`${errsCount} === ${names_1.default.errors}`);
        function checkAdditionalProperties() {
          gen.forIn("key", data, (key) => {
            if (!props.length && !patProps.length)
              additionalPropertyCode(key);
            else
              gen.if(isAdditional(key), () => additionalPropertyCode(key));
          });
        }
        function isAdditional(key) {
          let definedProp;
          if (props.length > 8) {
            const propsSchema = (0, util_1.schemaRefOrVal)(it, parentSchema.properties, "properties");
            definedProp = (0, code_1.isOwnProperty)(gen, propsSchema, key);
          } else if (props.length) {
            definedProp = (0, codegen_1.or)(...props.map((p) => (0, codegen_1._)`${key} === ${p}`));
          } else {
            definedProp = codegen_1.nil;
          }
          if (patProps.length) {
            definedProp = (0, codegen_1.or)(definedProp, ...patProps.map((p) => (0, codegen_1._)`${(0, code_1.usePattern)(cxt, p)}.test(${key})`));
          }
          return (0, codegen_1.not)(definedProp);
        }
        function deleteAdditional(key) {
          gen.code((0, codegen_1._)`delete ${data}[${key}]`);
        }
        function additionalPropertyCode(key) {
          if (opts.removeAdditional === "all" || opts.removeAdditional && schema2 === false) {
            deleteAdditional(key);
            return;
          }
          if (schema2 === false) {
            cxt.setParams({ additionalProperty: key });
            cxt.error();
            if (!allErrors)
              gen.break();
            return;
          }
          if (typeof schema2 == "object" && !(0, util_1.alwaysValidSchema)(it, schema2)) {
            const valid = gen.name("valid");
            if (opts.removeAdditional === "failing") {
              applyAdditionalSchema(key, valid, false);
              gen.if((0, codegen_1.not)(valid), () => {
                cxt.reset();
                deleteAdditional(key);
              });
            } else {
              applyAdditionalSchema(key, valid);
              if (!allErrors)
                gen.if((0, codegen_1.not)(valid), () => gen.break());
            }
          }
        }
        function applyAdditionalSchema(key, valid, errors2) {
          const subschema = {
            keyword: "additionalProperties",
            dataProp: key,
            dataPropType: util_1.Type.Str
          };
          if (errors2 === false) {
            Object.assign(subschema, {
              compositeRule: true,
              createErrors: false,
              allErrors: false
            });
          }
          cxt.subschema(subschema, valid);
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/properties.js
var require_properties = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/properties.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var validate_1 = require_validate();
    var code_1 = require_code2();
    var util_1 = require_util();
    var additionalProperties_1 = require_additionalProperties();
    var def = {
      keyword: "properties",
      type: "object",
      schemaType: "object",
      code(cxt) {
        const { gen, schema: schema2, parentSchema, data, it } = cxt;
        if (it.opts.removeAdditional === "all" && parentSchema.additionalProperties === void 0) {
          additionalProperties_1.default.code(new validate_1.KeywordCxt(it, additionalProperties_1.default, "additionalProperties"));
        }
        const allProps = (0, code_1.allSchemaProperties)(schema2);
        for (const prop of allProps) {
          it.definedProperties.add(prop);
        }
        if (it.opts.unevaluated && allProps.length && it.props !== true) {
          it.props = util_1.mergeEvaluated.props(gen, (0, util_1.toHash)(allProps), it.props);
        }
        const properties = allProps.filter((p) => !(0, util_1.alwaysValidSchema)(it, schema2[p]));
        if (properties.length === 0)
          return;
        const valid = gen.name("valid");
        for (const prop of properties) {
          if (hasDefault(prop)) {
            applyPropertySchema(prop);
          } else {
            gen.if((0, code_1.propertyInData)(gen, data, prop, it.opts.ownProperties));
            applyPropertySchema(prop);
            if (!it.allErrors)
              gen.else().var(valid, true);
            gen.endIf();
          }
          cxt.it.definedProperties.add(prop);
          cxt.ok(valid);
        }
        function hasDefault(prop) {
          return it.opts.useDefaults && !it.compositeRule && schema2[prop].default !== void 0;
        }
        function applyPropertySchema(prop) {
          cxt.subschema({
            keyword: "properties",
            schemaProp: prop,
            dataProp: prop
          }, valid);
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/patternProperties.js
var require_patternProperties = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/patternProperties.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var code_1 = require_code2();
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var util_2 = require_util();
    var def = {
      keyword: "patternProperties",
      type: "object",
      schemaType: "object",
      code(cxt) {
        const { gen, schema: schema2, data, parentSchema, it } = cxt;
        const { opts } = it;
        const patterns = (0, code_1.allSchemaProperties)(schema2);
        const alwaysValidPatterns = patterns.filter((p) => (0, util_1.alwaysValidSchema)(it, schema2[p]));
        if (patterns.length === 0 || alwaysValidPatterns.length === patterns.length && (!it.opts.unevaluated || it.props === true)) {
          return;
        }
        const checkProperties = opts.strictSchema && !opts.allowMatchingProperties && parentSchema.properties;
        const valid = gen.name("valid");
        if (it.props !== true && !(it.props instanceof codegen_1.Name)) {
          it.props = (0, util_2.evaluatedPropsToName)(gen, it.props);
        }
        const { props } = it;
        validatePatternProperties();
        function validatePatternProperties() {
          for (const pat of patterns) {
            if (checkProperties)
              checkMatchingProperties(pat);
            if (it.allErrors) {
              validateProperties(pat);
            } else {
              gen.var(valid, true);
              validateProperties(pat);
              gen.if(valid);
            }
          }
        }
        function checkMatchingProperties(pat) {
          for (const prop in checkProperties) {
            if (new RegExp(pat).test(prop)) {
              (0, util_1.checkStrictMode)(it, `property ${prop} matches pattern ${pat} (use allowMatchingProperties)`);
            }
          }
        }
        function validateProperties(pat) {
          gen.forIn("key", data, (key) => {
            gen.if((0, codegen_1._)`${(0, code_1.usePattern)(cxt, pat)}.test(${key})`, () => {
              const alwaysValid = alwaysValidPatterns.includes(pat);
              if (!alwaysValid) {
                cxt.subschema({
                  keyword: "patternProperties",
                  schemaProp: pat,
                  dataProp: key,
                  dataPropType: util_2.Type.Str
                }, valid);
              }
              if (it.opts.unevaluated && props !== true) {
                gen.assign((0, codegen_1._)`${props}[${key}]`, true);
              } else if (!alwaysValid && !it.allErrors) {
                gen.if((0, codegen_1.not)(valid), () => gen.break());
              }
            });
          });
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/not.js
var require_not = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/not.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var util_1 = require_util();
    var def = {
      keyword: "not",
      schemaType: ["object", "boolean"],
      trackErrors: true,
      code(cxt) {
        const { gen, schema: schema2, it } = cxt;
        if ((0, util_1.alwaysValidSchema)(it, schema2)) {
          cxt.fail();
          return;
        }
        const valid = gen.name("valid");
        cxt.subschema({
          keyword: "not",
          compositeRule: true,
          createErrors: false,
          allErrors: false
        }, valid);
        cxt.failResult(valid, () => cxt.reset(), () => cxt.error());
      },
      error: { message: "must NOT be valid" }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/anyOf.js
var require_anyOf = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/anyOf.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var code_1 = require_code2();
    var def = {
      keyword: "anyOf",
      schemaType: "array",
      trackErrors: true,
      code: code_1.validateUnion,
      error: { message: "must match a schema in anyOf" }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/oneOf.js
var require_oneOf = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/oneOf.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var error = {
      message: "must match exactly one schema in oneOf",
      params: ({ params }) => (0, codegen_1._)`{passingSchemas: ${params.passing}}`
    };
    var def = {
      keyword: "oneOf",
      schemaType: "array",
      trackErrors: true,
      error,
      code(cxt) {
        const { gen, schema: schema2, parentSchema, it } = cxt;
        if (!Array.isArray(schema2))
          throw new Error("ajv implementation error");
        if (it.opts.discriminator && parentSchema.discriminator)
          return;
        const schArr = schema2;
        const valid = gen.let("valid", false);
        const passing = gen.let("passing", null);
        const schValid = gen.name("_valid");
        cxt.setParams({ passing });
        gen.block(validateOneOf);
        cxt.result(valid, () => cxt.reset(), () => cxt.error(true));
        function validateOneOf() {
          schArr.forEach((sch, i) => {
            let schCxt;
            if ((0, util_1.alwaysValidSchema)(it, sch)) {
              gen.var(schValid, true);
            } else {
              schCxt = cxt.subschema({
                keyword: "oneOf",
                schemaProp: i,
                compositeRule: true
              }, schValid);
            }
            if (i > 0) {
              gen.if((0, codegen_1._)`${schValid} && ${valid}`).assign(valid, false).assign(passing, (0, codegen_1._)`[${passing}, ${i}]`).else();
            }
            gen.if(schValid, () => {
              gen.assign(valid, true);
              gen.assign(passing, i);
              if (schCxt)
                cxt.mergeEvaluated(schCxt, codegen_1.Name);
            });
          });
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/allOf.js
var require_allOf = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/allOf.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var util_1 = require_util();
    var def = {
      keyword: "allOf",
      schemaType: "array",
      code(cxt) {
        const { gen, schema: schema2, it } = cxt;
        if (!Array.isArray(schema2))
          throw new Error("ajv implementation error");
        const valid = gen.name("valid");
        schema2.forEach((sch, i) => {
          if ((0, util_1.alwaysValidSchema)(it, sch))
            return;
          const schCxt = cxt.subschema({ keyword: "allOf", schemaProp: i }, valid);
          cxt.ok(valid);
          cxt.mergeEvaluated(schCxt);
        });
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/if.js
var require_if = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/if.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var error = {
      message: ({ params }) => (0, codegen_1.str)`must match "${params.ifClause}" schema`,
      params: ({ params }) => (0, codegen_1._)`{failingKeyword: ${params.ifClause}}`
    };
    var def = {
      keyword: "if",
      schemaType: ["object", "boolean"],
      trackErrors: true,
      error,
      code(cxt) {
        const { gen, parentSchema, it } = cxt;
        if (parentSchema.then === void 0 && parentSchema.else === void 0) {
          (0, util_1.checkStrictMode)(it, '"if" without "then" and "else" is ignored');
        }
        const hasThen = hasSchema(it, "then");
        const hasElse = hasSchema(it, "else");
        if (!hasThen && !hasElse)
          return;
        const valid = gen.let("valid", true);
        const schValid = gen.name("_valid");
        validateIf();
        cxt.reset();
        if (hasThen && hasElse) {
          const ifClause = gen.let("ifClause");
          cxt.setParams({ ifClause });
          gen.if(schValid, validateClause("then", ifClause), validateClause("else", ifClause));
        } else if (hasThen) {
          gen.if(schValid, validateClause("then"));
        } else {
          gen.if((0, codegen_1.not)(schValid), validateClause("else"));
        }
        cxt.pass(valid, () => cxt.error(true));
        function validateIf() {
          const schCxt = cxt.subschema({
            keyword: "if",
            compositeRule: true,
            createErrors: false,
            allErrors: false
          }, schValid);
          cxt.mergeEvaluated(schCxt);
        }
        function validateClause(keyword, ifClause) {
          return () => {
            const schCxt = cxt.subschema({ keyword }, schValid);
            gen.assign(valid, schValid);
            cxt.mergeValidEvaluated(schCxt, valid);
            if (ifClause)
              gen.assign(ifClause, (0, codegen_1._)`${keyword}`);
            else
              cxt.setParams({ ifClause: keyword });
          };
        }
      }
    };
    function hasSchema(it, keyword) {
      const schema2 = it.schema[keyword];
      return schema2 !== void 0 && !(0, util_1.alwaysValidSchema)(it, schema2);
    }
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/thenElse.js
var require_thenElse = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/thenElse.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var util_1 = require_util();
    var def = {
      keyword: ["then", "else"],
      schemaType: ["object", "boolean"],
      code({ keyword, parentSchema, it }) {
        if (parentSchema.if === void 0)
          (0, util_1.checkStrictMode)(it, `"${keyword}" without "if" is ignored`);
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/index.js
var require_applicator = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var additionalItems_1 = require_additionalItems();
    var prefixItems_1 = require_prefixItems();
    var items_1 = require_items();
    var items2020_1 = require_items2020();
    var contains_1 = require_contains();
    var dependencies_1 = require_dependencies();
    var propertyNames_1 = require_propertyNames();
    var additionalProperties_1 = require_additionalProperties();
    var properties_1 = require_properties();
    var patternProperties_1 = require_patternProperties();
    var not_1 = require_not();
    var anyOf_1 = require_anyOf();
    var oneOf_1 = require_oneOf();
    var allOf_1 = require_allOf();
    var if_1 = require_if();
    var thenElse_1 = require_thenElse();
    function getApplicator(draft2020 = false) {
      const applicator = [
        // any
        not_1.default,
        anyOf_1.default,
        oneOf_1.default,
        allOf_1.default,
        if_1.default,
        thenElse_1.default,
        // object
        propertyNames_1.default,
        additionalProperties_1.default,
        dependencies_1.default,
        properties_1.default,
        patternProperties_1.default
      ];
      if (draft2020)
        applicator.push(prefixItems_1.default, items2020_1.default);
      else
        applicator.push(additionalItems_1.default, items_1.default);
      applicator.push(contains_1.default);
      return applicator;
    }
    exports.default = getApplicator;
  }
});

// node_modules/ajv/dist/vocabularies/dynamic/dynamicAnchor.js
var require_dynamicAnchor = __commonJS({
  "node_modules/ajv/dist/vocabularies/dynamic/dynamicAnchor.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.dynamicAnchor = void 0;
    var codegen_1 = require_codegen();
    var names_1 = require_names();
    var compile_1 = require_compile();
    var ref_1 = require_ref();
    var def = {
      keyword: "$dynamicAnchor",
      schemaType: "string",
      code: (cxt) => dynamicAnchor(cxt, cxt.schema)
    };
    function dynamicAnchor(cxt, anchor) {
      const { gen, it } = cxt;
      it.schemaEnv.root.dynamicAnchors[anchor] = true;
      const v = (0, codegen_1._)`${names_1.default.dynamicAnchors}${(0, codegen_1.getProperty)(anchor)}`;
      const validate = it.errSchemaPath === "#" ? it.validateName : _getValidate(cxt);
      gen.if((0, codegen_1._)`!${v}`, () => gen.assign(v, validate));
    }
    exports.dynamicAnchor = dynamicAnchor;
    function _getValidate(cxt) {
      const { schemaEnv, schema: schema2, self } = cxt.it;
      const { root, baseId, localRefs, meta } = schemaEnv.root;
      const { schemaId } = self.opts;
      const sch = new compile_1.SchemaEnv({ schema: schema2, schemaId, root, baseId, localRefs, meta });
      compile_1.compileSchema.call(self, sch);
      return (0, ref_1.getValidate)(cxt, sch);
    }
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/dynamic/dynamicRef.js
var require_dynamicRef = __commonJS({
  "node_modules/ajv/dist/vocabularies/dynamic/dynamicRef.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.dynamicRef = void 0;
    var codegen_1 = require_codegen();
    var names_1 = require_names();
    var ref_1 = require_ref();
    var def = {
      keyword: "$dynamicRef",
      schemaType: "string",
      code: (cxt) => dynamicRef(cxt, cxt.schema)
    };
    function dynamicRef(cxt, ref) {
      const { gen, keyword, it } = cxt;
      if (ref[0] !== "#")
        throw new Error(`"${keyword}" only supports hash fragment reference`);
      const anchor = ref.slice(1);
      if (it.allErrors) {
        _dynamicRef();
      } else {
        const valid = gen.let("valid", false);
        _dynamicRef(valid);
        cxt.ok(valid);
      }
      function _dynamicRef(valid) {
        if (it.schemaEnv.root.dynamicAnchors[anchor]) {
          const v = gen.let("_v", (0, codegen_1._)`${names_1.default.dynamicAnchors}${(0, codegen_1.getProperty)(anchor)}`);
          gen.if(v, _callRef(v, valid), _callRef(it.validateName, valid));
        } else {
          _callRef(it.validateName, valid)();
        }
      }
      function _callRef(validate, valid) {
        return valid ? () => gen.block(() => {
          (0, ref_1.callRef)(cxt, validate);
          gen.let(valid, true);
        }) : () => (0, ref_1.callRef)(cxt, validate);
      }
    }
    exports.dynamicRef = dynamicRef;
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/dynamic/recursiveAnchor.js
var require_recursiveAnchor = __commonJS({
  "node_modules/ajv/dist/vocabularies/dynamic/recursiveAnchor.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var dynamicAnchor_1 = require_dynamicAnchor();
    var util_1 = require_util();
    var def = {
      keyword: "$recursiveAnchor",
      schemaType: "boolean",
      code(cxt) {
        if (cxt.schema)
          (0, dynamicAnchor_1.dynamicAnchor)(cxt, "");
        else
          (0, util_1.checkStrictMode)(cxt.it, "$recursiveAnchor: false is ignored");
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/dynamic/recursiveRef.js
var require_recursiveRef = __commonJS({
  "node_modules/ajv/dist/vocabularies/dynamic/recursiveRef.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var dynamicRef_1 = require_dynamicRef();
    var def = {
      keyword: "$recursiveRef",
      schemaType: "string",
      code: (cxt) => (0, dynamicRef_1.dynamicRef)(cxt, cxt.schema)
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/dynamic/index.js
var require_dynamic = __commonJS({
  "node_modules/ajv/dist/vocabularies/dynamic/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var dynamicAnchor_1 = require_dynamicAnchor();
    var dynamicRef_1 = require_dynamicRef();
    var recursiveAnchor_1 = require_recursiveAnchor();
    var recursiveRef_1 = require_recursiveRef();
    var dynamic = [dynamicAnchor_1.default, dynamicRef_1.default, recursiveAnchor_1.default, recursiveRef_1.default];
    exports.default = dynamic;
  }
});

// node_modules/ajv/dist/vocabularies/validation/dependentRequired.js
var require_dependentRequired = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/dependentRequired.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var dependencies_1 = require_dependencies();
    var def = {
      keyword: "dependentRequired",
      type: "object",
      schemaType: "object",
      error: dependencies_1.error,
      code: (cxt) => (0, dependencies_1.validatePropertyDeps)(cxt)
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/applicator/dependentSchemas.js
var require_dependentSchemas = __commonJS({
  "node_modules/ajv/dist/vocabularies/applicator/dependentSchemas.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var dependencies_1 = require_dependencies();
    var def = {
      keyword: "dependentSchemas",
      type: "object",
      schemaType: "object",
      code: (cxt) => (0, dependencies_1.validateSchemaDeps)(cxt)
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/validation/limitContains.js
var require_limitContains = __commonJS({
  "node_modules/ajv/dist/vocabularies/validation/limitContains.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var util_1 = require_util();
    var def = {
      keyword: ["maxContains", "minContains"],
      type: "array",
      schemaType: "number",
      code({ keyword, parentSchema, it }) {
        if (parentSchema.contains === void 0) {
          (0, util_1.checkStrictMode)(it, `"${keyword}" without "contains" is ignored`);
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/next.js
var require_next = __commonJS({
  "node_modules/ajv/dist/vocabularies/next.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var dependentRequired_1 = require_dependentRequired();
    var dependentSchemas_1 = require_dependentSchemas();
    var limitContains_1 = require_limitContains();
    var next = [dependentRequired_1.default, dependentSchemas_1.default, limitContains_1.default];
    exports.default = next;
  }
});

// node_modules/ajv/dist/vocabularies/unevaluated/unevaluatedProperties.js
var require_unevaluatedProperties = __commonJS({
  "node_modules/ajv/dist/vocabularies/unevaluated/unevaluatedProperties.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var names_1 = require_names();
    var error = {
      message: "must NOT have unevaluated properties",
      params: ({ params }) => (0, codegen_1._)`{unevaluatedProperty: ${params.unevaluatedProperty}}`
    };
    var def = {
      keyword: "unevaluatedProperties",
      type: "object",
      schemaType: ["boolean", "object"],
      trackErrors: true,
      error,
      code(cxt) {
        const { gen, schema: schema2, data, errsCount, it } = cxt;
        if (!errsCount)
          throw new Error("ajv implementation error");
        const { allErrors, props } = it;
        if (props instanceof codegen_1.Name) {
          gen.if((0, codegen_1._)`${props} !== true`, () => gen.forIn("key", data, (key) => gen.if(unevaluatedDynamic(props, key), () => unevaluatedPropCode(key))));
        } else if (props !== true) {
          gen.forIn("key", data, (key) => props === void 0 ? unevaluatedPropCode(key) : gen.if(unevaluatedStatic(props, key), () => unevaluatedPropCode(key)));
        }
        it.props = true;
        cxt.ok((0, codegen_1._)`${errsCount} === ${names_1.default.errors}`);
        function unevaluatedPropCode(key) {
          if (schema2 === false) {
            cxt.setParams({ unevaluatedProperty: key });
            cxt.error();
            if (!allErrors)
              gen.break();
            return;
          }
          if (!(0, util_1.alwaysValidSchema)(it, schema2)) {
            const valid = gen.name("valid");
            cxt.subschema({
              keyword: "unevaluatedProperties",
              dataProp: key,
              dataPropType: util_1.Type.Str
            }, valid);
            if (!allErrors)
              gen.if((0, codegen_1.not)(valid), () => gen.break());
          }
        }
        function unevaluatedDynamic(evaluatedProps, key) {
          return (0, codegen_1._)`!${evaluatedProps} || !${evaluatedProps}[${key}]`;
        }
        function unevaluatedStatic(evaluatedProps, key) {
          const ps = [];
          for (const p in evaluatedProps) {
            if (evaluatedProps[p] === true)
              ps.push((0, codegen_1._)`${key} !== ${p}`);
          }
          return (0, codegen_1.and)(...ps);
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/unevaluated/unevaluatedItems.js
var require_unevaluatedItems = __commonJS({
  "node_modules/ajv/dist/vocabularies/unevaluated/unevaluatedItems.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var util_1 = require_util();
    var error = {
      message: ({ params: { len } }) => (0, codegen_1.str)`must NOT have more than ${len} items`,
      params: ({ params: { len } }) => (0, codegen_1._)`{limit: ${len}}`
    };
    var def = {
      keyword: "unevaluatedItems",
      type: "array",
      schemaType: ["boolean", "object"],
      error,
      code(cxt) {
        const { gen, schema: schema2, data, it } = cxt;
        const items = it.items || 0;
        if (items === true)
          return;
        const len = gen.const("len", (0, codegen_1._)`${data}.length`);
        if (schema2 === false) {
          cxt.setParams({ len: items });
          cxt.fail((0, codegen_1._)`${len} > ${items}`);
        } else if (typeof schema2 == "object" && !(0, util_1.alwaysValidSchema)(it, schema2)) {
          const valid = gen.var("valid", (0, codegen_1._)`${len} <= ${items}`);
          gen.if((0, codegen_1.not)(valid), () => validateItems(valid, items));
          cxt.ok(valid);
        }
        it.items = true;
        function validateItems(valid, from) {
          gen.forRange("i", from, len, (i) => {
            cxt.subschema({ keyword: "unevaluatedItems", dataProp: i, dataPropType: util_1.Type.Num }, valid);
            if (!it.allErrors)
              gen.if((0, codegen_1.not)(valid), () => gen.break());
          });
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/unevaluated/index.js
var require_unevaluated = __commonJS({
  "node_modules/ajv/dist/vocabularies/unevaluated/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var unevaluatedProperties_1 = require_unevaluatedProperties();
    var unevaluatedItems_1 = require_unevaluatedItems();
    var unevaluated = [unevaluatedProperties_1.default, unevaluatedItems_1.default];
    exports.default = unevaluated;
  }
});

// node_modules/ajv/dist/vocabularies/format/format.js
var require_format = __commonJS({
  "node_modules/ajv/dist/vocabularies/format/format.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var error = {
      message: ({ schemaCode }) => (0, codegen_1.str)`must match format "${schemaCode}"`,
      params: ({ schemaCode }) => (0, codegen_1._)`{format: ${schemaCode}}`
    };
    var def = {
      keyword: "format",
      type: ["number", "string"],
      schemaType: "string",
      $data: true,
      error,
      code(cxt, ruleType) {
        const { gen, data, $data, schema: schema2, schemaCode, it } = cxt;
        const { opts, errSchemaPath, schemaEnv, self } = it;
        if (!opts.validateFormats)
          return;
        if ($data)
          validate$DataFormat();
        else
          validateFormat();
        function validate$DataFormat() {
          const fmts = gen.scopeValue("formats", {
            ref: self.formats,
            code: opts.code.formats
          });
          const fDef = gen.const("fDef", (0, codegen_1._)`${fmts}[${schemaCode}]`);
          const fType = gen.let("fType");
          const format = gen.let("format");
          gen.if((0, codegen_1._)`typeof ${fDef} == "object" && !(${fDef} instanceof RegExp)`, () => gen.assign(fType, (0, codegen_1._)`${fDef}.type || "string"`).assign(format, (0, codegen_1._)`${fDef}.validate`), () => gen.assign(fType, (0, codegen_1._)`"string"`).assign(format, fDef));
          cxt.fail$data((0, codegen_1.or)(unknownFmt(), invalidFmt()));
          function unknownFmt() {
            if (opts.strictSchema === false)
              return codegen_1.nil;
            return (0, codegen_1._)`${schemaCode} && !${format}`;
          }
          function invalidFmt() {
            const callFormat = schemaEnv.$async ? (0, codegen_1._)`(${fDef}.async ? await ${format}(${data}) : ${format}(${data}))` : (0, codegen_1._)`${format}(${data})`;
            const validData = (0, codegen_1._)`(typeof ${format} == "function" ? ${callFormat} : ${format}.test(${data}))`;
            return (0, codegen_1._)`${format} && ${format} !== true && ${fType} === ${ruleType} && !${validData}`;
          }
        }
        function validateFormat() {
          const formatDef = self.formats[schema2];
          if (!formatDef) {
            unknownFormat();
            return;
          }
          if (formatDef === true)
            return;
          const [fmtType, format, fmtRef] = getFormat(formatDef);
          if (fmtType === ruleType)
            cxt.pass(validCondition());
          function unknownFormat() {
            if (opts.strictSchema === false) {
              self.logger.warn(unknownMsg());
              return;
            }
            throw new Error(unknownMsg());
            function unknownMsg() {
              return `unknown format "${schema2}" ignored in schema at path "${errSchemaPath}"`;
            }
          }
          function getFormat(fmtDef) {
            const code = fmtDef instanceof RegExp ? (0, codegen_1.regexpCode)(fmtDef) : opts.code.formats ? (0, codegen_1._)`${opts.code.formats}${(0, codegen_1.getProperty)(schema2)}` : void 0;
            const fmt = gen.scopeValue("formats", { key: schema2, ref: fmtDef, code });
            if (typeof fmtDef == "object" && !(fmtDef instanceof RegExp)) {
              return [fmtDef.type || "string", fmtDef.validate, (0, codegen_1._)`${fmt}.validate`];
            }
            return ["string", fmtDef, fmt];
          }
          function validCondition() {
            if (typeof formatDef == "object" && !(formatDef instanceof RegExp) && formatDef.async) {
              if (!schemaEnv.$async)
                throw new Error("async format in sync schema");
              return (0, codegen_1._)`await ${fmtRef}(${data})`;
            }
            return typeof format == "function" ? (0, codegen_1._)`${fmtRef}(${data})` : (0, codegen_1._)`${fmtRef}.test(${data})`;
          }
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/vocabularies/format/index.js
var require_format2 = __commonJS({
  "node_modules/ajv/dist/vocabularies/format/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var format_1 = require_format();
    var format = [format_1.default];
    exports.default = format;
  }
});

// node_modules/ajv/dist/vocabularies/metadata.js
var require_metadata = __commonJS({
  "node_modules/ajv/dist/vocabularies/metadata.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.contentVocabulary = exports.metadataVocabulary = void 0;
    exports.metadataVocabulary = [
      "title",
      "description",
      "default",
      "deprecated",
      "readOnly",
      "writeOnly",
      "examples"
    ];
    exports.contentVocabulary = [
      "contentMediaType",
      "contentEncoding",
      "contentSchema"
    ];
  }
});

// node_modules/ajv/dist/vocabularies/draft2020.js
var require_draft2020 = __commonJS({
  "node_modules/ajv/dist/vocabularies/draft2020.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var core_1 = require_core2();
    var validation_1 = require_validation();
    var applicator_1 = require_applicator();
    var dynamic_1 = require_dynamic();
    var next_1 = require_next();
    var unevaluated_1 = require_unevaluated();
    var format_1 = require_format2();
    var metadata_1 = require_metadata();
    var draft2020Vocabularies = [
      dynamic_1.default,
      core_1.default,
      validation_1.default,
      (0, applicator_1.default)(true),
      format_1.default,
      metadata_1.metadataVocabulary,
      metadata_1.contentVocabulary,
      next_1.default,
      unevaluated_1.default
    ];
    exports.default = draft2020Vocabularies;
  }
});

// node_modules/ajv/dist/vocabularies/discriminator/types.js
var require_types = __commonJS({
  "node_modules/ajv/dist/vocabularies/discriminator/types.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.DiscrError = void 0;
    var DiscrError;
    (function(DiscrError2) {
      DiscrError2["Tag"] = "tag";
      DiscrError2["Mapping"] = "mapping";
    })(DiscrError || (exports.DiscrError = DiscrError = {}));
  }
});

// node_modules/ajv/dist/vocabularies/discriminator/index.js
var require_discriminator = __commonJS({
  "node_modules/ajv/dist/vocabularies/discriminator/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var codegen_1 = require_codegen();
    var types_1 = require_types();
    var compile_1 = require_compile();
    var ref_error_1 = require_ref_error();
    var util_1 = require_util();
    var error = {
      message: ({ params: { discrError, tagName } }) => discrError === types_1.DiscrError.Tag ? `tag "${tagName}" must be string` : `value of tag "${tagName}" must be in oneOf`,
      params: ({ params: { discrError, tag, tagName } }) => (0, codegen_1._)`{error: ${discrError}, tag: ${tagName}, tagValue: ${tag}}`
    };
    var def = {
      keyword: "discriminator",
      type: "object",
      schemaType: "object",
      error,
      code(cxt) {
        const { gen, data, schema: schema2, parentSchema, it } = cxt;
        const { oneOf } = parentSchema;
        if (!it.opts.discriminator) {
          throw new Error("discriminator: requires discriminator option");
        }
        const tagName = schema2.propertyName;
        if (typeof tagName != "string")
          throw new Error("discriminator: requires propertyName");
        if (schema2.mapping)
          throw new Error("discriminator: mapping is not supported");
        if (!oneOf)
          throw new Error("discriminator: requires oneOf keyword");
        const valid = gen.let("valid", false);
        const tag = gen.const("tag", (0, codegen_1._)`${data}${(0, codegen_1.getProperty)(tagName)}`);
        gen.if((0, codegen_1._)`typeof ${tag} == "string"`, () => validateMapping(), () => cxt.error(false, { discrError: types_1.DiscrError.Tag, tag, tagName }));
        cxt.ok(valid);
        function validateMapping() {
          const mapping = getMapping();
          gen.if(false);
          for (const tagValue in mapping) {
            gen.elseIf((0, codegen_1._)`${tag} === ${tagValue}`);
            gen.assign(valid, applyTagSchema(mapping[tagValue]));
          }
          gen.else();
          cxt.error(false, { discrError: types_1.DiscrError.Mapping, tag, tagName });
          gen.endIf();
        }
        function applyTagSchema(schemaProp) {
          const _valid = gen.name("valid");
          const schCxt = cxt.subschema({ keyword: "oneOf", schemaProp }, _valid);
          cxt.mergeEvaluated(schCxt, codegen_1.Name);
          return _valid;
        }
        function getMapping() {
          var _a;
          const oneOfMapping = {};
          const topRequired = hasRequired(parentSchema);
          let tagRequired = true;
          for (let i = 0; i < oneOf.length; i++) {
            let sch = oneOf[i];
            if ((sch === null || sch === void 0 ? void 0 : sch.$ref) && !(0, util_1.schemaHasRulesButRef)(sch, it.self.RULES)) {
              const ref = sch.$ref;
              sch = compile_1.resolveRef.call(it.self, it.schemaEnv.root, it.baseId, ref);
              if (sch instanceof compile_1.SchemaEnv)
                sch = sch.schema;
              if (sch === void 0)
                throw new ref_error_1.default(it.opts.uriResolver, it.baseId, ref);
            }
            const propSch = (_a = sch === null || sch === void 0 ? void 0 : sch.properties) === null || _a === void 0 ? void 0 : _a[tagName];
            if (typeof propSch != "object") {
              throw new Error(`discriminator: oneOf subschemas (or referenced schemas) must have "properties/${tagName}"`);
            }
            tagRequired = tagRequired && (topRequired || hasRequired(sch));
            addMappings(propSch, i);
          }
          if (!tagRequired)
            throw new Error(`discriminator: "${tagName}" must be required`);
          return oneOfMapping;
          function hasRequired({ required }) {
            return Array.isArray(required) && required.includes(tagName);
          }
          function addMappings(sch, i) {
            if (sch.const) {
              addMapping(sch.const, i);
            } else if (sch.enum) {
              for (const tagValue of sch.enum) {
                addMapping(tagValue, i);
              }
            } else {
              throw new Error(`discriminator: "properties/${tagName}" must have "const" or "enum"`);
            }
          }
          function addMapping(tagValue, i) {
            if (typeof tagValue != "string" || tagValue in oneOfMapping) {
              throw new Error(`discriminator: "${tagName}" values must be unique strings`);
            }
            oneOfMapping[tagValue] = i;
          }
        }
      }
    };
    exports.default = def;
  }
});

// node_modules/ajv/dist/refs/json-schema-2020-12/schema.json
var require_schema = __commonJS({
  "node_modules/ajv/dist/refs/json-schema-2020-12/schema.json"(exports, module) {
    module.exports = {
      $schema: "https://json-schema.org/draft/2020-12/schema",
      $id: "https://json-schema.org/draft/2020-12/schema",
      $vocabulary: {
        "https://json-schema.org/draft/2020-12/vocab/core": true,
        "https://json-schema.org/draft/2020-12/vocab/applicator": true,
        "https://json-schema.org/draft/2020-12/vocab/unevaluated": true,
        "https://json-schema.org/draft/2020-12/vocab/validation": true,
        "https://json-schema.org/draft/2020-12/vocab/meta-data": true,
        "https://json-schema.org/draft/2020-12/vocab/format-annotation": true,
        "https://json-schema.org/draft/2020-12/vocab/content": true
      },
      $dynamicAnchor: "meta",
      title: "Core and Validation specifications meta-schema",
      allOf: [
        { $ref: "meta/core" },
        { $ref: "meta/applicator" },
        { $ref: "meta/unevaluated" },
        { $ref: "meta/validation" },
        { $ref: "meta/meta-data" },
        { $ref: "meta/format-annotation" },
        { $ref: "meta/content" }
      ],
      type: ["object", "boolean"],
      $comment: "This meta-schema also defines keywords that have appeared in previous drafts in order to prevent incompatible extensions as they remain in common use.",
      properties: {
        definitions: {
          $comment: '"definitions" has been replaced by "$defs".',
          type: "object",
          additionalProperties: { $dynamicRef: "#meta" },
          deprecated: true,
          default: {}
        },
        dependencies: {
          $comment: '"dependencies" has been split and replaced by "dependentSchemas" and "dependentRequired" in order to serve their differing semantics.',
          type: "object",
          additionalProperties: {
            anyOf: [{ $dynamicRef: "#meta" }, { $ref: "meta/validation#/$defs/stringArray" }]
          },
          deprecated: true,
          default: {}
        },
        $recursiveAnchor: {
          $comment: '"$recursiveAnchor" has been replaced by "$dynamicAnchor".',
          $ref: "meta/core#/$defs/anchorString",
          deprecated: true
        },
        $recursiveRef: {
          $comment: '"$recursiveRef" has been replaced by "$dynamicRef".',
          $ref: "meta/core#/$defs/uriReferenceString",
          deprecated: true
        }
      }
    };
  }
});

// node_modules/ajv/dist/refs/json-schema-2020-12/meta/applicator.json
var require_applicator2 = __commonJS({
  "node_modules/ajv/dist/refs/json-schema-2020-12/meta/applicator.json"(exports, module) {
    module.exports = {
      $schema: "https://json-schema.org/draft/2020-12/schema",
      $id: "https://json-schema.org/draft/2020-12/meta/applicator",
      $vocabulary: {
        "https://json-schema.org/draft/2020-12/vocab/applicator": true
      },
      $dynamicAnchor: "meta",
      title: "Applicator vocabulary meta-schema",
      type: ["object", "boolean"],
      properties: {
        prefixItems: { $ref: "#/$defs/schemaArray" },
        items: { $dynamicRef: "#meta" },
        contains: { $dynamicRef: "#meta" },
        additionalProperties: { $dynamicRef: "#meta" },
        properties: {
          type: "object",
          additionalProperties: { $dynamicRef: "#meta" },
          default: {}
        },
        patternProperties: {
          type: "object",
          additionalProperties: { $dynamicRef: "#meta" },
          propertyNames: { format: "regex" },
          default: {}
        },
        dependentSchemas: {
          type: "object",
          additionalProperties: { $dynamicRef: "#meta" },
          default: {}
        },
        propertyNames: { $dynamicRef: "#meta" },
        if: { $dynamicRef: "#meta" },
        then: { $dynamicRef: "#meta" },
        else: { $dynamicRef: "#meta" },
        allOf: { $ref: "#/$defs/schemaArray" },
        anyOf: { $ref: "#/$defs/schemaArray" },
        oneOf: { $ref: "#/$defs/schemaArray" },
        not: { $dynamicRef: "#meta" }
      },
      $defs: {
        schemaArray: {
          type: "array",
          minItems: 1,
          items: { $dynamicRef: "#meta" }
        }
      }
    };
  }
});

// node_modules/ajv/dist/refs/json-schema-2020-12/meta/unevaluated.json
var require_unevaluated2 = __commonJS({
  "node_modules/ajv/dist/refs/json-schema-2020-12/meta/unevaluated.json"(exports, module) {
    module.exports = {
      $schema: "https://json-schema.org/draft/2020-12/schema",
      $id: "https://json-schema.org/draft/2020-12/meta/unevaluated",
      $vocabulary: {
        "https://json-schema.org/draft/2020-12/vocab/unevaluated": true
      },
      $dynamicAnchor: "meta",
      title: "Unevaluated applicator vocabulary meta-schema",
      type: ["object", "boolean"],
      properties: {
        unevaluatedItems: { $dynamicRef: "#meta" },
        unevaluatedProperties: { $dynamicRef: "#meta" }
      }
    };
  }
});

// node_modules/ajv/dist/refs/json-schema-2020-12/meta/content.json
var require_content = __commonJS({
  "node_modules/ajv/dist/refs/json-schema-2020-12/meta/content.json"(exports, module) {
    module.exports = {
      $schema: "https://json-schema.org/draft/2020-12/schema",
      $id: "https://json-schema.org/draft/2020-12/meta/content",
      $vocabulary: {
        "https://json-schema.org/draft/2020-12/vocab/content": true
      },
      $dynamicAnchor: "meta",
      title: "Content vocabulary meta-schema",
      type: ["object", "boolean"],
      properties: {
        contentEncoding: { type: "string" },
        contentMediaType: { type: "string" },
        contentSchema: { $dynamicRef: "#meta" }
      }
    };
  }
});

// node_modules/ajv/dist/refs/json-schema-2020-12/meta/core.json
var require_core3 = __commonJS({
  "node_modules/ajv/dist/refs/json-schema-2020-12/meta/core.json"(exports, module) {
    module.exports = {
      $schema: "https://json-schema.org/draft/2020-12/schema",
      $id: "https://json-schema.org/draft/2020-12/meta/core",
      $vocabulary: {
        "https://json-schema.org/draft/2020-12/vocab/core": true
      },
      $dynamicAnchor: "meta",
      title: "Core vocabulary meta-schema",
      type: ["object", "boolean"],
      properties: {
        $id: {
          $ref: "#/$defs/uriReferenceString",
          $comment: "Non-empty fragments not allowed.",
          pattern: "^[^#]*#?$"
        },
        $schema: { $ref: "#/$defs/uriString" },
        $ref: { $ref: "#/$defs/uriReferenceString" },
        $anchor: { $ref: "#/$defs/anchorString" },
        $dynamicRef: { $ref: "#/$defs/uriReferenceString" },
        $dynamicAnchor: { $ref: "#/$defs/anchorString" },
        $vocabulary: {
          type: "object",
          propertyNames: { $ref: "#/$defs/uriString" },
          additionalProperties: {
            type: "boolean"
          }
        },
        $comment: {
          type: "string"
        },
        $defs: {
          type: "object",
          additionalProperties: { $dynamicRef: "#meta" }
        }
      },
      $defs: {
        anchorString: {
          type: "string",
          pattern: "^[A-Za-z_][-A-Za-z0-9._]*$"
        },
        uriString: {
          type: "string",
          format: "uri"
        },
        uriReferenceString: {
          type: "string",
          format: "uri-reference"
        }
      }
    };
  }
});

// node_modules/ajv/dist/refs/json-schema-2020-12/meta/format-annotation.json
var require_format_annotation = __commonJS({
  "node_modules/ajv/dist/refs/json-schema-2020-12/meta/format-annotation.json"(exports, module) {
    module.exports = {
      $schema: "https://json-schema.org/draft/2020-12/schema",
      $id: "https://json-schema.org/draft/2020-12/meta/format-annotation",
      $vocabulary: {
        "https://json-schema.org/draft/2020-12/vocab/format-annotation": true
      },
      $dynamicAnchor: "meta",
      title: "Format vocabulary meta-schema for annotation results",
      type: ["object", "boolean"],
      properties: {
        format: { type: "string" }
      }
    };
  }
});

// node_modules/ajv/dist/refs/json-schema-2020-12/meta/meta-data.json
var require_meta_data = __commonJS({
  "node_modules/ajv/dist/refs/json-schema-2020-12/meta/meta-data.json"(exports, module) {
    module.exports = {
      $schema: "https://json-schema.org/draft/2020-12/schema",
      $id: "https://json-schema.org/draft/2020-12/meta/meta-data",
      $vocabulary: {
        "https://json-schema.org/draft/2020-12/vocab/meta-data": true
      },
      $dynamicAnchor: "meta",
      title: "Meta-data vocabulary meta-schema",
      type: ["object", "boolean"],
      properties: {
        title: {
          type: "string"
        },
        description: {
          type: "string"
        },
        default: true,
        deprecated: {
          type: "boolean",
          default: false
        },
        readOnly: {
          type: "boolean",
          default: false
        },
        writeOnly: {
          type: "boolean",
          default: false
        },
        examples: {
          type: "array",
          items: true
        }
      }
    };
  }
});

// node_modules/ajv/dist/refs/json-schema-2020-12/meta/validation.json
var require_validation2 = __commonJS({
  "node_modules/ajv/dist/refs/json-schema-2020-12/meta/validation.json"(exports, module) {
    module.exports = {
      $schema: "https://json-schema.org/draft/2020-12/schema",
      $id: "https://json-schema.org/draft/2020-12/meta/validation",
      $vocabulary: {
        "https://json-schema.org/draft/2020-12/vocab/validation": true
      },
      $dynamicAnchor: "meta",
      title: "Validation vocabulary meta-schema",
      type: ["object", "boolean"],
      properties: {
        type: {
          anyOf: [
            { $ref: "#/$defs/simpleTypes" },
            {
              type: "array",
              items: { $ref: "#/$defs/simpleTypes" },
              minItems: 1,
              uniqueItems: true
            }
          ]
        },
        const: true,
        enum: {
          type: "array",
          items: true
        },
        multipleOf: {
          type: "number",
          exclusiveMinimum: 0
        },
        maximum: {
          type: "number"
        },
        exclusiveMaximum: {
          type: "number"
        },
        minimum: {
          type: "number"
        },
        exclusiveMinimum: {
          type: "number"
        },
        maxLength: { $ref: "#/$defs/nonNegativeInteger" },
        minLength: { $ref: "#/$defs/nonNegativeIntegerDefault0" },
        pattern: {
          type: "string",
          format: "regex"
        },
        maxItems: { $ref: "#/$defs/nonNegativeInteger" },
        minItems: { $ref: "#/$defs/nonNegativeIntegerDefault0" },
        uniqueItems: {
          type: "boolean",
          default: false
        },
        maxContains: { $ref: "#/$defs/nonNegativeInteger" },
        minContains: {
          $ref: "#/$defs/nonNegativeInteger",
          default: 1
        },
        maxProperties: { $ref: "#/$defs/nonNegativeInteger" },
        minProperties: { $ref: "#/$defs/nonNegativeIntegerDefault0" },
        required: { $ref: "#/$defs/stringArray" },
        dependentRequired: {
          type: "object",
          additionalProperties: {
            $ref: "#/$defs/stringArray"
          }
        }
      },
      $defs: {
        nonNegativeInteger: {
          type: "integer",
          minimum: 0
        },
        nonNegativeIntegerDefault0: {
          $ref: "#/$defs/nonNegativeInteger",
          default: 0
        },
        simpleTypes: {
          enum: ["array", "boolean", "integer", "null", "number", "object", "string"]
        },
        stringArray: {
          type: "array",
          items: { type: "string" },
          uniqueItems: true,
          default: []
        }
      }
    };
  }
});

// node_modules/ajv/dist/refs/json-schema-2020-12/index.js
var require_json_schema_2020_12 = __commonJS({
  "node_modules/ajv/dist/refs/json-schema-2020-12/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    var metaSchema = require_schema();
    var applicator = require_applicator2();
    var unevaluated = require_unevaluated2();
    var content = require_content();
    var core = require_core3();
    var format = require_format_annotation();
    var metadata = require_meta_data();
    var validation = require_validation2();
    var META_SUPPORT_DATA = ["/properties"];
    function addMetaSchema2020($data) {
      ;
      [
        metaSchema,
        applicator,
        unevaluated,
        content,
        core,
        with$data(this, format),
        metadata,
        with$data(this, validation)
      ].forEach((sch) => this.addMetaSchema(sch, void 0, false));
      return this;
      function with$data(ajv, sch) {
        return $data ? ajv.$dataMetaSchema(sch, META_SUPPORT_DATA) : sch;
      }
    }
    exports.default = addMetaSchema2020;
  }
});

// node_modules/ajv/dist/2020.js
var require__ = __commonJS({
  "node_modules/ajv/dist/2020.js"(exports, module) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.MissingRefError = exports.ValidationError = exports.CodeGen = exports.Name = exports.nil = exports.stringify = exports.str = exports._ = exports.KeywordCxt = exports.Ajv2020 = void 0;
    var core_1 = require_core();
    var draft2020_1 = require_draft2020();
    var discriminator_1 = require_discriminator();
    var json_schema_2020_12_1 = require_json_schema_2020_12();
    var META_SCHEMA_ID = "https://json-schema.org/draft/2020-12/schema";
    var Ajv20202 = class extends core_1.default {
      constructor(opts = {}) {
        super({
          ...opts,
          dynamicRef: true,
          next: true,
          unevaluated: true
        });
      }
      _addVocabularies() {
        super._addVocabularies();
        draft2020_1.default.forEach((v) => this.addVocabulary(v));
        if (this.opts.discriminator)
          this.addKeyword(discriminator_1.default);
      }
      _addDefaultMetaSchema() {
        super._addDefaultMetaSchema();
        const { $data, meta } = this.opts;
        if (!meta)
          return;
        json_schema_2020_12_1.default.call(this, $data);
        this.refs["http://json-schema.org/schema"] = META_SCHEMA_ID;
      }
      defaultMeta() {
        return this.opts.defaultMeta = super.defaultMeta() || (this.getSchema(META_SCHEMA_ID) ? META_SCHEMA_ID : void 0);
      }
    };
    exports.Ajv2020 = Ajv20202;
    module.exports = exports = Ajv20202;
    module.exports.Ajv2020 = Ajv20202;
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.default = Ajv20202;
    var validate_1 = require_validate();
    Object.defineProperty(exports, "KeywordCxt", { enumerable: true, get: function() {
      return validate_1.KeywordCxt;
    } });
    var codegen_1 = require_codegen();
    Object.defineProperty(exports, "_", { enumerable: true, get: function() {
      return codegen_1._;
    } });
    Object.defineProperty(exports, "str", { enumerable: true, get: function() {
      return codegen_1.str;
    } });
    Object.defineProperty(exports, "stringify", { enumerable: true, get: function() {
      return codegen_1.stringify;
    } });
    Object.defineProperty(exports, "nil", { enumerable: true, get: function() {
      return codegen_1.nil;
    } });
    Object.defineProperty(exports, "Name", { enumerable: true, get: function() {
      return codegen_1.Name;
    } });
    Object.defineProperty(exports, "CodeGen", { enumerable: true, get: function() {
      return codegen_1.CodeGen;
    } });
    var validation_error_1 = require_validation_error();
    Object.defineProperty(exports, "ValidationError", { enumerable: true, get: function() {
      return validation_error_1.default;
    } });
    var ref_error_1 = require_ref_error();
    Object.defineProperty(exports, "MissingRefError", { enumerable: true, get: function() {
      return ref_error_1.default;
    } });
  }
});

// node_modules/yaml/dist/nodes/identity.js
var require_identity = __commonJS({
  "node_modules/yaml/dist/nodes/identity.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var ALIAS = /* @__PURE__ */ Symbol.for("yaml.alias");
    var DOC = /* @__PURE__ */ Symbol.for("yaml.document");
    var MAP = /* @__PURE__ */ Symbol.for("yaml.map");
    var PAIR = /* @__PURE__ */ Symbol.for("yaml.pair");
    var SCALAR = /* @__PURE__ */ Symbol.for("yaml.scalar");
    var SEQ = /* @__PURE__ */ Symbol.for("yaml.seq");
    var NODE_TYPE = /* @__PURE__ */ Symbol.for("yaml.node.type");
    var isAlias = (node) => !!node && typeof node === "object" && node[NODE_TYPE] === ALIAS;
    var isDocument = (node) => !!node && typeof node === "object" && node[NODE_TYPE] === DOC;
    var isMap = (node) => !!node && typeof node === "object" && node[NODE_TYPE] === MAP;
    var isPair = (node) => !!node && typeof node === "object" && node[NODE_TYPE] === PAIR;
    var isScalar = (node) => !!node && typeof node === "object" && node[NODE_TYPE] === SCALAR;
    var isSeq = (node) => !!node && typeof node === "object" && node[NODE_TYPE] === SEQ;
    function isCollection(node) {
      if (node && typeof node === "object")
        switch (node[NODE_TYPE]) {
          case MAP:
          case SEQ:
            return true;
        }
      return false;
    }
    function isNode(node) {
      if (node && typeof node === "object")
        switch (node[NODE_TYPE]) {
          case ALIAS:
          case MAP:
          case SCALAR:
          case SEQ:
            return true;
        }
      return false;
    }
    var hasAnchor = (node) => (isScalar(node) || isCollection(node)) && !!node.anchor;
    exports.ALIAS = ALIAS;
    exports.DOC = DOC;
    exports.MAP = MAP;
    exports.NODE_TYPE = NODE_TYPE;
    exports.PAIR = PAIR;
    exports.SCALAR = SCALAR;
    exports.SEQ = SEQ;
    exports.hasAnchor = hasAnchor;
    exports.isAlias = isAlias;
    exports.isCollection = isCollection;
    exports.isDocument = isDocument;
    exports.isMap = isMap;
    exports.isNode = isNode;
    exports.isPair = isPair;
    exports.isScalar = isScalar;
    exports.isSeq = isSeq;
  }
});

// node_modules/yaml/dist/visit.js
var require_visit = __commonJS({
  "node_modules/yaml/dist/visit.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var BREAK = /* @__PURE__ */ Symbol("break visit");
    var SKIP = /* @__PURE__ */ Symbol("skip children");
    var REMOVE = /* @__PURE__ */ Symbol("remove node");
    function visit(node, visitor) {
      const visitor_ = initVisitor(visitor);
      if (identity.isDocument(node)) {
        const cd = visit_(null, node.contents, visitor_, Object.freeze([node]));
        if (cd === REMOVE)
          node.contents = null;
      } else
        visit_(null, node, visitor_, Object.freeze([]));
    }
    visit.BREAK = BREAK;
    visit.SKIP = SKIP;
    visit.REMOVE = REMOVE;
    function visit_(key, node, visitor, path2) {
      const ctrl = callVisitor(key, node, visitor, path2);
      if (identity.isNode(ctrl) || identity.isPair(ctrl)) {
        replaceNode(key, path2, ctrl);
        return visit_(key, ctrl, visitor, path2);
      }
      if (typeof ctrl !== "symbol") {
        if (identity.isCollection(node)) {
          path2 = Object.freeze(path2.concat(node));
          for (let i = 0; i < node.items.length; ++i) {
            const ci = visit_(i, node.items[i], visitor, path2);
            if (typeof ci === "number")
              i = ci - 1;
            else if (ci === BREAK)
              return BREAK;
            else if (ci === REMOVE) {
              node.items.splice(i, 1);
              i -= 1;
            }
          }
        } else if (identity.isPair(node)) {
          path2 = Object.freeze(path2.concat(node));
          const ck = visit_("key", node.key, visitor, path2);
          if (ck === BREAK)
            return BREAK;
          else if (ck === REMOVE)
            node.key = null;
          const cv = visit_("value", node.value, visitor, path2);
          if (cv === BREAK)
            return BREAK;
          else if (cv === REMOVE)
            node.value = null;
        }
      }
      return ctrl;
    }
    async function visitAsync(node, visitor) {
      const visitor_ = initVisitor(visitor);
      if (identity.isDocument(node)) {
        const cd = await visitAsync_(null, node.contents, visitor_, Object.freeze([node]));
        if (cd === REMOVE)
          node.contents = null;
      } else
        await visitAsync_(null, node, visitor_, Object.freeze([]));
    }
    visitAsync.BREAK = BREAK;
    visitAsync.SKIP = SKIP;
    visitAsync.REMOVE = REMOVE;
    async function visitAsync_(key, node, visitor, path2) {
      const ctrl = await callVisitor(key, node, visitor, path2);
      if (identity.isNode(ctrl) || identity.isPair(ctrl)) {
        replaceNode(key, path2, ctrl);
        return visitAsync_(key, ctrl, visitor, path2);
      }
      if (typeof ctrl !== "symbol") {
        if (identity.isCollection(node)) {
          path2 = Object.freeze(path2.concat(node));
          for (let i = 0; i < node.items.length; ++i) {
            const ci = await visitAsync_(i, node.items[i], visitor, path2);
            if (typeof ci === "number")
              i = ci - 1;
            else if (ci === BREAK)
              return BREAK;
            else if (ci === REMOVE) {
              node.items.splice(i, 1);
              i -= 1;
            }
          }
        } else if (identity.isPair(node)) {
          path2 = Object.freeze(path2.concat(node));
          const ck = await visitAsync_("key", node.key, visitor, path2);
          if (ck === BREAK)
            return BREAK;
          else if (ck === REMOVE)
            node.key = null;
          const cv = await visitAsync_("value", node.value, visitor, path2);
          if (cv === BREAK)
            return BREAK;
          else if (cv === REMOVE)
            node.value = null;
        }
      }
      return ctrl;
    }
    function initVisitor(visitor) {
      if (typeof visitor === "object" && (visitor.Collection || visitor.Node || visitor.Value)) {
        return Object.assign({
          Alias: visitor.Node,
          Map: visitor.Node,
          Scalar: visitor.Node,
          Seq: visitor.Node
        }, visitor.Value && {
          Map: visitor.Value,
          Scalar: visitor.Value,
          Seq: visitor.Value
        }, visitor.Collection && {
          Map: visitor.Collection,
          Seq: visitor.Collection
        }, visitor);
      }
      return visitor;
    }
    function callVisitor(key, node, visitor, path2) {
      if (typeof visitor === "function")
        return visitor(key, node, path2);
      if (identity.isMap(node))
        return visitor.Map?.(key, node, path2);
      if (identity.isSeq(node))
        return visitor.Seq?.(key, node, path2);
      if (identity.isPair(node))
        return visitor.Pair?.(key, node, path2);
      if (identity.isScalar(node))
        return visitor.Scalar?.(key, node, path2);
      if (identity.isAlias(node))
        return visitor.Alias?.(key, node, path2);
      return void 0;
    }
    function replaceNode(key, path2, node) {
      const parent = path2[path2.length - 1];
      if (identity.isCollection(parent)) {
        parent.items[key] = node;
      } else if (identity.isPair(parent)) {
        if (key === "key")
          parent.key = node;
        else
          parent.value = node;
      } else if (identity.isDocument(parent)) {
        parent.contents = node;
      } else {
        const pt = identity.isAlias(parent) ? "alias" : "scalar";
        throw new Error(`Cannot replace node with ${pt} parent`);
      }
    }
    exports.visit = visit;
    exports.visitAsync = visitAsync;
  }
});

// node_modules/yaml/dist/doc/directives.js
var require_directives = __commonJS({
  "node_modules/yaml/dist/doc/directives.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var visit = require_visit();
    var escapeChars = {
      "!": "%21",
      ",": "%2C",
      "[": "%5B",
      "]": "%5D",
      "{": "%7B",
      "}": "%7D"
    };
    var escapeTagName = (tn) => tn.replace(/[!,[\]{}]/g, (ch) => escapeChars[ch]);
    var Directives = class _Directives {
      constructor(yaml, tags) {
        this.docStart = null;
        this.docEnd = false;
        this.yaml = Object.assign({}, _Directives.defaultYaml, yaml);
        this.tags = Object.assign({}, _Directives.defaultTags, tags);
      }
      clone() {
        const copy = new _Directives(this.yaml, this.tags);
        copy.docStart = this.docStart;
        return copy;
      }
      /**
       * During parsing, get a Directives instance for the current document and
       * update the stream state according to the current version's spec.
       */
      atDocument() {
        const res = new _Directives(this.yaml, this.tags);
        switch (this.yaml.version) {
          case "1.1":
            this.atNextDocument = true;
            break;
          case "1.2":
            this.atNextDocument = false;
            this.yaml = {
              explicit: _Directives.defaultYaml.explicit,
              version: "1.2"
            };
            this.tags = Object.assign({}, _Directives.defaultTags);
            break;
        }
        return res;
      }
      /**
       * @param onError - May be called even if the action was successful
       * @returns `true` on success
       */
      add(line, onError) {
        if (this.atNextDocument) {
          this.yaml = { explicit: _Directives.defaultYaml.explicit, version: "1.1" };
          this.tags = Object.assign({}, _Directives.defaultTags);
          this.atNextDocument = false;
        }
        const parts = line.trim().split(/[ \t]+/);
        const name = parts.shift();
        switch (name) {
          case "%TAG": {
            if (parts.length !== 2) {
              onError(0, "%TAG directive should contain exactly two parts");
              if (parts.length < 2)
                return false;
            }
            const [handle, prefix] = parts;
            this.tags[handle] = prefix;
            return true;
          }
          case "%YAML": {
            this.yaml.explicit = true;
            if (parts.length !== 1) {
              onError(0, "%YAML directive should contain exactly one part");
              return false;
            }
            const [version] = parts;
            if (version === "1.1" || version === "1.2") {
              this.yaml.version = version;
              return true;
            } else {
              const isValid = /^\d+\.\d+$/.test(version);
              onError(6, `Unsupported YAML version ${version}`, isValid);
              return false;
            }
          }
          default:
            onError(0, `Unknown directive ${name}`, true);
            return false;
        }
      }
      /**
       * Resolves a tag, matching handles to those defined in %TAG directives.
       *
       * @returns Resolved tag, which may also be the non-specific tag `'!'` or a
       *   `'!local'` tag, or `null` if unresolvable.
       */
      tagName(source, onError) {
        if (source === "!")
          return "!";
        if (source[0] !== "!") {
          onError(`Not a valid tag: ${source}`);
          return null;
        }
        if (source[1] === "<") {
          const verbatim = source.slice(2, -1);
          if (verbatim === "!" || verbatim === "!!") {
            onError(`Verbatim tags aren't resolved, so ${source} is invalid.`);
            return null;
          }
          if (source[source.length - 1] !== ">")
            onError("Verbatim tags must end with a >");
          return verbatim;
        }
        const [, handle, suffix] = source.match(/^(.*!)([^!]*)$/s);
        if (!suffix)
          onError(`The ${source} tag has no suffix`);
        const prefix = this.tags[handle];
        if (prefix) {
          try {
            return prefix + decodeURIComponent(suffix);
          } catch (error) {
            onError(String(error));
            return null;
          }
        }
        if (handle === "!")
          return source;
        onError(`Could not resolve tag: ${source}`);
        return null;
      }
      /**
       * Given a fully resolved tag, returns its printable string form,
       * taking into account current tag prefixes and defaults.
       */
      tagString(tag) {
        for (const [handle, prefix] of Object.entries(this.tags)) {
          if (tag.startsWith(prefix))
            return handle + escapeTagName(tag.substring(prefix.length));
        }
        return tag[0] === "!" ? tag : `!<${tag}>`;
      }
      toString(doc) {
        const lines = this.yaml.explicit ? [`%YAML ${this.yaml.version || "1.2"}`] : [];
        const tagEntries = Object.entries(this.tags);
        let tagNames;
        if (doc && tagEntries.length > 0 && identity.isNode(doc.contents)) {
          const tags = {};
          visit.visit(doc.contents, (_key, node) => {
            if (identity.isNode(node) && node.tag)
              tags[node.tag] = true;
          });
          tagNames = Object.keys(tags);
        } else
          tagNames = [];
        for (const [handle, prefix] of tagEntries) {
          if (handle === "!!" && prefix === "tag:yaml.org,2002:")
            continue;
          if (!doc || tagNames.some((tn) => tn.startsWith(prefix)))
            lines.push(`%TAG ${handle} ${prefix}`);
        }
        return lines.join("\n");
      }
    };
    Directives.defaultYaml = { explicit: false, version: "1.2" };
    Directives.defaultTags = { "!!": "tag:yaml.org,2002:" };
    exports.Directives = Directives;
  }
});

// node_modules/yaml/dist/doc/anchors.js
var require_anchors = __commonJS({
  "node_modules/yaml/dist/doc/anchors.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var visit = require_visit();
    function anchorIsValid(anchor) {
      if (/[\x00-\x19\s,[\]{}]/.test(anchor)) {
        const sa = JSON.stringify(anchor);
        const msg = `Anchor must not contain whitespace or control characters: ${sa}`;
        throw new Error(msg);
      }
      return true;
    }
    function anchorNames(root) {
      const anchors = /* @__PURE__ */ new Set();
      visit.visit(root, {
        Value(_key, node) {
          if (node.anchor)
            anchors.add(node.anchor);
        }
      });
      return anchors;
    }
    function findNewAnchor(prefix, exclude) {
      for (let i = 1; true; ++i) {
        const name = `${prefix}${i}`;
        if (!exclude.has(name))
          return name;
      }
    }
    function createNodeAnchors(doc, prefix) {
      const aliasObjects = [];
      const sourceObjects = /* @__PURE__ */ new Map();
      let prevAnchors = null;
      return {
        onAnchor: (source) => {
          aliasObjects.push(source);
          prevAnchors ?? (prevAnchors = anchorNames(doc));
          const anchor = findNewAnchor(prefix, prevAnchors);
          prevAnchors.add(anchor);
          return anchor;
        },
        /**
         * With circular references, the source node is only resolved after all
         * of its child nodes are. This is why anchors are set only after all of
         * the nodes have been created.
         */
        setAnchors: () => {
          for (const source of aliasObjects) {
            const ref = sourceObjects.get(source);
            if (typeof ref === "object" && ref.anchor && (identity.isScalar(ref.node) || identity.isCollection(ref.node))) {
              ref.node.anchor = ref.anchor;
            } else {
              const error = new Error("Failed to resolve repeated object (this should not happen)");
              error.source = source;
              throw error;
            }
          }
        },
        sourceObjects
      };
    }
    exports.anchorIsValid = anchorIsValid;
    exports.anchorNames = anchorNames;
    exports.createNodeAnchors = createNodeAnchors;
    exports.findNewAnchor = findNewAnchor;
  }
});

// node_modules/yaml/dist/doc/applyReviver.js
var require_applyReviver = __commonJS({
  "node_modules/yaml/dist/doc/applyReviver.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    function applyReviver(reviver, obj, key, val) {
      if (val && typeof val === "object") {
        if (Array.isArray(val)) {
          for (let i = 0, len = val.length; i < len; ++i) {
            const v0 = val[i];
            const v1 = applyReviver(reviver, val, String(i), v0);
            if (v1 === void 0)
              delete val[i];
            else if (v1 !== v0)
              val[i] = v1;
          }
        } else if (val instanceof Map) {
          for (const k of Array.from(val.keys())) {
            const v0 = val.get(k);
            const v1 = applyReviver(reviver, val, k, v0);
            if (v1 === void 0)
              val.delete(k);
            else if (v1 !== v0)
              val.set(k, v1);
          }
        } else if (val instanceof Set) {
          for (const v0 of Array.from(val)) {
            const v1 = applyReviver(reviver, val, v0, v0);
            if (v1 === void 0)
              val.delete(v0);
            else if (v1 !== v0) {
              val.delete(v0);
              val.add(v1);
            }
          }
        } else {
          for (const [k, v0] of Object.entries(val)) {
            const v1 = applyReviver(reviver, val, k, v0);
            if (v1 === void 0)
              delete val[k];
            else if (v1 !== v0)
              val[k] = v1;
          }
        }
      }
      return reviver.call(obj, key, val);
    }
    exports.applyReviver = applyReviver;
  }
});

// node_modules/yaml/dist/nodes/toJS.js
var require_toJS = __commonJS({
  "node_modules/yaml/dist/nodes/toJS.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    function toJS(value2, arg, ctx) {
      if (Array.isArray(value2))
        return value2.map((v, i) => toJS(v, String(i), ctx));
      if (value2 && typeof value2.toJSON === "function") {
        if (!ctx || !identity.hasAnchor(value2))
          return value2.toJSON(arg, ctx);
        const data = { aliasCount: 0, count: 1, res: void 0 };
        ctx.anchors.set(value2, data);
        ctx.onCreate = (res2) => {
          data.res = res2;
          delete ctx.onCreate;
        };
        const res = value2.toJSON(arg, ctx);
        if (ctx.onCreate)
          ctx.onCreate(res);
        return res;
      }
      if (typeof value2 === "bigint" && !ctx?.keep)
        return Number(value2);
      return value2;
    }
    exports.toJS = toJS;
  }
});

// node_modules/yaml/dist/nodes/Node.js
var require_Node = __commonJS({
  "node_modules/yaml/dist/nodes/Node.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var applyReviver = require_applyReviver();
    var identity = require_identity();
    var toJS = require_toJS();
    var NodeBase = class {
      constructor(type) {
        Object.defineProperty(this, identity.NODE_TYPE, { value: type });
      }
      /** Create a copy of this node.  */
      clone() {
        const copy = Object.create(Object.getPrototypeOf(this), Object.getOwnPropertyDescriptors(this));
        if (this.range)
          copy.range = this.range.slice();
        return copy;
      }
      /** A plain JavaScript representation of this node. */
      toJS(doc, { mapAsMap, maxAliasCount, onAnchor, reviver } = {}) {
        if (!identity.isDocument(doc))
          throw new TypeError("A document argument is required");
        const ctx = {
          anchors: /* @__PURE__ */ new Map(),
          doc,
          keep: true,
          mapAsMap: mapAsMap === true,
          mapKeyWarned: false,
          maxAliasCount: typeof maxAliasCount === "number" ? maxAliasCount : 100
        };
        const res = toJS.toJS(this, "", ctx);
        if (typeof onAnchor === "function")
          for (const { count, res: res2 } of ctx.anchors.values())
            onAnchor(res2, count);
        return typeof reviver === "function" ? applyReviver.applyReviver(reviver, { "": res }, "", res) : res;
      }
    };
    exports.NodeBase = NodeBase;
  }
});

// node_modules/yaml/dist/nodes/Alias.js
var require_Alias = __commonJS({
  "node_modules/yaml/dist/nodes/Alias.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var anchors = require_anchors();
    var visit = require_visit();
    var identity = require_identity();
    var Node = require_Node();
    var toJS = require_toJS();
    var Alias = class extends Node.NodeBase {
      constructor(source) {
        super(identity.ALIAS);
        this.source = source;
        Object.defineProperty(this, "tag", {
          set() {
            throw new Error("Alias nodes cannot have tags");
          }
        });
      }
      /**
       * Resolve the value of this alias within `doc`, finding the last
       * instance of the `source` anchor before this node.
       */
      resolve(doc, ctx) {
        if (ctx?.maxAliasCount === 0)
          throw new ReferenceError("Alias resolution is disabled");
        let nodes;
        if (ctx?.aliasResolveCache) {
          nodes = ctx.aliasResolveCache;
        } else {
          nodes = [];
          visit.visit(doc, {
            Node: (_key, node) => {
              if (identity.isAlias(node) || identity.hasAnchor(node))
                nodes.push(node);
            }
          });
          if (ctx)
            ctx.aliasResolveCache = nodes;
        }
        let found = void 0;
        for (const node of nodes) {
          if (node === this)
            break;
          if (node.anchor === this.source)
            found = node;
        }
        return found;
      }
      toJSON(_arg, ctx) {
        if (!ctx)
          return { source: this.source };
        const { anchors: anchors2, doc, maxAliasCount } = ctx;
        const source = this.resolve(doc, ctx);
        if (!source) {
          const msg = `Unresolved alias (the anchor must be set before the alias): ${this.source}`;
          throw new ReferenceError(msg);
        }
        let data = anchors2.get(source);
        if (!data) {
          toJS.toJS(source, null, ctx);
          data = anchors2.get(source);
        }
        if (data?.res === void 0) {
          const msg = "This should not happen: Alias anchor was not resolved?";
          throw new ReferenceError(msg);
        }
        if (maxAliasCount >= 0) {
          data.count += 1;
          if (data.aliasCount === 0)
            data.aliasCount = getAliasCount(doc, source, anchors2);
          if (data.count * data.aliasCount > maxAliasCount) {
            const msg = "Excessive alias count indicates a resource exhaustion attack";
            throw new ReferenceError(msg);
          }
        }
        return data.res;
      }
      toString(ctx, _onComment, _onChompKeep) {
        const src = `*${this.source}`;
        if (ctx) {
          anchors.anchorIsValid(this.source);
          if (ctx.options.verifyAliasOrder && !ctx.anchors.has(this.source)) {
            const msg = `Unresolved alias (the anchor must be set before the alias): ${this.source}`;
            throw new Error(msg);
          }
          if (ctx.implicitKey)
            return `${src} `;
        }
        return src;
      }
    };
    function getAliasCount(doc, node, anchors2) {
      if (identity.isAlias(node)) {
        const source = node.resolve(doc);
        const anchor = anchors2 && source && anchors2.get(source);
        return anchor ? anchor.count * anchor.aliasCount : 0;
      } else if (identity.isCollection(node)) {
        let count = 0;
        for (const item of node.items) {
          const c = getAliasCount(doc, item, anchors2);
          if (c > count)
            count = c;
        }
        return count;
      } else if (identity.isPair(node)) {
        const kc = getAliasCount(doc, node.key, anchors2);
        const vc = getAliasCount(doc, node.value, anchors2);
        return Math.max(kc, vc);
      }
      return 1;
    }
    exports.Alias = Alias;
  }
});

// node_modules/yaml/dist/nodes/Scalar.js
var require_Scalar = __commonJS({
  "node_modules/yaml/dist/nodes/Scalar.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var Node = require_Node();
    var toJS = require_toJS();
    var isScalarValue = (value2) => !value2 || typeof value2 !== "function" && typeof value2 !== "object";
    var Scalar = class extends Node.NodeBase {
      constructor(value2) {
        super(identity.SCALAR);
        this.value = value2;
      }
      toJSON(arg, ctx) {
        return ctx?.keep ? this.value : toJS.toJS(this.value, arg, ctx);
      }
      toString() {
        return String(this.value);
      }
    };
    Scalar.BLOCK_FOLDED = "BLOCK_FOLDED";
    Scalar.BLOCK_LITERAL = "BLOCK_LITERAL";
    Scalar.PLAIN = "PLAIN";
    Scalar.QUOTE_DOUBLE = "QUOTE_DOUBLE";
    Scalar.QUOTE_SINGLE = "QUOTE_SINGLE";
    exports.Scalar = Scalar;
    exports.isScalarValue = isScalarValue;
  }
});

// node_modules/yaml/dist/doc/createNode.js
var require_createNode = __commonJS({
  "node_modules/yaml/dist/doc/createNode.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Alias = require_Alias();
    var identity = require_identity();
    var Scalar = require_Scalar();
    var defaultTagPrefix = "tag:yaml.org,2002:";
    function findTagObject(value2, tagName, tags) {
      if (tagName) {
        const match = tags.filter((t) => t.tag === tagName);
        const tagObj = match.find((t) => !t.format) ?? match[0];
        if (!tagObj)
          throw new Error(`Tag ${tagName} not found`);
        return tagObj;
      }
      return tags.find((t) => t.identify?.(value2) && !t.format);
    }
    function createNode(value2, tagName, ctx) {
      if (identity.isDocument(value2))
        value2 = value2.contents;
      if (identity.isNode(value2))
        return value2;
      if (identity.isPair(value2)) {
        const map = ctx.schema[identity.MAP].createNode?.(ctx.schema, null, ctx);
        map.items.push(value2);
        return map;
      }
      if (value2 instanceof String || value2 instanceof Number || value2 instanceof Boolean || typeof BigInt !== "undefined" && value2 instanceof BigInt) {
        value2 = value2.valueOf();
      }
      const { aliasDuplicateObjects, onAnchor, onTagObj, schema: schema2, sourceObjects } = ctx;
      let ref = void 0;
      if (aliasDuplicateObjects && value2 && typeof value2 === "object") {
        ref = sourceObjects.get(value2);
        if (ref) {
          ref.anchor ?? (ref.anchor = onAnchor(value2));
          return new Alias.Alias(ref.anchor);
        } else {
          ref = { anchor: null, node: null };
          sourceObjects.set(value2, ref);
        }
      }
      if (tagName?.startsWith("!!"))
        tagName = defaultTagPrefix + tagName.slice(2);
      let tagObj = findTagObject(value2, tagName, schema2.tags);
      if (!tagObj) {
        if (value2 && typeof value2.toJSON === "function") {
          value2 = value2.toJSON();
        }
        if (!value2 || typeof value2 !== "object") {
          const node2 = new Scalar.Scalar(value2);
          if (ref)
            ref.node = node2;
          return node2;
        }
        tagObj = value2 instanceof Map ? schema2[identity.MAP] : Symbol.iterator in Object(value2) ? schema2[identity.SEQ] : schema2[identity.MAP];
      }
      if (onTagObj) {
        onTagObj(tagObj);
        delete ctx.onTagObj;
      }
      const node = tagObj?.createNode ? tagObj.createNode(ctx.schema, value2, ctx) : typeof tagObj?.nodeClass?.from === "function" ? tagObj.nodeClass.from(ctx.schema, value2, ctx) : new Scalar.Scalar(value2);
      if (tagName)
        node.tag = tagName;
      else if (!tagObj.default)
        node.tag = tagObj.tag;
      if (ref)
        ref.node = node;
      return node;
    }
    exports.createNode = createNode;
  }
});

// node_modules/yaml/dist/nodes/Collection.js
var require_Collection = __commonJS({
  "node_modules/yaml/dist/nodes/Collection.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var createNode = require_createNode();
    var identity = require_identity();
    var Node = require_Node();
    function collectionFromPath(schema2, path2, value2) {
      let v = value2;
      for (let i = path2.length - 1; i >= 0; --i) {
        const k = path2[i];
        if (typeof k === "number" && Number.isInteger(k) && k >= 0) {
          const a = [];
          a[k] = v;
          v = a;
        } else {
          v = /* @__PURE__ */ new Map([[k, v]]);
        }
      }
      return createNode.createNode(v, void 0, {
        aliasDuplicateObjects: false,
        keepUndefined: false,
        onAnchor: () => {
          throw new Error("This should not happen, please report a bug.");
        },
        schema: schema2,
        sourceObjects: /* @__PURE__ */ new Map()
      });
    }
    var isEmptyPath = (path2) => path2 == null || typeof path2 === "object" && !!path2[Symbol.iterator]().next().done;
    var Collection = class extends Node.NodeBase {
      constructor(type, schema2) {
        super(type);
        Object.defineProperty(this, "schema", {
          value: schema2,
          configurable: true,
          enumerable: false,
          writable: true
        });
      }
      /**
       * Create a copy of this collection.
       *
       * @param schema - If defined, overwrites the original's schema
       */
      clone(schema2) {
        const copy = Object.create(Object.getPrototypeOf(this), Object.getOwnPropertyDescriptors(this));
        if (schema2)
          copy.schema = schema2;
        copy.items = copy.items.map((it) => identity.isNode(it) || identity.isPair(it) ? it.clone(schema2) : it);
        if (this.range)
          copy.range = this.range.slice();
        return copy;
      }
      /**
       * Adds a value to the collection. For `!!map` and `!!omap` the value must
       * be a Pair instance or a `{ key, value }` object, which may not have a key
       * that already exists in the map.
       */
      addIn(path2, value2) {
        if (isEmptyPath(path2))
          this.add(value2);
        else {
          const [key, ...rest] = path2;
          const node = this.get(key, true);
          if (identity.isCollection(node))
            node.addIn(rest, value2);
          else if (node === void 0 && this.schema)
            this.set(key, collectionFromPath(this.schema, rest, value2));
          else
            throw new Error(`Expected YAML collection at ${key}. Remaining path: ${rest}`);
        }
      }
      /**
       * Removes a value from the collection.
       * @returns `true` if the item was found and removed.
       */
      deleteIn(path2) {
        const [key, ...rest] = path2;
        if (rest.length === 0)
          return this.delete(key);
        const node = this.get(key, true);
        if (identity.isCollection(node))
          return node.deleteIn(rest);
        else
          throw new Error(`Expected YAML collection at ${key}. Remaining path: ${rest}`);
      }
      /**
       * Returns item at `key`, or `undefined` if not found. By default unwraps
       * scalar values from their surrounding node; to disable set `keepScalar` to
       * `true` (collections are always returned intact).
       */
      getIn(path2, keepScalar) {
        const [key, ...rest] = path2;
        const node = this.get(key, true);
        if (rest.length === 0)
          return !keepScalar && identity.isScalar(node) ? node.value : node;
        else
          return identity.isCollection(node) ? node.getIn(rest, keepScalar) : void 0;
      }
      hasAllNullValues(allowScalar) {
        return this.items.every((node) => {
          if (!identity.isPair(node))
            return false;
          const n = node.value;
          return n == null || allowScalar && identity.isScalar(n) && n.value == null && !n.commentBefore && !n.comment && !n.tag;
        });
      }
      /**
       * Checks if the collection includes a value with the key `key`.
       */
      hasIn(path2) {
        const [key, ...rest] = path2;
        if (rest.length === 0)
          return this.has(key);
        const node = this.get(key, true);
        return identity.isCollection(node) ? node.hasIn(rest) : false;
      }
      /**
       * Sets a value in this collection. For `!!set`, `value` needs to be a
       * boolean to add/remove the item from the set.
       */
      setIn(path2, value2) {
        const [key, ...rest] = path2;
        if (rest.length === 0) {
          this.set(key, value2);
        } else {
          const node = this.get(key, true);
          if (identity.isCollection(node))
            node.setIn(rest, value2);
          else if (node === void 0 && this.schema)
            this.set(key, collectionFromPath(this.schema, rest, value2));
          else
            throw new Error(`Expected YAML collection at ${key}. Remaining path: ${rest}`);
        }
      }
    };
    exports.Collection = Collection;
    exports.collectionFromPath = collectionFromPath;
    exports.isEmptyPath = isEmptyPath;
  }
});

// node_modules/yaml/dist/stringify/stringifyComment.js
var require_stringifyComment = __commonJS({
  "node_modules/yaml/dist/stringify/stringifyComment.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var stringifyComment = (str) => str.replace(/^(?!$)(?: $)?/gm, "#");
    function indentComment(comment, indent) {
      if (/^\n+$/.test(comment))
        return comment.substring(1);
      return indent ? comment.replace(/^(?! *$)/gm, indent) : comment;
    }
    var lineComment = (str, indent, comment) => str.endsWith("\n") ? indentComment(comment, indent) : comment.includes("\n") ? "\n" + indentComment(comment, indent) : (str.endsWith(" ") ? "" : " ") + comment;
    exports.indentComment = indentComment;
    exports.lineComment = lineComment;
    exports.stringifyComment = stringifyComment;
  }
});

// node_modules/yaml/dist/stringify/foldFlowLines.js
var require_foldFlowLines = __commonJS({
  "node_modules/yaml/dist/stringify/foldFlowLines.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var FOLD_FLOW = "flow";
    var FOLD_BLOCK = "block";
    var FOLD_QUOTED = "quoted";
    function foldFlowLines(text, indent, mode = "flow", { indentAtStart, lineWidth = 80, minContentWidth = 20, onFold, onOverflow } = {}) {
      if (!lineWidth || lineWidth < 0)
        return text;
      if (lineWidth < minContentWidth)
        minContentWidth = 0;
      const endStep = Math.max(1 + minContentWidth, 1 + lineWidth - indent.length);
      if (text.length <= endStep)
        return text;
      const folds = [];
      const escapedFolds = {};
      let end = lineWidth - indent.length;
      if (typeof indentAtStart === "number") {
        if (indentAtStart > lineWidth - Math.max(2, minContentWidth))
          folds.push(0);
        else
          end = lineWidth - indentAtStart;
      }
      let split = void 0;
      let prev = void 0;
      let overflow = false;
      let i = -1;
      let escStart = -1;
      let escEnd = -1;
      if (mode === FOLD_BLOCK) {
        i = consumeMoreIndentedLines(text, i, indent.length);
        if (i !== -1)
          end = i + endStep;
      }
      for (let ch; ch = text[i += 1]; ) {
        if (mode === FOLD_QUOTED && ch === "\\") {
          escStart = i;
          switch (text[i + 1]) {
            case "x":
              i += 3;
              break;
            case "u":
              i += 5;
              break;
            case "U":
              i += 9;
              break;
            default:
              i += 1;
          }
          escEnd = i;
        }
        if (ch === "\n") {
          if (mode === FOLD_BLOCK)
            i = consumeMoreIndentedLines(text, i, indent.length);
          end = i + indent.length + endStep;
          split = void 0;
        } else {
          if (ch === " " && prev && prev !== " " && prev !== "\n" && prev !== "	") {
            const next = text[i + 1];
            if (next && next !== " " && next !== "\n" && next !== "	")
              split = i;
          }
          if (i >= end) {
            if (split) {
              folds.push(split);
              end = split + endStep;
              split = void 0;
            } else if (mode === FOLD_QUOTED) {
              while (prev === " " || prev === "	") {
                prev = ch;
                ch = text[i += 1];
                overflow = true;
              }
              const j = i > escEnd + 1 ? i - 2 : escStart - 1;
              if (escapedFolds[j])
                return text;
              folds.push(j);
              escapedFolds[j] = true;
              end = j + endStep;
              split = void 0;
            } else {
              overflow = true;
            }
          }
        }
        prev = ch;
      }
      if (overflow && onOverflow)
        onOverflow();
      if (folds.length === 0)
        return text;
      if (onFold)
        onFold();
      let res = text.slice(0, folds[0]);
      for (let i2 = 0; i2 < folds.length; ++i2) {
        const fold = folds[i2];
        const end2 = folds[i2 + 1] || text.length;
        if (fold === 0)
          res = `
${indent}${text.slice(0, end2)}`;
        else {
          if (mode === FOLD_QUOTED && escapedFolds[fold])
            res += `${text[fold]}\\`;
          res += `
${indent}${text.slice(fold + 1, end2)}`;
        }
      }
      return res;
    }
    function consumeMoreIndentedLines(text, i, indent) {
      let end = i;
      let start = i + 1;
      let ch = text[start];
      while (ch === " " || ch === "	") {
        if (i < start + indent) {
          ch = text[++i];
        } else {
          do {
            ch = text[++i];
          } while (ch && ch !== "\n");
          end = i;
          start = i + 1;
          ch = text[start];
        }
      }
      return end;
    }
    exports.FOLD_BLOCK = FOLD_BLOCK;
    exports.FOLD_FLOW = FOLD_FLOW;
    exports.FOLD_QUOTED = FOLD_QUOTED;
    exports.foldFlowLines = foldFlowLines;
  }
});

// node_modules/yaml/dist/stringify/stringifyString.js
var require_stringifyString = __commonJS({
  "node_modules/yaml/dist/stringify/stringifyString.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Scalar = require_Scalar();
    var foldFlowLines = require_foldFlowLines();
    var getFoldOptions = (ctx, isBlock) => ({
      indentAtStart: isBlock ? ctx.indent.length : ctx.indentAtStart,
      lineWidth: ctx.options.lineWidth,
      minContentWidth: ctx.options.minContentWidth
    });
    var containsDocumentMarker = (str) => /^(%|---|\.\.\.)/m.test(str);
    function lineLengthOverLimit(str, lineWidth, indentLength) {
      if (!lineWidth || lineWidth < 0)
        return false;
      const limit = lineWidth - indentLength;
      const strLen = str.length;
      if (strLen <= limit)
        return false;
      for (let i = 0, start = 0; i < strLen; ++i) {
        if (str[i] === "\n") {
          if (i - start > limit)
            return true;
          start = i + 1;
          if (strLen - start <= limit)
            return false;
        }
      }
      return true;
    }
    function doubleQuotedString(value2, ctx) {
      const json = JSON.stringify(value2);
      if (ctx.options.doubleQuotedAsJSON)
        return json;
      const { implicitKey } = ctx;
      const minMultiLineLength = ctx.options.doubleQuotedMinMultiLineLength;
      const indent = ctx.indent || (containsDocumentMarker(value2) ? "  " : "");
      let str = "";
      let start = 0;
      for (let i = 0, ch = json[i]; ch; ch = json[++i]) {
        if (ch === " " && json[i + 1] === "\\" && json[i + 2] === "n") {
          str += json.slice(start, i) + "\\ ";
          i += 1;
          start = i;
          ch = "\\";
        }
        if (ch === "\\")
          switch (json[i + 1]) {
            case "u":
              {
                str += json.slice(start, i);
                const code = json.substr(i + 2, 4);
                switch (code) {
                  case "0000":
                    str += "\\0";
                    break;
                  case "0007":
                    str += "\\a";
                    break;
                  case "000b":
                    str += "\\v";
                    break;
                  case "001b":
                    str += "\\e";
                    break;
                  case "0085":
                    str += "\\N";
                    break;
                  case "00a0":
                    str += "\\_";
                    break;
                  case "2028":
                    str += "\\L";
                    break;
                  case "2029":
                    str += "\\P";
                    break;
                  default:
                    if (code.substr(0, 2) === "00")
                      str += "\\x" + code.substr(2);
                    else
                      str += json.substr(i, 6);
                }
                i += 5;
                start = i + 1;
              }
              break;
            case "n":
              if (implicitKey || json[i + 2] === '"' || json.length < minMultiLineLength) {
                i += 1;
              } else {
                str += json.slice(start, i) + "\n\n";
                while (json[i + 2] === "\\" && json[i + 3] === "n" && json[i + 4] !== '"') {
                  str += "\n";
                  i += 2;
                }
                str += indent;
                if (json[i + 2] === " ")
                  str += "\\";
                i += 1;
                start = i + 1;
              }
              break;
            default:
              i += 1;
          }
      }
      str = start ? str + json.slice(start) : json;
      return implicitKey ? str : foldFlowLines.foldFlowLines(str, indent, foldFlowLines.FOLD_QUOTED, getFoldOptions(ctx, false));
    }
    function singleQuotedString(value2, ctx) {
      if (ctx.options.singleQuote === false || ctx.implicitKey && value2.includes("\n") || /[ \t]\n|\n[ \t]/.test(value2))
        return doubleQuotedString(value2, ctx);
      const indent = ctx.indent || (containsDocumentMarker(value2) ? "  " : "");
      const res = "'" + value2.replace(/'/g, "''").replace(/\n+/g, `$&
${indent}`) + "'";
      return ctx.implicitKey ? res : foldFlowLines.foldFlowLines(res, indent, foldFlowLines.FOLD_FLOW, getFoldOptions(ctx, false));
    }
    function quotedString(value2, ctx) {
      const { singleQuote } = ctx.options;
      let qs;
      if (singleQuote === false)
        qs = doubleQuotedString;
      else {
        const hasDouble = value2.includes('"');
        const hasSingle = value2.includes("'");
        if (hasDouble && !hasSingle)
          qs = singleQuotedString;
        else if (hasSingle && !hasDouble)
          qs = doubleQuotedString;
        else
          qs = singleQuote ? singleQuotedString : doubleQuotedString;
      }
      return qs(value2, ctx);
    }
    var blockEndNewlines;
    try {
      blockEndNewlines = new RegExp("(^|(?<!\n))\n+(?!\n|$)", "g");
    } catch {
      blockEndNewlines = /\n+(?!\n|$)/g;
    }
    function blockString({ comment, type, value: value2 }, ctx, onComment, onChompKeep) {
      const { blockQuote, commentString, lineWidth } = ctx.options;
      if (!blockQuote || /\n[\t ]+$/.test(value2)) {
        return quotedString(value2, ctx);
      }
      const indent = ctx.indent || (ctx.forceBlockIndent || containsDocumentMarker(value2) ? "  " : "");
      const literal = blockQuote === "literal" ? true : blockQuote === "folded" || type === Scalar.Scalar.BLOCK_FOLDED ? false : type === Scalar.Scalar.BLOCK_LITERAL ? true : !lineLengthOverLimit(value2, lineWidth, indent.length);
      if (!value2)
        return literal ? "|\n" : ">\n";
      let chomp;
      let endStart;
      for (endStart = value2.length; endStart > 0; --endStart) {
        const ch = value2[endStart - 1];
        if (ch !== "\n" && ch !== "	" && ch !== " ")
          break;
      }
      let end = value2.substring(endStart);
      const endNlPos = end.indexOf("\n");
      if (endNlPos === -1) {
        chomp = "-";
      } else if (value2 === end || endNlPos !== end.length - 1) {
        chomp = "+";
        if (onChompKeep)
          onChompKeep();
      } else {
        chomp = "";
      }
      if (end) {
        value2 = value2.slice(0, -end.length);
        if (end[end.length - 1] === "\n")
          end = end.slice(0, -1);
        end = end.replace(blockEndNewlines, `$&${indent}`);
      }
      let startWithSpace = false;
      let startEnd;
      let startNlPos = -1;
      for (startEnd = 0; startEnd < value2.length; ++startEnd) {
        const ch = value2[startEnd];
        if (ch === " ")
          startWithSpace = true;
        else if (ch === "\n")
          startNlPos = startEnd;
        else
          break;
      }
      let start = value2.substring(0, startNlPos < startEnd ? startNlPos + 1 : startEnd);
      if (start) {
        value2 = value2.substring(start.length);
        start = start.replace(/\n+/g, `$&${indent}`);
      }
      const indentSize = indent ? "2" : "1";
      let header = (startWithSpace ? indentSize : "") + chomp;
      if (comment) {
        header += " " + commentString(comment.replace(/ ?[\r\n]+/g, " "));
        if (onComment)
          onComment();
      }
      if (!literal) {
        const foldedValue = value2.replace(/\n+/g, "\n$&").replace(/(?:^|\n)([\t ].*)(?:([\n\t ]*)\n(?![\n\t ]))?/g, "$1$2").replace(/\n+/g, `$&${indent}`);
        let literalFallback = false;
        const foldOptions = getFoldOptions(ctx, true);
        if (blockQuote !== "folded" && type !== Scalar.Scalar.BLOCK_FOLDED) {
          foldOptions.onOverflow = () => {
            literalFallback = true;
          };
        }
        const body = foldFlowLines.foldFlowLines(`${start}${foldedValue}${end}`, indent, foldFlowLines.FOLD_BLOCK, foldOptions);
        if (!literalFallback)
          return `>${header}
${indent}${body}`;
      }
      value2 = value2.replace(/\n+/g, `$&${indent}`);
      return `|${header}
${indent}${start}${value2}${end}`;
    }
    function plainString(item, ctx, onComment, onChompKeep) {
      const { type, value: value2 } = item;
      const { actualString, implicitKey, indent, indentStep, inFlow } = ctx;
      if (implicitKey && value2.includes("\n") || inFlow && /[[\]{},]/.test(value2)) {
        return quotedString(value2, ctx);
      }
      if (/^[\n\t ,[\]{}#&*!|>'"%@`]|^[?-]$|^[?-][ \t]|[\n:][ \t]|[ \t]\n|[\n\t ]#|[\n\t :]$/.test(value2)) {
        return implicitKey || inFlow || !value2.includes("\n") ? quotedString(value2, ctx) : blockString(item, ctx, onComment, onChompKeep);
      }
      if (!implicitKey && !inFlow && type !== Scalar.Scalar.PLAIN && value2.includes("\n")) {
        return blockString(item, ctx, onComment, onChompKeep);
      }
      if (containsDocumentMarker(value2)) {
        if (indent === "") {
          ctx.forceBlockIndent = true;
          return blockString(item, ctx, onComment, onChompKeep);
        } else if (implicitKey && indent === indentStep) {
          return quotedString(value2, ctx);
        }
      }
      const str = value2.replace(/\n+/g, `$&
${indent}`);
      if (actualString) {
        const test = (tag) => tag.default && tag.tag !== "tag:yaml.org,2002:str" && tag.test?.test(str);
        const { compat, tags } = ctx.doc.schema;
        if (tags.some(test) || compat?.some(test))
          return quotedString(value2, ctx);
      }
      return implicitKey ? str : foldFlowLines.foldFlowLines(str, indent, foldFlowLines.FOLD_FLOW, getFoldOptions(ctx, false));
    }
    function stringifyString(item, ctx, onComment, onChompKeep) {
      const { implicitKey, inFlow } = ctx;
      const ss = typeof item.value === "string" ? item : Object.assign({}, item, { value: String(item.value) });
      let { type } = item;
      if (type !== Scalar.Scalar.QUOTE_DOUBLE) {
        if (/[\x00-\x08\x0b-\x1f\x7f-\x9f\u{D800}-\u{DFFF}]/u.test(ss.value))
          type = Scalar.Scalar.QUOTE_DOUBLE;
      }
      const _stringify = (_type) => {
        switch (_type) {
          case Scalar.Scalar.BLOCK_FOLDED:
          case Scalar.Scalar.BLOCK_LITERAL:
            return implicitKey || inFlow ? quotedString(ss.value, ctx) : blockString(ss, ctx, onComment, onChompKeep);
          case Scalar.Scalar.QUOTE_DOUBLE:
            return doubleQuotedString(ss.value, ctx);
          case Scalar.Scalar.QUOTE_SINGLE:
            return singleQuotedString(ss.value, ctx);
          case Scalar.Scalar.PLAIN:
            return plainString(ss, ctx, onComment, onChompKeep);
          default:
            return null;
        }
      };
      let res = _stringify(type);
      if (res === null) {
        const { defaultKeyType, defaultStringType } = ctx.options;
        const t = implicitKey && defaultKeyType || defaultStringType;
        res = _stringify(t);
        if (res === null)
          throw new Error(`Unsupported default string type ${t}`);
      }
      return res;
    }
    exports.stringifyString = stringifyString;
  }
});

// node_modules/yaml/dist/stringify/stringify.js
var require_stringify = __commonJS({
  "node_modules/yaml/dist/stringify/stringify.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var anchors = require_anchors();
    var identity = require_identity();
    var stringifyComment = require_stringifyComment();
    var stringifyString = require_stringifyString();
    function createStringifyContext(doc, options) {
      const opt = Object.assign({
        blockQuote: true,
        commentString: stringifyComment.stringifyComment,
        defaultKeyType: null,
        defaultStringType: "PLAIN",
        directives: null,
        doubleQuotedAsJSON: false,
        doubleQuotedMinMultiLineLength: 40,
        falseStr: "false",
        flowCollectionPadding: true,
        indentSeq: true,
        lineWidth: 80,
        minContentWidth: 20,
        nullStr: "null",
        simpleKeys: false,
        singleQuote: null,
        trailingComma: false,
        trueStr: "true",
        verifyAliasOrder: true
      }, doc.schema.toStringOptions, options);
      let inFlow;
      switch (opt.collectionStyle) {
        case "block":
          inFlow = false;
          break;
        case "flow":
          inFlow = true;
          break;
        default:
          inFlow = null;
      }
      return {
        anchors: /* @__PURE__ */ new Set(),
        doc,
        flowCollectionPadding: opt.flowCollectionPadding ? " " : "",
        indent: "",
        indentStep: typeof opt.indent === "number" ? " ".repeat(opt.indent) : "  ",
        inFlow,
        options: opt
      };
    }
    function getTagObject(tags, item) {
      if (item.tag) {
        const match = tags.filter((t) => t.tag === item.tag);
        if (match.length > 0)
          return match.find((t) => t.format === item.format) ?? match[0];
      }
      let tagObj = void 0;
      let obj;
      if (identity.isScalar(item)) {
        obj = item.value;
        let match = tags.filter((t) => t.identify?.(obj));
        if (match.length > 1) {
          const testMatch = match.filter((t) => t.test);
          if (testMatch.length > 0)
            match = testMatch;
        }
        tagObj = match.find((t) => t.format === item.format) ?? match.find((t) => !t.format);
      } else {
        obj = item;
        tagObj = tags.find((t) => t.nodeClass && obj instanceof t.nodeClass);
      }
      if (!tagObj) {
        const name = obj?.constructor?.name ?? (obj === null ? "null" : typeof obj);
        throw new Error(`Tag not resolved for ${name} value`);
      }
      return tagObj;
    }
    function stringifyProps(node, tagObj, { anchors: anchors$1, doc }) {
      if (!doc.directives)
        return "";
      const props = [];
      const anchor = (identity.isScalar(node) || identity.isCollection(node)) && node.anchor;
      if (anchor && anchors.anchorIsValid(anchor)) {
        anchors$1.add(anchor);
        props.push(`&${anchor}`);
      }
      const tag = node.tag ?? (tagObj.default ? null : tagObj.tag);
      if (tag)
        props.push(doc.directives.tagString(tag));
      return props.join(" ");
    }
    function stringify(item, ctx, onComment, onChompKeep) {
      if (identity.isPair(item))
        return item.toString(ctx, onComment, onChompKeep);
      if (identity.isAlias(item)) {
        if (ctx.doc.directives)
          return item.toString(ctx);
        if (ctx.resolvedAliases?.has(item)) {
          throw new TypeError(`Cannot stringify circular structure without alias nodes`);
        } else {
          if (ctx.resolvedAliases)
            ctx.resolvedAliases.add(item);
          else
            ctx.resolvedAliases = /* @__PURE__ */ new Set([item]);
          item = item.resolve(ctx.doc);
        }
      }
      let tagObj = void 0;
      const node = identity.isNode(item) ? item : ctx.doc.createNode(item, { onTagObj: (o) => tagObj = o });
      tagObj ?? (tagObj = getTagObject(ctx.doc.schema.tags, node));
      const props = stringifyProps(node, tagObj, ctx);
      if (props.length > 0)
        ctx.indentAtStart = (ctx.indentAtStart ?? 0) + props.length + 1;
      const str = typeof tagObj.stringify === "function" ? tagObj.stringify(node, ctx, onComment, onChompKeep) : identity.isScalar(node) ? stringifyString.stringifyString(node, ctx, onComment, onChompKeep) : node.toString(ctx, onComment, onChompKeep);
      if (!props)
        return str;
      return identity.isScalar(node) || str[0] === "{" || str[0] === "[" ? `${props} ${str}` : `${props}
${ctx.indent}${str}`;
    }
    exports.createStringifyContext = createStringifyContext;
    exports.stringify = stringify;
  }
});

// node_modules/yaml/dist/stringify/stringifyPair.js
var require_stringifyPair = __commonJS({
  "node_modules/yaml/dist/stringify/stringifyPair.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var Scalar = require_Scalar();
    var stringify = require_stringify();
    var stringifyComment = require_stringifyComment();
    function stringifyPair({ key, value: value2 }, ctx, onComment, onChompKeep) {
      const { allNullValues, doc, indent, indentStep, options: { commentString, indentSeq, simpleKeys } } = ctx;
      let keyComment = identity.isNode(key) && key.comment || null;
      if (simpleKeys) {
        if (keyComment) {
          throw new Error("With simple keys, key nodes cannot have comments");
        }
        if (identity.isCollection(key) || !identity.isNode(key) && typeof key === "object") {
          const msg = "With simple keys, collection cannot be used as a key value";
          throw new Error(msg);
        }
      }
      let explicitKey = !simpleKeys && (!key || keyComment && value2 == null && !ctx.inFlow || identity.isCollection(key) || (identity.isScalar(key) ? key.type === Scalar.Scalar.BLOCK_FOLDED || key.type === Scalar.Scalar.BLOCK_LITERAL : typeof key === "object"));
      ctx = Object.assign({}, ctx, {
        allNullValues: false,
        implicitKey: !explicitKey && (simpleKeys || !allNullValues),
        indent: indent + indentStep
      });
      let keyCommentDone = false;
      let chompKeep = false;
      let str = stringify.stringify(key, ctx, () => keyCommentDone = true, () => chompKeep = true);
      if (!explicitKey && !ctx.inFlow && str.length > 1024) {
        if (simpleKeys)
          throw new Error("With simple keys, single line scalar must not span more than 1024 characters");
        explicitKey = true;
      }
      if (ctx.inFlow) {
        if (allNullValues || value2 == null) {
          if (keyCommentDone && onComment)
            onComment();
          return str === "" ? "?" : explicitKey ? `? ${str}` : str;
        }
      } else if (allNullValues && !simpleKeys || value2 == null && explicitKey) {
        str = `? ${str}`;
        if (keyComment && !keyCommentDone) {
          str += stringifyComment.lineComment(str, ctx.indent, commentString(keyComment));
        } else if (chompKeep && onChompKeep)
          onChompKeep();
        return str;
      }
      if (keyCommentDone)
        keyComment = null;
      if (explicitKey) {
        if (keyComment)
          str += stringifyComment.lineComment(str, ctx.indent, commentString(keyComment));
        str = `? ${str}
${indent}:`;
      } else {
        str = `${str}:`;
        if (keyComment)
          str += stringifyComment.lineComment(str, ctx.indent, commentString(keyComment));
      }
      let vsb, vcb, valueComment;
      if (identity.isNode(value2)) {
        vsb = !!value2.spaceBefore;
        vcb = value2.commentBefore;
        valueComment = value2.comment;
      } else {
        vsb = false;
        vcb = null;
        valueComment = null;
        if (value2 && typeof value2 === "object")
          value2 = doc.createNode(value2);
      }
      ctx.implicitKey = false;
      if (!explicitKey && !keyComment && identity.isScalar(value2))
        ctx.indentAtStart = str.length + 1;
      chompKeep = false;
      if (!indentSeq && indentStep.length >= 2 && !ctx.inFlow && !explicitKey && identity.isSeq(value2) && !value2.flow && !value2.tag && !value2.anchor) {
        ctx.indent = ctx.indent.substring(2);
      }
      let valueCommentDone = false;
      const valueStr = stringify.stringify(value2, ctx, () => valueCommentDone = true, () => chompKeep = true);
      let ws = " ";
      if (keyComment || vsb || vcb) {
        ws = vsb ? "\n" : "";
        if (vcb) {
          const cs = commentString(vcb);
          ws += `
${stringifyComment.indentComment(cs, ctx.indent)}`;
        }
        if (valueStr === "" && !ctx.inFlow) {
          if (ws === "\n" && valueComment)
            ws = "\n\n";
        } else {
          ws += `
${ctx.indent}`;
        }
      } else if (!explicitKey && identity.isCollection(value2)) {
        const vs0 = valueStr[0];
        const nl0 = valueStr.indexOf("\n");
        const hasNewline = nl0 !== -1;
        const flow = ctx.inFlow ?? value2.flow ?? value2.items.length === 0;
        if (hasNewline || !flow) {
          let hasPropsLine = false;
          if (hasNewline && (vs0 === "&" || vs0 === "!")) {
            let sp0 = valueStr.indexOf(" ");
            if (vs0 === "&" && sp0 !== -1 && sp0 < nl0 && valueStr[sp0 + 1] === "!") {
              sp0 = valueStr.indexOf(" ", sp0 + 1);
            }
            if (sp0 === -1 || nl0 < sp0)
              hasPropsLine = true;
          }
          if (!hasPropsLine)
            ws = `
${ctx.indent}`;
        }
      } else if (valueStr === "" || valueStr[0] === "\n") {
        ws = "";
      }
      str += ws + valueStr;
      if (ctx.inFlow) {
        if (valueCommentDone && onComment)
          onComment();
      } else if (valueComment && !valueCommentDone) {
        str += stringifyComment.lineComment(str, ctx.indent, commentString(valueComment));
      } else if (chompKeep && onChompKeep) {
        onChompKeep();
      }
      return str;
    }
    exports.stringifyPair = stringifyPair;
  }
});

// node_modules/yaml/dist/log.js
var require_log = __commonJS({
  "node_modules/yaml/dist/log.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var node_process = __require("process");
    function debug(logLevel, ...messages) {
      if (logLevel === "debug")
        console.log(...messages);
    }
    function warn(logLevel, warning) {
      if (logLevel === "debug" || logLevel === "warn") {
        if (typeof node_process.emitWarning === "function")
          node_process.emitWarning(warning);
        else
          console.warn(warning);
      }
    }
    exports.debug = debug;
    exports.warn = warn;
  }
});

// node_modules/yaml/dist/schema/yaml-1.1/merge.js
var require_merge = __commonJS({
  "node_modules/yaml/dist/schema/yaml-1.1/merge.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var Scalar = require_Scalar();
    var MERGE_KEY = "<<";
    var merge = {
      identify: (value2) => value2 === MERGE_KEY || typeof value2 === "symbol" && value2.description === MERGE_KEY,
      default: "key",
      tag: "tag:yaml.org,2002:merge",
      test: /^<<$/,
      resolve: () => Object.assign(new Scalar.Scalar(Symbol(MERGE_KEY)), {
        addToJSMap: addMergeToJSMap
      }),
      stringify: () => MERGE_KEY
    };
    var isMergeKey = (ctx, key) => (merge.identify(key) || identity.isScalar(key) && (!key.type || key.type === Scalar.Scalar.PLAIN) && merge.identify(key.value)) && ctx?.doc.schema.tags.some((tag) => tag.tag === merge.tag && tag.default);
    function addMergeToJSMap(ctx, map, value2) {
      const source = resolveAliasValue(ctx, value2);
      if (identity.isSeq(source))
        for (const it of source.items)
          mergeValue(ctx, map, it);
      else if (Array.isArray(source))
        for (const it of source)
          mergeValue(ctx, map, it);
      else
        mergeValue(ctx, map, source);
    }
    function mergeValue(ctx, map, value2) {
      const source = resolveAliasValue(ctx, value2);
      if (!identity.isMap(source))
        throw new Error("Merge sources must be maps or map aliases");
      const srcMap = source.toJSON(null, ctx, Map);
      for (const [key, value3] of srcMap) {
        if (map instanceof Map) {
          if (!map.has(key))
            map.set(key, value3);
        } else if (map instanceof Set) {
          map.add(key);
        } else if (!Object.prototype.hasOwnProperty.call(map, key)) {
          Object.defineProperty(map, key, {
            value: value3,
            writable: true,
            enumerable: true,
            configurable: true
          });
        }
      }
      return map;
    }
    function resolveAliasValue(ctx, value2) {
      return ctx && identity.isAlias(value2) ? value2.resolve(ctx.doc, ctx) : value2;
    }
    exports.addMergeToJSMap = addMergeToJSMap;
    exports.isMergeKey = isMergeKey;
    exports.merge = merge;
  }
});

// node_modules/yaml/dist/nodes/addPairToJSMap.js
var require_addPairToJSMap = __commonJS({
  "node_modules/yaml/dist/nodes/addPairToJSMap.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var log = require_log();
    var merge = require_merge();
    var stringify = require_stringify();
    var identity = require_identity();
    var toJS = require_toJS();
    function addPairToJSMap(ctx, map, { key, value: value2 }) {
      if (identity.isNode(key) && key.addToJSMap)
        key.addToJSMap(ctx, map, value2);
      else if (merge.isMergeKey(ctx, key))
        merge.addMergeToJSMap(ctx, map, value2);
      else {
        const jsKey = toJS.toJS(key, "", ctx);
        if (map instanceof Map) {
          map.set(jsKey, toJS.toJS(value2, jsKey, ctx));
        } else if (map instanceof Set) {
          map.add(jsKey);
        } else {
          const stringKey = stringifyKey(key, jsKey, ctx);
          const jsValue = toJS.toJS(value2, stringKey, ctx);
          if (stringKey in map)
            Object.defineProperty(map, stringKey, {
              value: jsValue,
              writable: true,
              enumerable: true,
              configurable: true
            });
          else
            map[stringKey] = jsValue;
        }
      }
      return map;
    }
    function stringifyKey(key, jsKey, ctx) {
      if (jsKey === null)
        return "";
      if (typeof jsKey !== "object")
        return String(jsKey);
      if (identity.isNode(key) && ctx?.doc) {
        const strCtx = stringify.createStringifyContext(ctx.doc, {});
        strCtx.anchors = /* @__PURE__ */ new Set();
        for (const node of ctx.anchors.keys())
          strCtx.anchors.add(node.anchor);
        strCtx.inFlow = true;
        strCtx.inStringifyKey = true;
        const strKey = key.toString(strCtx);
        if (!ctx.mapKeyWarned) {
          let jsonStr = JSON.stringify(strKey);
          if (jsonStr.length > 40)
            jsonStr = jsonStr.substring(0, 36) + '..."';
          log.warn(ctx.doc.options.logLevel, `Keys with collection values will be stringified due to JS Object restrictions: ${jsonStr}. Set mapAsMap: true to use object keys.`);
          ctx.mapKeyWarned = true;
        }
        return strKey;
      }
      return JSON.stringify(jsKey);
    }
    exports.addPairToJSMap = addPairToJSMap;
  }
});

// node_modules/yaml/dist/nodes/Pair.js
var require_Pair = __commonJS({
  "node_modules/yaml/dist/nodes/Pair.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var createNode = require_createNode();
    var stringifyPair = require_stringifyPair();
    var addPairToJSMap = require_addPairToJSMap();
    var identity = require_identity();
    function createPair(key, value2, ctx) {
      const k = createNode.createNode(key, void 0, ctx);
      const v = createNode.createNode(value2, void 0, ctx);
      return new Pair(k, v);
    }
    var Pair = class _Pair {
      constructor(key, value2 = null) {
        Object.defineProperty(this, identity.NODE_TYPE, { value: identity.PAIR });
        this.key = key;
        this.value = value2;
      }
      clone(schema2) {
        let { key, value: value2 } = this;
        if (identity.isNode(key))
          key = key.clone(schema2);
        if (identity.isNode(value2))
          value2 = value2.clone(schema2);
        return new _Pair(key, value2);
      }
      toJSON(_, ctx) {
        const pair = ctx?.mapAsMap ? /* @__PURE__ */ new Map() : {};
        return addPairToJSMap.addPairToJSMap(ctx, pair, this);
      }
      toString(ctx, onComment, onChompKeep) {
        return ctx?.doc ? stringifyPair.stringifyPair(this, ctx, onComment, onChompKeep) : JSON.stringify(this);
      }
    };
    exports.Pair = Pair;
    exports.createPair = createPair;
  }
});

// node_modules/yaml/dist/stringify/stringifyCollection.js
var require_stringifyCollection = __commonJS({
  "node_modules/yaml/dist/stringify/stringifyCollection.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var stringify = require_stringify();
    var stringifyComment = require_stringifyComment();
    function stringifyCollection(collection, ctx, options) {
      const flow = ctx.inFlow ?? collection.flow;
      const stringify2 = flow ? stringifyFlowCollection : stringifyBlockCollection;
      return stringify2(collection, ctx, options);
    }
    function stringifyBlockCollection({ comment, items }, ctx, { blockItemPrefix, flowChars, itemIndent, onChompKeep, onComment }) {
      const { indent, options: { commentString } } = ctx;
      const itemCtx = Object.assign({}, ctx, { indent: itemIndent, type: null });
      let chompKeep = false;
      const lines = [];
      for (let i = 0; i < items.length; ++i) {
        const item = items[i];
        let comment2 = null;
        if (identity.isNode(item)) {
          if (!chompKeep && item.spaceBefore)
            lines.push("");
          addCommentBefore(ctx, lines, item.commentBefore, chompKeep);
          if (item.comment)
            comment2 = item.comment;
        } else if (identity.isPair(item)) {
          const ik = identity.isNode(item.key) ? item.key : null;
          if (ik) {
            if (!chompKeep && ik.spaceBefore)
              lines.push("");
            addCommentBefore(ctx, lines, ik.commentBefore, chompKeep);
          }
        }
        chompKeep = false;
        let str2 = stringify.stringify(item, itemCtx, () => comment2 = null, () => chompKeep = true);
        if (comment2)
          str2 += stringifyComment.lineComment(str2, itemIndent, commentString(comment2));
        if (chompKeep && comment2)
          chompKeep = false;
        lines.push(blockItemPrefix + str2);
      }
      let str;
      if (lines.length === 0) {
        str = flowChars.start + flowChars.end;
      } else {
        str = lines[0];
        for (let i = 1; i < lines.length; ++i) {
          const line = lines[i];
          str += line ? `
${indent}${line}` : "\n";
        }
      }
      if (comment) {
        str += "\n" + stringifyComment.indentComment(commentString(comment), indent);
        if (onComment)
          onComment();
      } else if (chompKeep && onChompKeep)
        onChompKeep();
      return str;
    }
    function stringifyFlowCollection({ items }, ctx, { flowChars, itemIndent }) {
      const { indent, indentStep, flowCollectionPadding: fcPadding, options: { commentString } } = ctx;
      itemIndent += indentStep;
      const itemCtx = Object.assign({}, ctx, {
        indent: itemIndent,
        inFlow: true,
        type: null
      });
      let reqNewline = false;
      let linesAtValue = 0;
      const lines = [];
      for (let i = 0; i < items.length; ++i) {
        const item = items[i];
        let comment = null;
        if (identity.isNode(item)) {
          if (item.spaceBefore)
            lines.push("");
          addCommentBefore(ctx, lines, item.commentBefore, false);
          if (item.comment)
            comment = item.comment;
        } else if (identity.isPair(item)) {
          const ik = identity.isNode(item.key) ? item.key : null;
          if (ik) {
            if (ik.spaceBefore)
              lines.push("");
            addCommentBefore(ctx, lines, ik.commentBefore, false);
            if (ik.comment)
              reqNewline = true;
          }
          const iv = identity.isNode(item.value) ? item.value : null;
          if (iv) {
            if (iv.comment)
              comment = iv.comment;
            if (iv.commentBefore)
              reqNewline = true;
          } else if (item.value == null && ik?.comment) {
            comment = ik.comment;
          }
        }
        if (comment)
          reqNewline = true;
        let str = stringify.stringify(item, itemCtx, () => comment = null);
        reqNewline || (reqNewline = lines.length > linesAtValue || str.includes("\n"));
        if (i < items.length - 1) {
          str += ",";
        } else if (ctx.options.trailingComma) {
          if (ctx.options.lineWidth > 0) {
            reqNewline || (reqNewline = lines.reduce((sum, line) => sum + line.length + 2, 2) + (str.length + 2) > ctx.options.lineWidth);
          }
          if (reqNewline) {
            str += ",";
          }
        }
        if (comment)
          str += stringifyComment.lineComment(str, itemIndent, commentString(comment));
        lines.push(str);
        linesAtValue = lines.length;
      }
      const { start, end } = flowChars;
      if (lines.length === 0) {
        return start + end;
      } else {
        if (!reqNewline) {
          const len = lines.reduce((sum, line) => sum + line.length + 2, 2);
          reqNewline = ctx.options.lineWidth > 0 && len > ctx.options.lineWidth;
        }
        if (reqNewline) {
          let str = start;
          for (const line of lines)
            str += line ? `
${indentStep}${indent}${line}` : "\n";
          return `${str}
${indent}${end}`;
        } else {
          return `${start}${fcPadding}${lines.join(" ")}${fcPadding}${end}`;
        }
      }
    }
    function addCommentBefore({ indent, options: { commentString } }, lines, comment, chompKeep) {
      if (comment && chompKeep)
        comment = comment.replace(/^\n+/, "");
      if (comment) {
        const ic = stringifyComment.indentComment(commentString(comment), indent);
        lines.push(ic.trimStart());
      }
    }
    exports.stringifyCollection = stringifyCollection;
  }
});

// node_modules/yaml/dist/nodes/YAMLMap.js
var require_YAMLMap = __commonJS({
  "node_modules/yaml/dist/nodes/YAMLMap.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var stringifyCollection = require_stringifyCollection();
    var addPairToJSMap = require_addPairToJSMap();
    var Collection = require_Collection();
    var identity = require_identity();
    var Pair = require_Pair();
    var Scalar = require_Scalar();
    function findPair(items, key) {
      const k = identity.isScalar(key) ? key.value : key;
      for (const it of items) {
        if (identity.isPair(it)) {
          if (it.key === key || it.key === k)
            return it;
          if (identity.isScalar(it.key) && it.key.value === k)
            return it;
        }
      }
      return void 0;
    }
    var YAMLMap = class extends Collection.Collection {
      static get tagName() {
        return "tag:yaml.org,2002:map";
      }
      constructor(schema2) {
        super(identity.MAP, schema2);
        this.items = [];
      }
      /**
       * A generic collection parsing method that can be extended
       * to other node classes that inherit from YAMLMap
       */
      static from(schema2, obj, ctx) {
        const { keepUndefined, replacer } = ctx;
        const map = new this(schema2);
        const add = (key, value2) => {
          if (typeof replacer === "function")
            value2 = replacer.call(obj, key, value2);
          else if (Array.isArray(replacer) && !replacer.includes(key))
            return;
          if (value2 !== void 0 || keepUndefined)
            map.items.push(Pair.createPair(key, value2, ctx));
        };
        if (obj instanceof Map) {
          for (const [key, value2] of obj)
            add(key, value2);
        } else if (obj && typeof obj === "object") {
          for (const key of Object.keys(obj))
            add(key, obj[key]);
        }
        if (typeof schema2.sortMapEntries === "function") {
          map.items.sort(schema2.sortMapEntries);
        }
        return map;
      }
      /**
       * Adds a value to the collection.
       *
       * @param overwrite - If not set `true`, using a key that is already in the
       *   collection will throw. Otherwise, overwrites the previous value.
       */
      add(pair, overwrite) {
        let _pair;
        if (identity.isPair(pair))
          _pair = pair;
        else if (!pair || typeof pair !== "object" || !("key" in pair)) {
          _pair = new Pair.Pair(pair, pair?.value);
        } else
          _pair = new Pair.Pair(pair.key, pair.value);
        const prev = findPair(this.items, _pair.key);
        const sortEntries = this.schema?.sortMapEntries;
        if (prev) {
          if (!overwrite)
            throw new Error(`Key ${_pair.key} already set`);
          if (identity.isScalar(prev.value) && Scalar.isScalarValue(_pair.value))
            prev.value.value = _pair.value;
          else
            prev.value = _pair.value;
        } else if (sortEntries) {
          const i = this.items.findIndex((item) => sortEntries(_pair, item) < 0);
          if (i === -1)
            this.items.push(_pair);
          else
            this.items.splice(i, 0, _pair);
        } else {
          this.items.push(_pair);
        }
      }
      delete(key) {
        const it = findPair(this.items, key);
        if (!it)
          return false;
        const del = this.items.splice(this.items.indexOf(it), 1);
        return del.length > 0;
      }
      get(key, keepScalar) {
        const it = findPair(this.items, key);
        const node = it?.value;
        return (!keepScalar && identity.isScalar(node) ? node.value : node) ?? void 0;
      }
      has(key) {
        return !!findPair(this.items, key);
      }
      set(key, value2) {
        this.add(new Pair.Pair(key, value2), true);
      }
      /**
       * @param ctx - Conversion context, originally set in Document#toJS()
       * @param {Class} Type - If set, forces the returned collection type
       * @returns Instance of Type, Map, or Object
       */
      toJSON(_, ctx, Type) {
        const map = Type ? new Type() : ctx?.mapAsMap ? /* @__PURE__ */ new Map() : {};
        if (ctx?.onCreate)
          ctx.onCreate(map);
        for (const item of this.items)
          addPairToJSMap.addPairToJSMap(ctx, map, item);
        return map;
      }
      toString(ctx, onComment, onChompKeep) {
        if (!ctx)
          return JSON.stringify(this);
        for (const item of this.items) {
          if (!identity.isPair(item))
            throw new Error(`Map items must all be pairs; found ${JSON.stringify(item)} instead`);
        }
        if (!ctx.allNullValues && this.hasAllNullValues(false))
          ctx = Object.assign({}, ctx, { allNullValues: true });
        return stringifyCollection.stringifyCollection(this, ctx, {
          blockItemPrefix: "",
          flowChars: { start: "{", end: "}" },
          itemIndent: ctx.indent || "",
          onChompKeep,
          onComment
        });
      }
    };
    exports.YAMLMap = YAMLMap;
    exports.findPair = findPair;
  }
});

// node_modules/yaml/dist/schema/common/map.js
var require_map = __commonJS({
  "node_modules/yaml/dist/schema/common/map.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var YAMLMap = require_YAMLMap();
    var map = {
      collection: "map",
      default: true,
      nodeClass: YAMLMap.YAMLMap,
      tag: "tag:yaml.org,2002:map",
      resolve(map2, onError) {
        if (!identity.isMap(map2))
          onError("Expected a mapping for this tag");
        return map2;
      },
      createNode: (schema2, obj, ctx) => YAMLMap.YAMLMap.from(schema2, obj, ctx)
    };
    exports.map = map;
  }
});

// node_modules/yaml/dist/nodes/YAMLSeq.js
var require_YAMLSeq = __commonJS({
  "node_modules/yaml/dist/nodes/YAMLSeq.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var createNode = require_createNode();
    var stringifyCollection = require_stringifyCollection();
    var Collection = require_Collection();
    var identity = require_identity();
    var Scalar = require_Scalar();
    var toJS = require_toJS();
    var YAMLSeq = class extends Collection.Collection {
      static get tagName() {
        return "tag:yaml.org,2002:seq";
      }
      constructor(schema2) {
        super(identity.SEQ, schema2);
        this.items = [];
      }
      add(value2) {
        this.items.push(value2);
      }
      /**
       * Removes a value from the collection.
       *
       * `key` must contain a representation of an integer for this to succeed.
       * It may be wrapped in a `Scalar`.
       *
       * @returns `true` if the item was found and removed.
       */
      delete(key) {
        const idx = asItemIndex(key);
        if (typeof idx !== "number")
          return false;
        const del = this.items.splice(idx, 1);
        return del.length > 0;
      }
      get(key, keepScalar) {
        const idx = asItemIndex(key);
        if (typeof idx !== "number")
          return void 0;
        const it = this.items[idx];
        return !keepScalar && identity.isScalar(it) ? it.value : it;
      }
      /**
       * Checks if the collection includes a value with the key `key`.
       *
       * `key` must contain a representation of an integer for this to succeed.
       * It may be wrapped in a `Scalar`.
       */
      has(key) {
        const idx = asItemIndex(key);
        return typeof idx === "number" && idx < this.items.length;
      }
      /**
       * Sets a value in this collection. For `!!set`, `value` needs to be a
       * boolean to add/remove the item from the set.
       *
       * If `key` does not contain a representation of an integer, this will throw.
       * It may be wrapped in a `Scalar`.
       */
      set(key, value2) {
        const idx = asItemIndex(key);
        if (typeof idx !== "number")
          throw new Error(`Expected a valid index, not ${key}.`);
        const prev = this.items[idx];
        if (identity.isScalar(prev) && Scalar.isScalarValue(value2))
          prev.value = value2;
        else
          this.items[idx] = value2;
      }
      toJSON(_, ctx) {
        const seq = [];
        if (ctx?.onCreate)
          ctx.onCreate(seq);
        let i = 0;
        for (const item of this.items)
          seq.push(toJS.toJS(item, String(i++), ctx));
        return seq;
      }
      toString(ctx, onComment, onChompKeep) {
        if (!ctx)
          return JSON.stringify(this);
        return stringifyCollection.stringifyCollection(this, ctx, {
          blockItemPrefix: "- ",
          flowChars: { start: "[", end: "]" },
          itemIndent: (ctx.indent || "") + "  ",
          onChompKeep,
          onComment
        });
      }
      static from(schema2, obj, ctx) {
        const { replacer } = ctx;
        const seq = new this(schema2);
        if (obj && Symbol.iterator in Object(obj)) {
          let i = 0;
          for (let it of obj) {
            if (typeof replacer === "function") {
              const key = obj instanceof Set ? it : String(i++);
              it = replacer.call(obj, key, it);
            }
            seq.items.push(createNode.createNode(it, void 0, ctx));
          }
        }
        return seq;
      }
    };
    function asItemIndex(key) {
      let idx = identity.isScalar(key) ? key.value : key;
      if (idx && typeof idx === "string")
        idx = Number(idx);
      return typeof idx === "number" && Number.isInteger(idx) && idx >= 0 ? idx : null;
    }
    exports.YAMLSeq = YAMLSeq;
  }
});

// node_modules/yaml/dist/schema/common/seq.js
var require_seq = __commonJS({
  "node_modules/yaml/dist/schema/common/seq.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var YAMLSeq = require_YAMLSeq();
    var seq = {
      collection: "seq",
      default: true,
      nodeClass: YAMLSeq.YAMLSeq,
      tag: "tag:yaml.org,2002:seq",
      resolve(seq2, onError) {
        if (!identity.isSeq(seq2))
          onError("Expected a sequence for this tag");
        return seq2;
      },
      createNode: (schema2, obj, ctx) => YAMLSeq.YAMLSeq.from(schema2, obj, ctx)
    };
    exports.seq = seq;
  }
});

// node_modules/yaml/dist/schema/common/string.js
var require_string = __commonJS({
  "node_modules/yaml/dist/schema/common/string.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var stringifyString = require_stringifyString();
    var string = {
      identify: (value2) => typeof value2 === "string",
      default: true,
      tag: "tag:yaml.org,2002:str",
      resolve: (str) => str,
      stringify(item, ctx, onComment, onChompKeep) {
        ctx = Object.assign({ actualString: true }, ctx);
        return stringifyString.stringifyString(item, ctx, onComment, onChompKeep);
      }
    };
    exports.string = string;
  }
});

// node_modules/yaml/dist/schema/common/null.js
var require_null = __commonJS({
  "node_modules/yaml/dist/schema/common/null.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Scalar = require_Scalar();
    var nullTag = {
      identify: (value2) => value2 == null,
      createNode: () => new Scalar.Scalar(null),
      default: true,
      tag: "tag:yaml.org,2002:null",
      test: /^(?:~|[Nn]ull|NULL)?$/,
      resolve: () => new Scalar.Scalar(null),
      stringify: ({ source }, ctx) => typeof source === "string" && nullTag.test.test(source) ? source : ctx.options.nullStr
    };
    exports.nullTag = nullTag;
  }
});

// node_modules/yaml/dist/schema/core/bool.js
var require_bool = __commonJS({
  "node_modules/yaml/dist/schema/core/bool.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Scalar = require_Scalar();
    var boolTag = {
      identify: (value2) => typeof value2 === "boolean",
      default: true,
      tag: "tag:yaml.org,2002:bool",
      test: /^(?:[Tt]rue|TRUE|[Ff]alse|FALSE)$/,
      resolve: (str) => new Scalar.Scalar(str[0] === "t" || str[0] === "T"),
      stringify({ source, value: value2 }, ctx) {
        if (source && boolTag.test.test(source)) {
          const sv = source[0] === "t" || source[0] === "T";
          if (value2 === sv)
            return source;
        }
        return value2 ? ctx.options.trueStr : ctx.options.falseStr;
      }
    };
    exports.boolTag = boolTag;
  }
});

// node_modules/yaml/dist/stringify/stringifyNumber.js
var require_stringifyNumber = __commonJS({
  "node_modules/yaml/dist/stringify/stringifyNumber.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    function stringifyNumber({ format, minFractionDigits, tag, value: value2 }) {
      if (typeof value2 === "bigint")
        return String(value2);
      const num = typeof value2 === "number" ? value2 : Number(value2);
      if (!isFinite(num))
        return isNaN(num) ? ".nan" : num < 0 ? "-.inf" : ".inf";
      let n = Object.is(value2, -0) ? "-0" : JSON.stringify(value2);
      if (!format && minFractionDigits && (!tag || tag === "tag:yaml.org,2002:float") && /^-?\d/.test(n) && !n.includes("e")) {
        let i = n.indexOf(".");
        if (i < 0) {
          i = n.length;
          n += ".";
        }
        let d = minFractionDigits - (n.length - i - 1);
        while (d-- > 0)
          n += "0";
      }
      return n;
    }
    exports.stringifyNumber = stringifyNumber;
  }
});

// node_modules/yaml/dist/schema/core/float.js
var require_float = __commonJS({
  "node_modules/yaml/dist/schema/core/float.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Scalar = require_Scalar();
    var stringifyNumber = require_stringifyNumber();
    var floatNaN = {
      identify: (value2) => typeof value2 === "number",
      default: true,
      tag: "tag:yaml.org,2002:float",
      test: /^(?:[-+]?\.(?:inf|Inf|INF)|\.nan|\.NaN|\.NAN)$/,
      resolve: (str) => str.slice(-3).toLowerCase() === "nan" ? NaN : str[0] === "-" ? Number.NEGATIVE_INFINITY : Number.POSITIVE_INFINITY,
      stringify: stringifyNumber.stringifyNumber
    };
    var floatExp = {
      identify: (value2) => typeof value2 === "number",
      default: true,
      tag: "tag:yaml.org,2002:float",
      format: "EXP",
      test: /^[-+]?(?:\.[0-9]+|[0-9]+(?:\.[0-9]*)?)[eE][-+]?[0-9]+$/,
      resolve: (str) => parseFloat(str),
      stringify(node) {
        const num = Number(node.value);
        return isFinite(num) ? num.toExponential() : stringifyNumber.stringifyNumber(node);
      }
    };
    var float = {
      identify: (value2) => typeof value2 === "number",
      default: true,
      tag: "tag:yaml.org,2002:float",
      test: /^[-+]?(?:\.[0-9]+|[0-9]+\.[0-9]*)$/,
      resolve(str) {
        const node = new Scalar.Scalar(parseFloat(str));
        const dot = str.indexOf(".");
        if (dot !== -1 && str[str.length - 1] === "0")
          node.minFractionDigits = str.length - dot - 1;
        return node;
      },
      stringify: stringifyNumber.stringifyNumber
    };
    exports.float = float;
    exports.floatExp = floatExp;
    exports.floatNaN = floatNaN;
  }
});

// node_modules/yaml/dist/schema/core/int.js
var require_int = __commonJS({
  "node_modules/yaml/dist/schema/core/int.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var stringifyNumber = require_stringifyNumber();
    var intIdentify = (value2) => typeof value2 === "bigint" || Number.isInteger(value2);
    var intResolve = (str, offset, radix, { intAsBigInt }) => intAsBigInt ? BigInt(str) : parseInt(str.substring(offset), radix);
    function intStringify(node, radix, prefix) {
      const { value: value2 } = node;
      if (intIdentify(value2) && value2 >= 0)
        return prefix + value2.toString(radix);
      return stringifyNumber.stringifyNumber(node);
    }
    var intOct = {
      identify: (value2) => intIdentify(value2) && value2 >= 0,
      default: true,
      tag: "tag:yaml.org,2002:int",
      format: "OCT",
      test: /^0o[0-7]+$/,
      resolve: (str, _onError, opt) => intResolve(str, 2, 8, opt),
      stringify: (node) => intStringify(node, 8, "0o")
    };
    var int = {
      identify: intIdentify,
      default: true,
      tag: "tag:yaml.org,2002:int",
      test: /^[-+]?[0-9]+$/,
      resolve: (str, _onError, opt) => intResolve(str, 0, 10, opt),
      stringify: stringifyNumber.stringifyNumber
    };
    var intHex = {
      identify: (value2) => intIdentify(value2) && value2 >= 0,
      default: true,
      tag: "tag:yaml.org,2002:int",
      format: "HEX",
      test: /^0x[0-9a-fA-F]+$/,
      resolve: (str, _onError, opt) => intResolve(str, 2, 16, opt),
      stringify: (node) => intStringify(node, 16, "0x")
    };
    exports.int = int;
    exports.intHex = intHex;
    exports.intOct = intOct;
  }
});

// node_modules/yaml/dist/schema/core/schema.js
var require_schema2 = __commonJS({
  "node_modules/yaml/dist/schema/core/schema.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var map = require_map();
    var _null = require_null();
    var seq = require_seq();
    var string = require_string();
    var bool = require_bool();
    var float = require_float();
    var int = require_int();
    var schema2 = [
      map.map,
      seq.seq,
      string.string,
      _null.nullTag,
      bool.boolTag,
      int.intOct,
      int.int,
      int.intHex,
      float.floatNaN,
      float.floatExp,
      float.float
    ];
    exports.schema = schema2;
  }
});

// node_modules/yaml/dist/schema/json/schema.js
var require_schema3 = __commonJS({
  "node_modules/yaml/dist/schema/json/schema.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Scalar = require_Scalar();
    var map = require_map();
    var seq = require_seq();
    function intIdentify(value2) {
      return typeof value2 === "bigint" || Number.isInteger(value2);
    }
    var stringifyJSON = ({ value: value2 }) => JSON.stringify(value2);
    var jsonScalars = [
      {
        identify: (value2) => typeof value2 === "string",
        default: true,
        tag: "tag:yaml.org,2002:str",
        resolve: (str) => str,
        stringify: stringifyJSON
      },
      {
        identify: (value2) => value2 == null,
        createNode: () => new Scalar.Scalar(null),
        default: true,
        tag: "tag:yaml.org,2002:null",
        test: /^null$/,
        resolve: () => null,
        stringify: stringifyJSON
      },
      {
        identify: (value2) => typeof value2 === "boolean",
        default: true,
        tag: "tag:yaml.org,2002:bool",
        test: /^true$|^false$/,
        resolve: (str) => str === "true",
        stringify: stringifyJSON
      },
      {
        identify: intIdentify,
        default: true,
        tag: "tag:yaml.org,2002:int",
        test: /^-?(?:0|[1-9][0-9]*)$/,
        resolve: (str, _onError, { intAsBigInt }) => intAsBigInt ? BigInt(str) : parseInt(str, 10),
        stringify: ({ value: value2 }) => intIdentify(value2) ? value2.toString() : JSON.stringify(value2)
      },
      {
        identify: (value2) => typeof value2 === "number",
        default: true,
        tag: "tag:yaml.org,2002:float",
        test: /^-?(?:0|[1-9][0-9]*)(?:\.[0-9]*)?(?:[eE][-+]?[0-9]+)?$/,
        resolve: (str) => parseFloat(str),
        stringify: stringifyJSON
      }
    ];
    var jsonError = {
      default: true,
      tag: "",
      test: /^/,
      resolve(str, onError) {
        onError(`Unresolved plain scalar ${JSON.stringify(str)}`);
        return str;
      }
    };
    var schema2 = [map.map, seq.seq].concat(jsonScalars, jsonError);
    exports.schema = schema2;
  }
});

// node_modules/yaml/dist/schema/yaml-1.1/binary.js
var require_binary = __commonJS({
  "node_modules/yaml/dist/schema/yaml-1.1/binary.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var node_buffer = __require("buffer");
    var Scalar = require_Scalar();
    var stringifyString = require_stringifyString();
    var binary = {
      identify: (value2) => value2 instanceof Uint8Array,
      // Buffer inherits from Uint8Array
      default: false,
      tag: "tag:yaml.org,2002:binary",
      /**
       * Returns a Buffer in node and an Uint8Array in browsers
       *
       * To use the resulting buffer as an image, you'll want to do something like:
       *
       *   const blob = new Blob([buffer], { type: 'image/jpeg' })
       *   document.querySelector('#photo').src = URL.createObjectURL(blob)
       */
      resolve(src, onError) {
        if (typeof node_buffer.Buffer === "function") {
          return node_buffer.Buffer.from(src, "base64");
        } else if (typeof atob === "function") {
          const str = atob(src.replace(/[\n\r]/g, ""));
          const buffer = new Uint8Array(str.length);
          for (let i = 0; i < str.length; ++i)
            buffer[i] = str.charCodeAt(i);
          return buffer;
        } else {
          onError("This environment does not support reading binary tags; either Buffer or atob is required");
          return src;
        }
      },
      stringify({ comment, type, value: value2 }, ctx, onComment, onChompKeep) {
        if (!value2)
          return "";
        const buf = value2;
        let str;
        if (typeof node_buffer.Buffer === "function") {
          str = buf instanceof node_buffer.Buffer ? buf.toString("base64") : node_buffer.Buffer.from(buf.buffer).toString("base64");
        } else if (typeof btoa === "function") {
          let s = "";
          for (let i = 0; i < buf.length; ++i)
            s += String.fromCharCode(buf[i]);
          str = btoa(s);
        } else {
          throw new Error("This environment does not support writing binary tags; either Buffer or btoa is required");
        }
        type ?? (type = Scalar.Scalar.BLOCK_LITERAL);
        if (type !== Scalar.Scalar.QUOTE_DOUBLE) {
          const lineWidth = Math.max(ctx.options.lineWidth - ctx.indent.length, ctx.options.minContentWidth);
          const n = Math.ceil(str.length / lineWidth);
          const lines = new Array(n);
          for (let i = 0, o = 0; i < n; ++i, o += lineWidth) {
            lines[i] = str.substr(o, lineWidth);
          }
          str = lines.join(type === Scalar.Scalar.BLOCK_LITERAL ? "\n" : " ");
        }
        return stringifyString.stringifyString({ comment, type, value: str }, ctx, onComment, onChompKeep);
      }
    };
    exports.binary = binary;
  }
});

// node_modules/yaml/dist/schema/yaml-1.1/pairs.js
var require_pairs = __commonJS({
  "node_modules/yaml/dist/schema/yaml-1.1/pairs.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var Pair = require_Pair();
    var Scalar = require_Scalar();
    var YAMLSeq = require_YAMLSeq();
    function resolvePairs(seq, onError) {
      if (identity.isSeq(seq)) {
        for (let i = 0; i < seq.items.length; ++i) {
          let item = seq.items[i];
          if (identity.isPair(item))
            continue;
          else if (identity.isMap(item)) {
            if (item.items.length > 1)
              onError("Each pair must have its own sequence indicator");
            const pair = item.items[0] || new Pair.Pair(new Scalar.Scalar(null));
            if (item.commentBefore)
              pair.key.commentBefore = pair.key.commentBefore ? `${item.commentBefore}
${pair.key.commentBefore}` : item.commentBefore;
            if (item.comment) {
              const cn = pair.value ?? pair.key;
              cn.comment = cn.comment ? `${item.comment}
${cn.comment}` : item.comment;
            }
            item = pair;
          }
          seq.items[i] = identity.isPair(item) ? item : new Pair.Pair(item);
        }
      } else
        onError("Expected a sequence for this tag");
      return seq;
    }
    function createPairs(schema2, iterable, ctx) {
      const { replacer } = ctx;
      const pairs2 = new YAMLSeq.YAMLSeq(schema2);
      pairs2.tag = "tag:yaml.org,2002:pairs";
      let i = 0;
      if (iterable && Symbol.iterator in Object(iterable))
        for (let it of iterable) {
          if (typeof replacer === "function")
            it = replacer.call(iterable, String(i++), it);
          let key, value2;
          if (Array.isArray(it)) {
            if (it.length === 2) {
              key = it[0];
              value2 = it[1];
            } else
              throw new TypeError(`Expected [key, value] tuple: ${it}`);
          } else if (it && it instanceof Object) {
            const keys = Object.keys(it);
            if (keys.length === 1) {
              key = keys[0];
              value2 = it[key];
            } else {
              throw new TypeError(`Expected tuple with one key, not ${keys.length} keys`);
            }
          } else {
            key = it;
          }
          pairs2.items.push(Pair.createPair(key, value2, ctx));
        }
      return pairs2;
    }
    var pairs = {
      collection: "seq",
      default: false,
      tag: "tag:yaml.org,2002:pairs",
      resolve: resolvePairs,
      createNode: createPairs
    };
    exports.createPairs = createPairs;
    exports.pairs = pairs;
    exports.resolvePairs = resolvePairs;
  }
});

// node_modules/yaml/dist/schema/yaml-1.1/omap.js
var require_omap = __commonJS({
  "node_modules/yaml/dist/schema/yaml-1.1/omap.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var toJS = require_toJS();
    var YAMLMap = require_YAMLMap();
    var YAMLSeq = require_YAMLSeq();
    var pairs = require_pairs();
    var YAMLOMap = class _YAMLOMap extends YAMLSeq.YAMLSeq {
      constructor() {
        super();
        this.add = YAMLMap.YAMLMap.prototype.add.bind(this);
        this.delete = YAMLMap.YAMLMap.prototype.delete.bind(this);
        this.get = YAMLMap.YAMLMap.prototype.get.bind(this);
        this.has = YAMLMap.YAMLMap.prototype.has.bind(this);
        this.set = YAMLMap.YAMLMap.prototype.set.bind(this);
        this.tag = _YAMLOMap.tag;
      }
      /**
       * If `ctx` is given, the return type is actually `Map<unknown, unknown>`,
       * but TypeScript won't allow widening the signature of a child method.
       */
      toJSON(_, ctx) {
        if (!ctx)
          return super.toJSON(_);
        const map = /* @__PURE__ */ new Map();
        if (ctx?.onCreate)
          ctx.onCreate(map);
        for (const pair of this.items) {
          let key, value2;
          if (identity.isPair(pair)) {
            key = toJS.toJS(pair.key, "", ctx);
            value2 = toJS.toJS(pair.value, key, ctx);
          } else {
            key = toJS.toJS(pair, "", ctx);
          }
          if (map.has(key))
            throw new Error("Ordered maps must not include duplicate keys");
          map.set(key, value2);
        }
        return map;
      }
      static from(schema2, iterable, ctx) {
        const pairs$1 = pairs.createPairs(schema2, iterable, ctx);
        const omap2 = new this();
        omap2.items = pairs$1.items;
        return omap2;
      }
    };
    YAMLOMap.tag = "tag:yaml.org,2002:omap";
    var omap = {
      collection: "seq",
      identify: (value2) => value2 instanceof Map,
      nodeClass: YAMLOMap,
      default: false,
      tag: "tag:yaml.org,2002:omap",
      resolve(seq, onError) {
        const pairs$1 = pairs.resolvePairs(seq, onError);
        const seenKeys = [];
        for (const { key } of pairs$1.items) {
          if (identity.isScalar(key)) {
            if (seenKeys.includes(key.value)) {
              onError(`Ordered maps must not include duplicate keys: ${key.value}`);
            } else {
              seenKeys.push(key.value);
            }
          }
        }
        return Object.assign(new YAMLOMap(), pairs$1);
      },
      createNode: (schema2, iterable, ctx) => YAMLOMap.from(schema2, iterable, ctx)
    };
    exports.YAMLOMap = YAMLOMap;
    exports.omap = omap;
  }
});

// node_modules/yaml/dist/schema/yaml-1.1/bool.js
var require_bool2 = __commonJS({
  "node_modules/yaml/dist/schema/yaml-1.1/bool.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Scalar = require_Scalar();
    function boolStringify({ value: value2, source }, ctx) {
      const boolObj = value2 ? trueTag : falseTag;
      if (source && boolObj.test.test(source))
        return source;
      return value2 ? ctx.options.trueStr : ctx.options.falseStr;
    }
    var trueTag = {
      identify: (value2) => value2 === true,
      default: true,
      tag: "tag:yaml.org,2002:bool",
      test: /^(?:Y|y|[Yy]es|YES|[Tt]rue|TRUE|[Oo]n|ON)$/,
      resolve: () => new Scalar.Scalar(true),
      stringify: boolStringify
    };
    var falseTag = {
      identify: (value2) => value2 === false,
      default: true,
      tag: "tag:yaml.org,2002:bool",
      test: /^(?:N|n|[Nn]o|NO|[Ff]alse|FALSE|[Oo]ff|OFF)$/,
      resolve: () => new Scalar.Scalar(false),
      stringify: boolStringify
    };
    exports.falseTag = falseTag;
    exports.trueTag = trueTag;
  }
});

// node_modules/yaml/dist/schema/yaml-1.1/float.js
var require_float2 = __commonJS({
  "node_modules/yaml/dist/schema/yaml-1.1/float.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Scalar = require_Scalar();
    var stringifyNumber = require_stringifyNumber();
    var floatNaN = {
      identify: (value2) => typeof value2 === "number",
      default: true,
      tag: "tag:yaml.org,2002:float",
      test: /^(?:[-+]?\.(?:inf|Inf|INF)|\.nan|\.NaN|\.NAN)$/,
      resolve: (str) => str.slice(-3).toLowerCase() === "nan" ? NaN : str[0] === "-" ? Number.NEGATIVE_INFINITY : Number.POSITIVE_INFINITY,
      stringify: stringifyNumber.stringifyNumber
    };
    var floatExp = {
      identify: (value2) => typeof value2 === "number",
      default: true,
      tag: "tag:yaml.org,2002:float",
      format: "EXP",
      test: /^[-+]?(?:[0-9][0-9_]*)?(?:\.[0-9_]*)?[eE][-+]?[0-9]+$/,
      resolve: (str) => parseFloat(str.replace(/_/g, "")),
      stringify(node) {
        const num = Number(node.value);
        return isFinite(num) ? num.toExponential() : stringifyNumber.stringifyNumber(node);
      }
    };
    var float = {
      identify: (value2) => typeof value2 === "number",
      default: true,
      tag: "tag:yaml.org,2002:float",
      test: /^[-+]?(?:[0-9][0-9_]*)?\.[0-9_]*$/,
      resolve(str) {
        const node = new Scalar.Scalar(parseFloat(str.replace(/_/g, "")));
        const dot = str.indexOf(".");
        if (dot !== -1) {
          const f = str.substring(dot + 1).replace(/_/g, "");
          if (f[f.length - 1] === "0")
            node.minFractionDigits = f.length;
        }
        return node;
      },
      stringify: stringifyNumber.stringifyNumber
    };
    exports.float = float;
    exports.floatExp = floatExp;
    exports.floatNaN = floatNaN;
  }
});

// node_modules/yaml/dist/schema/yaml-1.1/int.js
var require_int2 = __commonJS({
  "node_modules/yaml/dist/schema/yaml-1.1/int.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var stringifyNumber = require_stringifyNumber();
    var intIdentify = (value2) => typeof value2 === "bigint" || Number.isInteger(value2);
    function intResolve(str, offset, radix, { intAsBigInt }) {
      const sign = str[0];
      if (sign === "-" || sign === "+")
        offset += 1;
      str = str.substring(offset).replace(/_/g, "");
      if (intAsBigInt) {
        switch (radix) {
          case 2:
            str = `0b${str}`;
            break;
          case 8:
            str = `0o${str}`;
            break;
          case 16:
            str = `0x${str}`;
            break;
        }
        const n2 = BigInt(str);
        return sign === "-" ? BigInt(-1) * n2 : n2;
      }
      const n = parseInt(str, radix);
      return sign === "-" ? -1 * n : n;
    }
    function intStringify(node, radix, prefix) {
      const { value: value2 } = node;
      if (intIdentify(value2)) {
        const str = value2.toString(radix);
        return value2 < 0 ? "-" + prefix + str.substr(1) : prefix + str;
      }
      return stringifyNumber.stringifyNumber(node);
    }
    var intBin = {
      identify: intIdentify,
      default: true,
      tag: "tag:yaml.org,2002:int",
      format: "BIN",
      test: /^[-+]?0b[0-1_]+$/,
      resolve: (str, _onError, opt) => intResolve(str, 2, 2, opt),
      stringify: (node) => intStringify(node, 2, "0b")
    };
    var intOct = {
      identify: intIdentify,
      default: true,
      tag: "tag:yaml.org,2002:int",
      format: "OCT",
      test: /^[-+]?0[0-7_]+$/,
      resolve: (str, _onError, opt) => intResolve(str, 1, 8, opt),
      stringify: (node) => intStringify(node, 8, "0")
    };
    var int = {
      identify: intIdentify,
      default: true,
      tag: "tag:yaml.org,2002:int",
      test: /^[-+]?[0-9][0-9_]*$/,
      resolve: (str, _onError, opt) => intResolve(str, 0, 10, opt),
      stringify: stringifyNumber.stringifyNumber
    };
    var intHex = {
      identify: intIdentify,
      default: true,
      tag: "tag:yaml.org,2002:int",
      format: "HEX",
      test: /^[-+]?0x[0-9a-fA-F_]+$/,
      resolve: (str, _onError, opt) => intResolve(str, 2, 16, opt),
      stringify: (node) => intStringify(node, 16, "0x")
    };
    exports.int = int;
    exports.intBin = intBin;
    exports.intHex = intHex;
    exports.intOct = intOct;
  }
});

// node_modules/yaml/dist/schema/yaml-1.1/set.js
var require_set = __commonJS({
  "node_modules/yaml/dist/schema/yaml-1.1/set.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var Pair = require_Pair();
    var YAMLMap = require_YAMLMap();
    var YAMLSet = class _YAMLSet extends YAMLMap.YAMLMap {
      constructor(schema2) {
        super(schema2);
        this.tag = _YAMLSet.tag;
      }
      add(key) {
        let pair;
        if (identity.isPair(key))
          pair = key;
        else if (key && typeof key === "object" && "key" in key && "value" in key && key.value === null)
          pair = new Pair.Pair(key.key, null);
        else
          pair = new Pair.Pair(key, null);
        const prev = YAMLMap.findPair(this.items, pair.key);
        if (!prev)
          this.items.push(pair);
      }
      /**
       * If `keepPair` is `true`, returns the Pair matching `key`.
       * Otherwise, returns the value of that Pair's key.
       */
      get(key, keepPair) {
        const pair = YAMLMap.findPair(this.items, key);
        return !keepPair && identity.isPair(pair) ? identity.isScalar(pair.key) ? pair.key.value : pair.key : pair;
      }
      set(key, value2) {
        if (typeof value2 !== "boolean")
          throw new Error(`Expected boolean value for set(key, value) in a YAML set, not ${typeof value2}`);
        const prev = YAMLMap.findPair(this.items, key);
        if (prev && !value2) {
          this.items.splice(this.items.indexOf(prev), 1);
        } else if (!prev && value2) {
          this.items.push(new Pair.Pair(key));
        }
      }
      toJSON(_, ctx) {
        return super.toJSON(_, ctx, Set);
      }
      toString(ctx, onComment, onChompKeep) {
        if (!ctx)
          return JSON.stringify(this);
        if (this.hasAllNullValues(true))
          return super.toString(Object.assign({}, ctx, { allNullValues: true }), onComment, onChompKeep);
        else
          throw new Error("Set items must all have null values");
      }
      static from(schema2, iterable, ctx) {
        const { replacer } = ctx;
        const set2 = new this(schema2);
        if (iterable && Symbol.iterator in Object(iterable))
          for (let value2 of iterable) {
            if (typeof replacer === "function")
              value2 = replacer.call(iterable, value2, value2);
            set2.items.push(Pair.createPair(value2, null, ctx));
          }
        return set2;
      }
    };
    YAMLSet.tag = "tag:yaml.org,2002:set";
    var set = {
      collection: "map",
      identify: (value2) => value2 instanceof Set,
      nodeClass: YAMLSet,
      default: false,
      tag: "tag:yaml.org,2002:set",
      createNode: (schema2, iterable, ctx) => YAMLSet.from(schema2, iterable, ctx),
      resolve(map, onError) {
        if (identity.isMap(map)) {
          if (map.hasAllNullValues(true))
            return Object.assign(new YAMLSet(), map);
          else
            onError("Set items must all have null values");
        } else
          onError("Expected a mapping for this tag");
        return map;
      }
    };
    exports.YAMLSet = YAMLSet;
    exports.set = set;
  }
});

// node_modules/yaml/dist/schema/yaml-1.1/timestamp.js
var require_timestamp = __commonJS({
  "node_modules/yaml/dist/schema/yaml-1.1/timestamp.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var stringifyNumber = require_stringifyNumber();
    function parseSexagesimal(str, asBigInt) {
      const sign = str[0];
      const parts = sign === "-" || sign === "+" ? str.substring(1) : str;
      const num = (n) => asBigInt ? BigInt(n) : Number(n);
      const res = parts.replace(/_/g, "").split(":").reduce((res2, p) => res2 * num(60) + num(p), num(0));
      return sign === "-" ? num(-1) * res : res;
    }
    function stringifySexagesimal(node) {
      let { value: value2 } = node;
      let num = (n) => n;
      if (typeof value2 === "bigint")
        num = (n) => BigInt(n);
      else if (isNaN(value2) || !isFinite(value2))
        return stringifyNumber.stringifyNumber(node);
      let sign = "";
      if (value2 < 0) {
        sign = "-";
        value2 *= num(-1);
      }
      const _60 = num(60);
      const parts = [value2 % _60];
      if (value2 < 60) {
        parts.unshift(0);
      } else {
        value2 = (value2 - parts[0]) / _60;
        parts.unshift(value2 % _60);
        if (value2 >= 60) {
          value2 = (value2 - parts[0]) / _60;
          parts.unshift(value2);
        }
      }
      return sign + parts.map((n) => String(n).padStart(2, "0")).join(":").replace(/000000\d*$/, "");
    }
    var intTime = {
      identify: (value2) => typeof value2 === "bigint" || Number.isInteger(value2),
      default: true,
      tag: "tag:yaml.org,2002:int",
      format: "TIME",
      test: /^[-+]?[0-9][0-9_]*(?::[0-5]?[0-9])+$/,
      resolve: (str, _onError, { intAsBigInt }) => parseSexagesimal(str, intAsBigInt),
      stringify: stringifySexagesimal
    };
    var floatTime = {
      identify: (value2) => typeof value2 === "number",
      default: true,
      tag: "tag:yaml.org,2002:float",
      format: "TIME",
      test: /^[-+]?[0-9][0-9_]*(?::[0-5]?[0-9])+\.[0-9_]*$/,
      resolve: (str) => parseSexagesimal(str, false),
      stringify: stringifySexagesimal
    };
    var timestamp = {
      identify: (value2) => value2 instanceof Date,
      default: true,
      tag: "tag:yaml.org,2002:timestamp",
      // If the time zone is omitted, the timestamp is assumed to be specified in UTC. The time part
      // may be omitted altogether, resulting in a date format. In such a case, the time part is
      // assumed to be 00:00:00Z (start of day, UTC).
      test: RegExp("^([0-9]{4})-([0-9]{1,2})-([0-9]{1,2})(?:(?:t|T|[ \\t]+)([0-9]{1,2}):([0-9]{1,2}):([0-9]{1,2}(\\.[0-9]+)?)(?:[ \\t]*(Z|[-+][012]?[0-9](?::[0-9]{2})?))?)?$"),
      resolve(str) {
        const match = str.match(timestamp.test);
        if (!match)
          throw new Error("!!timestamp expects a date, starting with yyyy-mm-dd");
        const [, year, month, day, hour, minute, second] = match.map(Number);
        const millisec = match[7] ? Number((match[7] + "00").substr(1, 3)) : 0;
        let date = Date.UTC(year, month - 1, day, hour || 0, minute || 0, second || 0, millisec);
        const tz = match[8];
        if (tz && tz !== "Z") {
          let d = parseSexagesimal(tz, false);
          if (Math.abs(d) < 30)
            d *= 60;
          date -= 6e4 * d;
        }
        return new Date(date);
      },
      stringify: ({ value: value2 }) => value2?.toISOString().replace(/(T00:00:00)?\.000Z$/, "") ?? ""
    };
    exports.floatTime = floatTime;
    exports.intTime = intTime;
    exports.timestamp = timestamp;
  }
});

// node_modules/yaml/dist/schema/yaml-1.1/schema.js
var require_schema4 = __commonJS({
  "node_modules/yaml/dist/schema/yaml-1.1/schema.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var map = require_map();
    var _null = require_null();
    var seq = require_seq();
    var string = require_string();
    var binary = require_binary();
    var bool = require_bool2();
    var float = require_float2();
    var int = require_int2();
    var merge = require_merge();
    var omap = require_omap();
    var pairs = require_pairs();
    var set = require_set();
    var timestamp = require_timestamp();
    var schema2 = [
      map.map,
      seq.seq,
      string.string,
      _null.nullTag,
      bool.trueTag,
      bool.falseTag,
      int.intBin,
      int.intOct,
      int.int,
      int.intHex,
      float.floatNaN,
      float.floatExp,
      float.float,
      binary.binary,
      merge.merge,
      omap.omap,
      pairs.pairs,
      set.set,
      timestamp.intTime,
      timestamp.floatTime,
      timestamp.timestamp
    ];
    exports.schema = schema2;
  }
});

// node_modules/yaml/dist/schema/tags.js
var require_tags = __commonJS({
  "node_modules/yaml/dist/schema/tags.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var map = require_map();
    var _null = require_null();
    var seq = require_seq();
    var string = require_string();
    var bool = require_bool();
    var float = require_float();
    var int = require_int();
    var schema2 = require_schema2();
    var schema$1 = require_schema3();
    var binary = require_binary();
    var merge = require_merge();
    var omap = require_omap();
    var pairs = require_pairs();
    var schema$2 = require_schema4();
    var set = require_set();
    var timestamp = require_timestamp();
    var schemas = /* @__PURE__ */ new Map([
      ["core", schema2.schema],
      ["failsafe", [map.map, seq.seq, string.string]],
      ["json", schema$1.schema],
      ["yaml11", schema$2.schema],
      ["yaml-1.1", schema$2.schema]
    ]);
    var tagsByName = {
      binary: binary.binary,
      bool: bool.boolTag,
      float: float.float,
      floatExp: float.floatExp,
      floatNaN: float.floatNaN,
      floatTime: timestamp.floatTime,
      int: int.int,
      intHex: int.intHex,
      intOct: int.intOct,
      intTime: timestamp.intTime,
      map: map.map,
      merge: merge.merge,
      null: _null.nullTag,
      omap: omap.omap,
      pairs: pairs.pairs,
      seq: seq.seq,
      set: set.set,
      timestamp: timestamp.timestamp
    };
    var coreKnownTags = {
      "tag:yaml.org,2002:binary": binary.binary,
      "tag:yaml.org,2002:merge": merge.merge,
      "tag:yaml.org,2002:omap": omap.omap,
      "tag:yaml.org,2002:pairs": pairs.pairs,
      "tag:yaml.org,2002:set": set.set,
      "tag:yaml.org,2002:timestamp": timestamp.timestamp
    };
    function getTags(customTags, schemaName, addMergeTag) {
      const schemaTags = schemas.get(schemaName);
      if (schemaTags && !customTags) {
        return addMergeTag && !schemaTags.includes(merge.merge) ? schemaTags.concat(merge.merge) : schemaTags.slice();
      }
      let tags = schemaTags;
      if (!tags) {
        if (Array.isArray(customTags))
          tags = [];
        else {
          const keys = Array.from(schemas.keys()).filter((key) => key !== "yaml11").map((key) => JSON.stringify(key)).join(", ");
          throw new Error(`Unknown schema "${schemaName}"; use one of ${keys} or define customTags array`);
        }
      }
      if (Array.isArray(customTags)) {
        for (const tag of customTags)
          tags = tags.concat(tag);
      } else if (typeof customTags === "function") {
        tags = customTags(tags.slice());
      }
      if (addMergeTag)
        tags = tags.concat(merge.merge);
      return tags.reduce((tags2, tag) => {
        const tagObj = typeof tag === "string" ? tagsByName[tag] : tag;
        if (!tagObj) {
          const tagName = JSON.stringify(tag);
          const keys = Object.keys(tagsByName).map((key) => JSON.stringify(key)).join(", ");
          throw new Error(`Unknown custom tag ${tagName}; use one of ${keys}`);
        }
        if (!tags2.includes(tagObj))
          tags2.push(tagObj);
        return tags2;
      }, []);
    }
    exports.coreKnownTags = coreKnownTags;
    exports.getTags = getTags;
  }
});

// node_modules/yaml/dist/schema/Schema.js
var require_Schema = __commonJS({
  "node_modules/yaml/dist/schema/Schema.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var map = require_map();
    var seq = require_seq();
    var string = require_string();
    var tags = require_tags();
    var sortMapEntriesByKey = (a, b) => a.key < b.key ? -1 : a.key > b.key ? 1 : 0;
    var Schema = class _Schema {
      constructor({ compat, customTags, merge, resolveKnownTags, schema: schema2, sortMapEntries, toStringDefaults }) {
        this.compat = Array.isArray(compat) ? tags.getTags(compat, "compat") : compat ? tags.getTags(null, compat) : null;
        this.name = typeof schema2 === "string" && schema2 || "core";
        this.knownTags = resolveKnownTags ? tags.coreKnownTags : {};
        this.tags = tags.getTags(customTags, this.name, merge);
        this.toStringOptions = toStringDefaults ?? null;
        Object.defineProperty(this, identity.MAP, { value: map.map });
        Object.defineProperty(this, identity.SCALAR, { value: string.string });
        Object.defineProperty(this, identity.SEQ, { value: seq.seq });
        this.sortMapEntries = typeof sortMapEntries === "function" ? sortMapEntries : sortMapEntries === true ? sortMapEntriesByKey : null;
      }
      clone() {
        const copy = Object.create(_Schema.prototype, Object.getOwnPropertyDescriptors(this));
        copy.tags = this.tags.slice();
        return copy;
      }
    };
    exports.Schema = Schema;
  }
});

// node_modules/yaml/dist/stringify/stringifyDocument.js
var require_stringifyDocument = __commonJS({
  "node_modules/yaml/dist/stringify/stringifyDocument.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var stringify = require_stringify();
    var stringifyComment = require_stringifyComment();
    function stringifyDocument(doc, options) {
      const lines = [];
      let hasDirectives = options.directives === true;
      if (options.directives !== false && doc.directives) {
        const dir = doc.directives.toString(doc);
        if (dir) {
          lines.push(dir);
          hasDirectives = true;
        } else if (doc.directives.docStart)
          hasDirectives = true;
      }
      if (hasDirectives)
        lines.push("---");
      const ctx = stringify.createStringifyContext(doc, options);
      const { commentString } = ctx.options;
      if (doc.commentBefore) {
        if (lines.length !== 1)
          lines.unshift("");
        const cs = commentString(doc.commentBefore);
        lines.unshift(stringifyComment.indentComment(cs, ""));
      }
      let chompKeep = false;
      let contentComment = null;
      if (doc.contents) {
        if (identity.isNode(doc.contents)) {
          if (doc.contents.spaceBefore && hasDirectives)
            lines.push("");
          if (doc.contents.commentBefore) {
            const cs = commentString(doc.contents.commentBefore);
            lines.push(stringifyComment.indentComment(cs, ""));
          }
          ctx.forceBlockIndent = !!doc.comment;
          contentComment = doc.contents.comment;
        }
        const onChompKeep = contentComment ? void 0 : () => chompKeep = true;
        let body = stringify.stringify(doc.contents, ctx, () => contentComment = null, onChompKeep);
        if (contentComment)
          body += stringifyComment.lineComment(body, "", commentString(contentComment));
        if ((body[0] === "|" || body[0] === ">") && lines[lines.length - 1] === "---") {
          lines[lines.length - 1] = `--- ${body}`;
        } else
          lines.push(body);
      } else {
        lines.push(stringify.stringify(doc.contents, ctx));
      }
      if (doc.directives?.docEnd) {
        if (doc.comment) {
          const cs = commentString(doc.comment);
          if (cs.includes("\n")) {
            lines.push("...");
            lines.push(stringifyComment.indentComment(cs, ""));
          } else {
            lines.push(`... ${cs}`);
          }
        } else {
          lines.push("...");
        }
      } else {
        let dc = doc.comment;
        if (dc && chompKeep)
          dc = dc.replace(/^\n+/, "");
        if (dc) {
          if ((!chompKeep || contentComment) && lines[lines.length - 1] !== "")
            lines.push("");
          lines.push(stringifyComment.indentComment(commentString(dc), ""));
        }
      }
      return lines.join("\n") + "\n";
    }
    exports.stringifyDocument = stringifyDocument;
  }
});

// node_modules/yaml/dist/doc/Document.js
var require_Document = __commonJS({
  "node_modules/yaml/dist/doc/Document.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Alias = require_Alias();
    var Collection = require_Collection();
    var identity = require_identity();
    var Pair = require_Pair();
    var toJS = require_toJS();
    var Schema = require_Schema();
    var stringifyDocument = require_stringifyDocument();
    var anchors = require_anchors();
    var applyReviver = require_applyReviver();
    var createNode = require_createNode();
    var directives = require_directives();
    var Document = class _Document {
      constructor(value2, replacer, options) {
        this.commentBefore = null;
        this.comment = null;
        this.errors = [];
        this.warnings = [];
        Object.defineProperty(this, identity.NODE_TYPE, { value: identity.DOC });
        let _replacer = null;
        if (typeof replacer === "function" || Array.isArray(replacer)) {
          _replacer = replacer;
        } else if (options === void 0 && replacer) {
          options = replacer;
          replacer = void 0;
        }
        const opt = Object.assign({
          intAsBigInt: false,
          keepSourceTokens: false,
          logLevel: "warn",
          prettyErrors: true,
          strict: true,
          stringKeys: false,
          uniqueKeys: true,
          version: "1.2"
        }, options);
        this.options = opt;
        let { version } = opt;
        if (options?._directives) {
          this.directives = options._directives.atDocument();
          if (this.directives.yaml.explicit)
            version = this.directives.yaml.version;
        } else
          this.directives = new directives.Directives({ version });
        this.setSchema(version, options);
        this.contents = value2 === void 0 ? null : this.createNode(value2, _replacer, options);
      }
      /**
       * Create a deep copy of this Document and its contents.
       *
       * Custom Node values that inherit from `Object` still refer to their original instances.
       */
      clone() {
        const copy = Object.create(_Document.prototype, {
          [identity.NODE_TYPE]: { value: identity.DOC }
        });
        copy.commentBefore = this.commentBefore;
        copy.comment = this.comment;
        copy.errors = this.errors.slice();
        copy.warnings = this.warnings.slice();
        copy.options = Object.assign({}, this.options);
        if (this.directives)
          copy.directives = this.directives.clone();
        copy.schema = this.schema.clone();
        copy.contents = identity.isNode(this.contents) ? this.contents.clone(copy.schema) : this.contents;
        if (this.range)
          copy.range = this.range.slice();
        return copy;
      }
      /** Adds a value to the document. */
      add(value2) {
        if (assertCollection(this.contents))
          this.contents.add(value2);
      }
      /** Adds a value to the document. */
      addIn(path2, value2) {
        if (assertCollection(this.contents))
          this.contents.addIn(path2, value2);
      }
      /**
       * Create a new `Alias` node, ensuring that the target `node` has the required anchor.
       *
       * If `node` already has an anchor, `name` is ignored.
       * Otherwise, the `node.anchor` value will be set to `name`,
       * or if an anchor with that name is already present in the document,
       * `name` will be used as a prefix for a new unique anchor.
       * If `name` is undefined, the generated anchor will use 'a' as a prefix.
       */
      createAlias(node, name) {
        if (!node.anchor) {
          const prev = anchors.anchorNames(this);
          node.anchor = // eslint-disable-next-line @typescript-eslint/prefer-nullish-coalescing
          !name || prev.has(name) ? anchors.findNewAnchor(name || "a", prev) : name;
        }
        return new Alias.Alias(node.anchor);
      }
      createNode(value2, replacer, options) {
        let _replacer = void 0;
        if (typeof replacer === "function") {
          value2 = replacer.call({ "": value2 }, "", value2);
          _replacer = replacer;
        } else if (Array.isArray(replacer)) {
          const keyToStr = (v) => typeof v === "number" || v instanceof String || v instanceof Number;
          const asStr = replacer.filter(keyToStr).map(String);
          if (asStr.length > 0)
            replacer = replacer.concat(asStr);
          _replacer = replacer;
        } else if (options === void 0 && replacer) {
          options = replacer;
          replacer = void 0;
        }
        const { aliasDuplicateObjects, anchorPrefix, flow, keepUndefined, onTagObj, tag } = options ?? {};
        const { onAnchor, setAnchors, sourceObjects } = anchors.createNodeAnchors(
          this,
          // eslint-disable-next-line @typescript-eslint/prefer-nullish-coalescing
          anchorPrefix || "a"
        );
        const ctx = {
          aliasDuplicateObjects: aliasDuplicateObjects ?? true,
          keepUndefined: keepUndefined ?? false,
          onAnchor,
          onTagObj,
          replacer: _replacer,
          schema: this.schema,
          sourceObjects
        };
        const node = createNode.createNode(value2, tag, ctx);
        if (flow && identity.isCollection(node))
          node.flow = true;
        setAnchors();
        return node;
      }
      /**
       * Convert a key and a value into a `Pair` using the current schema,
       * recursively wrapping all values as `Scalar` or `Collection` nodes.
       */
      createPair(key, value2, options = {}) {
        const k = this.createNode(key, null, options);
        const v = this.createNode(value2, null, options);
        return new Pair.Pair(k, v);
      }
      /**
       * Removes a value from the document.
       * @returns `true` if the item was found and removed.
       */
      delete(key) {
        return assertCollection(this.contents) ? this.contents.delete(key) : false;
      }
      /**
       * Removes a value from the document.
       * @returns `true` if the item was found and removed.
       */
      deleteIn(path2) {
        if (Collection.isEmptyPath(path2)) {
          if (this.contents == null)
            return false;
          this.contents = null;
          return true;
        }
        return assertCollection(this.contents) ? this.contents.deleteIn(path2) : false;
      }
      /**
       * Returns item at `key`, or `undefined` if not found. By default unwraps
       * scalar values from their surrounding node; to disable set `keepScalar` to
       * `true` (collections are always returned intact).
       */
      get(key, keepScalar) {
        return identity.isCollection(this.contents) ? this.contents.get(key, keepScalar) : void 0;
      }
      /**
       * Returns item at `path`, or `undefined` if not found. By default unwraps
       * scalar values from their surrounding node; to disable set `keepScalar` to
       * `true` (collections are always returned intact).
       */
      getIn(path2, keepScalar) {
        if (Collection.isEmptyPath(path2))
          return !keepScalar && identity.isScalar(this.contents) ? this.contents.value : this.contents;
        return identity.isCollection(this.contents) ? this.contents.getIn(path2, keepScalar) : void 0;
      }
      /**
       * Checks if the document includes a value with the key `key`.
       */
      has(key) {
        return identity.isCollection(this.contents) ? this.contents.has(key) : false;
      }
      /**
       * Checks if the document includes a value at `path`.
       */
      hasIn(path2) {
        if (Collection.isEmptyPath(path2))
          return this.contents !== void 0;
        return identity.isCollection(this.contents) ? this.contents.hasIn(path2) : false;
      }
      /**
       * Sets a value in this document. For `!!set`, `value` needs to be a
       * boolean to add/remove the item from the set.
       */
      set(key, value2) {
        if (this.contents == null) {
          this.contents = Collection.collectionFromPath(this.schema, [key], value2);
        } else if (assertCollection(this.contents)) {
          this.contents.set(key, value2);
        }
      }
      /**
       * Sets a value in this document. For `!!set`, `value` needs to be a
       * boolean to add/remove the item from the set.
       */
      setIn(path2, value2) {
        if (Collection.isEmptyPath(path2)) {
          this.contents = value2;
        } else if (this.contents == null) {
          this.contents = Collection.collectionFromPath(this.schema, Array.from(path2), value2);
        } else if (assertCollection(this.contents)) {
          this.contents.setIn(path2, value2);
        }
      }
      /**
       * Change the YAML version and schema used by the document.
       * A `null` version disables support for directives, explicit tags, anchors, and aliases.
       * It also requires the `schema` option to be given as a `Schema` instance value.
       *
       * Overrides all previously set schema options.
       */
      setSchema(version, options = {}) {
        if (typeof version === "number")
          version = String(version);
        let opt;
        switch (version) {
          case "1.1":
            if (this.directives)
              this.directives.yaml.version = "1.1";
            else
              this.directives = new directives.Directives({ version: "1.1" });
            opt = { resolveKnownTags: false, schema: "yaml-1.1" };
            break;
          case "1.2":
          case "next":
            if (this.directives)
              this.directives.yaml.version = version;
            else
              this.directives = new directives.Directives({ version });
            opt = { resolveKnownTags: true, schema: "core" };
            break;
          case null:
            if (this.directives)
              delete this.directives;
            opt = null;
            break;
          default: {
            const sv = JSON.stringify(version);
            throw new Error(`Expected '1.1', '1.2' or null as first argument, but found: ${sv}`);
          }
        }
        if (options.schema instanceof Object)
          this.schema = options.schema;
        else if (opt)
          this.schema = new Schema.Schema(Object.assign(opt, options));
        else
          throw new Error(`With a null YAML version, the { schema: Schema } option is required`);
      }
      // json & jsonArg are only used from toJSON()
      toJS({ json, jsonArg, mapAsMap, maxAliasCount, onAnchor, reviver } = {}) {
        const ctx = {
          anchors: /* @__PURE__ */ new Map(),
          doc: this,
          keep: !json,
          mapAsMap: mapAsMap === true,
          mapKeyWarned: false,
          maxAliasCount: typeof maxAliasCount === "number" ? maxAliasCount : 100
        };
        const res = toJS.toJS(this.contents, jsonArg ?? "", ctx);
        if (typeof onAnchor === "function")
          for (const { count, res: res2 } of ctx.anchors.values())
            onAnchor(res2, count);
        return typeof reviver === "function" ? applyReviver.applyReviver(reviver, { "": res }, "", res) : res;
      }
      /**
       * A JSON representation of the document `contents`.
       *
       * @param jsonArg Used by `JSON.stringify` to indicate the array index or
       *   property name.
       */
      toJSON(jsonArg, onAnchor) {
        return this.toJS({ json: true, jsonArg, mapAsMap: false, onAnchor });
      }
      /** A YAML representation of the document. */
      toString(options = {}) {
        if (this.errors.length > 0)
          throw new Error("Document with errors cannot be stringified");
        if ("indent" in options && (!Number.isInteger(options.indent) || Number(options.indent) <= 0)) {
          const s = JSON.stringify(options.indent);
          throw new Error(`"indent" option must be a positive integer, not ${s}`);
        }
        return stringifyDocument.stringifyDocument(this, options);
      }
    };
    function assertCollection(contents) {
      if (identity.isCollection(contents))
        return true;
      throw new Error("Expected a YAML collection as document contents");
    }
    exports.Document = Document;
  }
});

// node_modules/yaml/dist/errors.js
var require_errors2 = __commonJS({
  "node_modules/yaml/dist/errors.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var YAMLError = class extends Error {
      constructor(name, pos, code, message) {
        super();
        this.name = name;
        this.code = code;
        this.message = message;
        this.pos = pos;
      }
    };
    var YAMLParseError = class extends YAMLError {
      constructor(pos, code, message) {
        super("YAMLParseError", pos, code, message);
      }
    };
    var YAMLWarning = class extends YAMLError {
      constructor(pos, code, message) {
        super("YAMLWarning", pos, code, message);
      }
    };
    var prettifyError = (src, lc) => (error) => {
      if (error.pos[0] === -1)
        return;
      error.linePos = error.pos.map((pos) => lc.linePos(pos));
      const { line, col } = error.linePos[0];
      error.message += ` at line ${line}, column ${col}`;
      let ci = col - 1;
      let lineStr = src.substring(lc.lineStarts[line - 1], lc.lineStarts[line]).replace(/[\n\r]+$/, "");
      if (ci >= 60 && lineStr.length > 80) {
        const trimStart = Math.min(ci - 39, lineStr.length - 79);
        lineStr = "\u2026" + lineStr.substring(trimStart);
        ci -= trimStart - 1;
      }
      if (lineStr.length > 80)
        lineStr = lineStr.substring(0, 79) + "\u2026";
      if (line > 1 && /^ *$/.test(lineStr.substring(0, ci))) {
        let prev = src.substring(lc.lineStarts[line - 2], lc.lineStarts[line - 1]);
        if (prev.length > 80)
          prev = prev.substring(0, 79) + "\u2026\n";
        lineStr = prev + lineStr;
      }
      if (/[^ ]/.test(lineStr)) {
        let count = 1;
        const end = error.linePos[1];
        if (end?.line === line && end.col > col) {
          count = Math.max(1, Math.min(end.col - col, 80 - ci));
        }
        const pointer = " ".repeat(ci) + "^".repeat(count);
        error.message += `:

${lineStr}
${pointer}
`;
      }
    };
    exports.YAMLError = YAMLError;
    exports.YAMLParseError = YAMLParseError;
    exports.YAMLWarning = YAMLWarning;
    exports.prettifyError = prettifyError;
  }
});

// node_modules/yaml/dist/compose/resolve-props.js
var require_resolve_props = __commonJS({
  "node_modules/yaml/dist/compose/resolve-props.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    function resolveProps(tokens, { flow, indicator, next, offset, onError, parentIndent, startOnNewline }) {
      let spaceBefore = false;
      let atNewline = startOnNewline;
      let hasSpace = startOnNewline;
      let comment = "";
      let commentSep = "";
      let hasNewline = false;
      let reqSpace = false;
      let tab = null;
      let anchor = null;
      let tag = null;
      let newlineAfterProp = null;
      let comma = null;
      let found = null;
      let start = null;
      for (const token of tokens) {
        if (reqSpace) {
          if (token.type !== "space" && token.type !== "newline" && token.type !== "comma")
            onError(token.offset, "MISSING_CHAR", "Tags and anchors must be separated from the next token by white space");
          reqSpace = false;
        }
        if (tab) {
          if (atNewline && token.type !== "comment" && token.type !== "newline") {
            onError(tab, "TAB_AS_INDENT", "Tabs are not allowed as indentation");
          }
          tab = null;
        }
        switch (token.type) {
          case "space":
            if (!flow && (indicator !== "doc-start" || next?.type !== "flow-collection") && token.source.includes("	")) {
              tab = token;
            }
            hasSpace = true;
            break;
          case "comment": {
            if (!hasSpace)
              onError(token, "MISSING_CHAR", "Comments must be separated from other tokens by white space characters");
            const cb = token.source.substring(1) || " ";
            if (!comment)
              comment = cb;
            else
              comment += commentSep + cb;
            commentSep = "";
            atNewline = false;
            break;
          }
          case "newline":
            if (atNewline) {
              if (comment)
                comment += token.source;
              else if (!found || indicator !== "seq-item-ind")
                spaceBefore = true;
            } else
              commentSep += token.source;
            atNewline = true;
            hasNewline = true;
            if (anchor || tag)
              newlineAfterProp = token;
            hasSpace = true;
            break;
          case "anchor":
            if (anchor)
              onError(token, "MULTIPLE_ANCHORS", "A node can have at most one anchor");
            if (token.source.endsWith(":"))
              onError(token.offset + token.source.length - 1, "BAD_ALIAS", "Anchor ending in : is ambiguous", true);
            anchor = token;
            start ?? (start = token.offset);
            atNewline = false;
            hasSpace = false;
            reqSpace = true;
            break;
          case "tag": {
            if (tag)
              onError(token, "MULTIPLE_TAGS", "A node can have at most one tag");
            tag = token;
            start ?? (start = token.offset);
            atNewline = false;
            hasSpace = false;
            reqSpace = true;
            break;
          }
          case indicator:
            if (anchor || tag)
              onError(token, "BAD_PROP_ORDER", `Anchors and tags must be after the ${token.source} indicator`);
            if (found)
              onError(token, "UNEXPECTED_TOKEN", `Unexpected ${token.source} in ${flow ?? "collection"}`);
            found = token;
            atNewline = indicator === "seq-item-ind" || indicator === "explicit-key-ind";
            hasSpace = false;
            break;
          case "comma":
            if (flow) {
              if (comma)
                onError(token, "UNEXPECTED_TOKEN", `Unexpected , in ${flow}`);
              comma = token;
              atNewline = false;
              hasSpace = false;
              break;
            }
          // else fallthrough
          default:
            onError(token, "UNEXPECTED_TOKEN", `Unexpected ${token.type} token`);
            atNewline = false;
            hasSpace = false;
        }
      }
      const last = tokens[tokens.length - 1];
      const end = last ? last.offset + last.source.length : offset;
      if (reqSpace && next && next.type !== "space" && next.type !== "newline" && next.type !== "comma" && (next.type !== "scalar" || next.source !== "")) {
        onError(next.offset, "MISSING_CHAR", "Tags and anchors must be separated from the next token by white space");
      }
      if (tab && (atNewline && tab.indent <= parentIndent || next?.type === "block-map" || next?.type === "block-seq"))
        onError(tab, "TAB_AS_INDENT", "Tabs are not allowed as indentation");
      return {
        comma,
        found,
        spaceBefore,
        comment,
        hasNewline,
        anchor,
        tag,
        newlineAfterProp,
        end,
        start: start ?? end
      };
    }
    exports.resolveProps = resolveProps;
  }
});

// node_modules/yaml/dist/compose/util-contains-newline.js
var require_util_contains_newline = __commonJS({
  "node_modules/yaml/dist/compose/util-contains-newline.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    function containsNewline(key) {
      if (!key)
        return null;
      switch (key.type) {
        case "alias":
        case "scalar":
        case "double-quoted-scalar":
        case "single-quoted-scalar":
          if (key.source.includes("\n"))
            return true;
          if (key.end) {
            for (const st of key.end)
              if (st.type === "newline")
                return true;
          }
          return false;
        case "flow-collection":
          for (const it of key.items) {
            for (const st of it.start)
              if (st.type === "newline")
                return true;
            if (it.sep) {
              for (const st of it.sep)
                if (st.type === "newline")
                  return true;
            }
            if (containsNewline(it.key) || containsNewline(it.value))
              return true;
          }
          return false;
        default:
          return true;
      }
    }
    exports.containsNewline = containsNewline;
  }
});

// node_modules/yaml/dist/compose/util-flow-indent-check.js
var require_util_flow_indent_check = __commonJS({
  "node_modules/yaml/dist/compose/util-flow-indent-check.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var utilContainsNewline = require_util_contains_newline();
    function flowIndentCheck(indent, fc, onError) {
      if (fc?.type === "flow-collection") {
        const end = fc.end[0];
        if (end.indent === indent && (end.source === "]" || end.source === "}") && utilContainsNewline.containsNewline(fc)) {
          const msg = "Flow end indicator should be more indented than parent";
          onError(end, "BAD_INDENT", msg, true);
        }
      }
    }
    exports.flowIndentCheck = flowIndentCheck;
  }
});

// node_modules/yaml/dist/compose/util-map-includes.js
var require_util_map_includes = __commonJS({
  "node_modules/yaml/dist/compose/util-map-includes.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    function mapIncludes(ctx, items, search) {
      const { uniqueKeys } = ctx.options;
      if (uniqueKeys === false)
        return false;
      const isEqual = typeof uniqueKeys === "function" ? uniqueKeys : (a, b) => a === b || identity.isScalar(a) && identity.isScalar(b) && a.value === b.value;
      return items.some((pair) => isEqual(pair.key, search));
    }
    exports.mapIncludes = mapIncludes;
  }
});

// node_modules/yaml/dist/compose/resolve-block-map.js
var require_resolve_block_map = __commonJS({
  "node_modules/yaml/dist/compose/resolve-block-map.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Pair = require_Pair();
    var YAMLMap = require_YAMLMap();
    var resolveProps = require_resolve_props();
    var utilContainsNewline = require_util_contains_newline();
    var utilFlowIndentCheck = require_util_flow_indent_check();
    var utilMapIncludes = require_util_map_includes();
    var startColMsg = "All mapping items must start at the same column";
    function resolveBlockMap({ composeNode, composeEmptyNode }, ctx, bm, onError, tag) {
      const NodeClass = tag?.nodeClass ?? YAMLMap.YAMLMap;
      const map = new NodeClass(ctx.schema);
      if (ctx.atRoot)
        ctx.atRoot = false;
      let offset = bm.offset;
      let commentEnd = null;
      for (const collItem of bm.items) {
        const { start, key, sep: sep3, value: value2 } = collItem;
        const keyProps = resolveProps.resolveProps(start, {
          indicator: "explicit-key-ind",
          next: key ?? sep3?.[0],
          offset,
          onError,
          parentIndent: bm.indent,
          startOnNewline: true
        });
        const implicitKey = !keyProps.found;
        if (implicitKey) {
          if (key) {
            if (key.type === "block-seq")
              onError(offset, "BLOCK_AS_IMPLICIT_KEY", "A block sequence may not be used as an implicit map key");
            else if ("indent" in key && key.indent !== bm.indent)
              onError(offset, "BAD_INDENT", startColMsg);
          }
          if (!keyProps.anchor && !keyProps.tag && !sep3) {
            commentEnd = keyProps.end;
            if (keyProps.comment) {
              if (map.comment)
                map.comment += "\n" + keyProps.comment;
              else
                map.comment = keyProps.comment;
            }
            continue;
          }
          if (keyProps.newlineAfterProp || utilContainsNewline.containsNewline(key)) {
            onError(key ?? start[start.length - 1], "MULTILINE_IMPLICIT_KEY", "Implicit keys need to be on a single line");
          }
        } else if (keyProps.found?.indent !== bm.indent) {
          onError(offset, "BAD_INDENT", startColMsg);
        }
        ctx.atKey = true;
        const keyStart = keyProps.end;
        const keyNode = key ? composeNode(ctx, key, keyProps, onError) : composeEmptyNode(ctx, keyStart, start, null, keyProps, onError);
        if (ctx.schema.compat)
          utilFlowIndentCheck.flowIndentCheck(bm.indent, key, onError);
        ctx.atKey = false;
        if (utilMapIncludes.mapIncludes(ctx, map.items, keyNode))
          onError(keyStart, "DUPLICATE_KEY", "Map keys must be unique");
        const valueProps = resolveProps.resolveProps(sep3 ?? [], {
          indicator: "map-value-ind",
          next: value2,
          offset: keyNode.range[2],
          onError,
          parentIndent: bm.indent,
          startOnNewline: !key || key.type === "block-scalar"
        });
        offset = valueProps.end;
        if (valueProps.found) {
          if (implicitKey) {
            if (value2?.type === "block-map" && !valueProps.hasNewline)
              onError(offset, "BLOCK_AS_IMPLICIT_KEY", "Nested mappings are not allowed in compact mappings");
            if (ctx.options.strict && keyProps.start < valueProps.found.offset - 1024)
              onError(keyNode.range, "KEY_OVER_1024_CHARS", "The : indicator must be at most 1024 chars after the start of an implicit block mapping key");
          }
          const valueNode = value2 ? composeNode(ctx, value2, valueProps, onError) : composeEmptyNode(ctx, offset, sep3, null, valueProps, onError);
          if (ctx.schema.compat)
            utilFlowIndentCheck.flowIndentCheck(bm.indent, value2, onError);
          offset = valueNode.range[2];
          const pair = new Pair.Pair(keyNode, valueNode);
          if (ctx.options.keepSourceTokens)
            pair.srcToken = collItem;
          map.items.push(pair);
        } else {
          if (implicitKey)
            onError(keyNode.range, "MISSING_CHAR", "Implicit map keys need to be followed by map values");
          if (valueProps.comment) {
            if (keyNode.comment)
              keyNode.comment += "\n" + valueProps.comment;
            else
              keyNode.comment = valueProps.comment;
          }
          const pair = new Pair.Pair(keyNode);
          if (ctx.options.keepSourceTokens)
            pair.srcToken = collItem;
          map.items.push(pair);
        }
      }
      if (commentEnd && commentEnd < offset)
        onError(commentEnd, "IMPOSSIBLE", "Map comment with trailing content");
      map.range = [bm.offset, offset, commentEnd ?? offset];
      return map;
    }
    exports.resolveBlockMap = resolveBlockMap;
  }
});

// node_modules/yaml/dist/compose/resolve-block-seq.js
var require_resolve_block_seq = __commonJS({
  "node_modules/yaml/dist/compose/resolve-block-seq.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var YAMLSeq = require_YAMLSeq();
    var resolveProps = require_resolve_props();
    var utilFlowIndentCheck = require_util_flow_indent_check();
    function resolveBlockSeq({ composeNode, composeEmptyNode }, ctx, bs, onError, tag) {
      const NodeClass = tag?.nodeClass ?? YAMLSeq.YAMLSeq;
      const seq = new NodeClass(ctx.schema);
      if (ctx.atRoot)
        ctx.atRoot = false;
      if (ctx.atKey)
        ctx.atKey = false;
      let offset = bs.offset;
      let commentEnd = null;
      for (const { start, value: value2 } of bs.items) {
        const props = resolveProps.resolveProps(start, {
          indicator: "seq-item-ind",
          next: value2,
          offset,
          onError,
          parentIndent: bs.indent,
          startOnNewline: true
        });
        if (!props.found) {
          if (props.anchor || props.tag || value2) {
            if (value2?.type === "block-seq")
              onError(props.end, "BAD_INDENT", "All sequence items must start at the same column");
            else
              onError(offset, "MISSING_CHAR", "Sequence item without - indicator");
          } else {
            commentEnd = props.end;
            if (props.comment)
              seq.comment = props.comment;
            continue;
          }
        }
        const node = value2 ? composeNode(ctx, value2, props, onError) : composeEmptyNode(ctx, props.end, start, null, props, onError);
        if (ctx.schema.compat)
          utilFlowIndentCheck.flowIndentCheck(bs.indent, value2, onError);
        offset = node.range[2];
        seq.items.push(node);
      }
      seq.range = [bs.offset, offset, commentEnd ?? offset];
      return seq;
    }
    exports.resolveBlockSeq = resolveBlockSeq;
  }
});

// node_modules/yaml/dist/compose/resolve-end.js
var require_resolve_end = __commonJS({
  "node_modules/yaml/dist/compose/resolve-end.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    function resolveEnd(end, offset, reqSpace, onError) {
      let comment = "";
      if (end) {
        let hasSpace = false;
        let sep3 = "";
        for (const token of end) {
          const { source, type } = token;
          switch (type) {
            case "space":
              hasSpace = true;
              break;
            case "comment": {
              if (reqSpace && !hasSpace)
                onError(token, "MISSING_CHAR", "Comments must be separated from other tokens by white space characters");
              const cb = source.substring(1) || " ";
              if (!comment)
                comment = cb;
              else
                comment += sep3 + cb;
              sep3 = "";
              break;
            }
            case "newline":
              if (comment)
                sep3 += source;
              hasSpace = true;
              break;
            default:
              onError(token, "UNEXPECTED_TOKEN", `Unexpected ${type} at node end`);
          }
          offset += source.length;
        }
      }
      return { comment, offset };
    }
    exports.resolveEnd = resolveEnd;
  }
});

// node_modules/yaml/dist/compose/resolve-flow-collection.js
var require_resolve_flow_collection = __commonJS({
  "node_modules/yaml/dist/compose/resolve-flow-collection.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var Pair = require_Pair();
    var YAMLMap = require_YAMLMap();
    var YAMLSeq = require_YAMLSeq();
    var resolveEnd = require_resolve_end();
    var resolveProps = require_resolve_props();
    var utilContainsNewline = require_util_contains_newline();
    var utilMapIncludes = require_util_map_includes();
    var blockMsg = "Block collections are not allowed within flow collections";
    var isBlock = (token) => token && (token.type === "block-map" || token.type === "block-seq");
    function resolveFlowCollection({ composeNode, composeEmptyNode }, ctx, fc, onError, tag) {
      const isMap = fc.start.source === "{";
      const fcName = isMap ? "flow map" : "flow sequence";
      const NodeClass = tag?.nodeClass ?? (isMap ? YAMLMap.YAMLMap : YAMLSeq.YAMLSeq);
      const coll = new NodeClass(ctx.schema);
      coll.flow = true;
      const atRoot = ctx.atRoot;
      if (atRoot)
        ctx.atRoot = false;
      if (ctx.atKey)
        ctx.atKey = false;
      let offset = fc.offset + fc.start.source.length;
      for (let i = 0; i < fc.items.length; ++i) {
        const collItem = fc.items[i];
        const { start, key, sep: sep3, value: value2 } = collItem;
        const props = resolveProps.resolveProps(start, {
          flow: fcName,
          indicator: "explicit-key-ind",
          next: key ?? sep3?.[0],
          offset,
          onError,
          parentIndent: fc.indent,
          startOnNewline: false
        });
        if (!props.found) {
          if (!props.anchor && !props.tag && !sep3 && !value2) {
            if (i === 0 && props.comma)
              onError(props.comma, "UNEXPECTED_TOKEN", `Unexpected , in ${fcName}`);
            else if (i < fc.items.length - 1)
              onError(props.start, "UNEXPECTED_TOKEN", `Unexpected empty item in ${fcName}`);
            if (props.comment) {
              if (coll.comment)
                coll.comment += "\n" + props.comment;
              else
                coll.comment = props.comment;
            }
            offset = props.end;
            continue;
          }
          if (!isMap && ctx.options.strict && utilContainsNewline.containsNewline(key))
            onError(
              key,
              // checked by containsNewline()
              "MULTILINE_IMPLICIT_KEY",
              "Implicit keys of flow sequence pairs need to be on a single line"
            );
        }
        if (i === 0) {
          if (props.comma)
            onError(props.comma, "UNEXPECTED_TOKEN", `Unexpected , in ${fcName}`);
        } else {
          if (!props.comma)
            onError(props.start, "MISSING_CHAR", `Missing , between ${fcName} items`);
          if (props.comment) {
            let prevItemComment = "";
            loop: for (const st of start) {
              switch (st.type) {
                case "comma":
                case "space":
                  break;
                case "comment":
                  prevItemComment = st.source.substring(1);
                  break loop;
                default:
                  break loop;
              }
            }
            if (prevItemComment) {
              let prev = coll.items[coll.items.length - 1];
              if (identity.isPair(prev))
                prev = prev.value ?? prev.key;
              if (prev.comment)
                prev.comment += "\n" + prevItemComment;
              else
                prev.comment = prevItemComment;
              props.comment = props.comment.substring(prevItemComment.length + 1);
            }
          }
        }
        if (!isMap && !sep3 && !props.found) {
          const valueNode = value2 ? composeNode(ctx, value2, props, onError) : composeEmptyNode(ctx, props.end, sep3, null, props, onError);
          coll.items.push(valueNode);
          offset = valueNode.range[2];
          if (isBlock(value2))
            onError(valueNode.range, "BLOCK_IN_FLOW", blockMsg);
        } else {
          ctx.atKey = true;
          const keyStart = props.end;
          const keyNode = key ? composeNode(ctx, key, props, onError) : composeEmptyNode(ctx, keyStart, start, null, props, onError);
          if (isBlock(key))
            onError(keyNode.range, "BLOCK_IN_FLOW", blockMsg);
          ctx.atKey = false;
          const valueProps = resolveProps.resolveProps(sep3 ?? [], {
            flow: fcName,
            indicator: "map-value-ind",
            next: value2,
            offset: keyNode.range[2],
            onError,
            parentIndent: fc.indent,
            startOnNewline: false
          });
          if (valueProps.found) {
            if (!isMap && !props.found && ctx.options.strict) {
              if (sep3)
                for (const st of sep3) {
                  if (st === valueProps.found)
                    break;
                  if (st.type === "newline") {
                    onError(st, "MULTILINE_IMPLICIT_KEY", "Implicit keys of flow sequence pairs need to be on a single line");
                    break;
                  }
                }
              if (props.start < valueProps.found.offset - 1024)
                onError(valueProps.found, "KEY_OVER_1024_CHARS", "The : indicator must be at most 1024 chars after the start of an implicit flow sequence key");
            }
          } else if (value2) {
            if ("source" in value2 && value2.source?.[0] === ":")
              onError(value2, "MISSING_CHAR", `Missing space after : in ${fcName}`);
            else
              onError(valueProps.start, "MISSING_CHAR", `Missing , or : between ${fcName} items`);
          }
          const valueNode = value2 ? composeNode(ctx, value2, valueProps, onError) : valueProps.found ? composeEmptyNode(ctx, valueProps.end, sep3, null, valueProps, onError) : null;
          if (valueNode) {
            if (isBlock(value2))
              onError(valueNode.range, "BLOCK_IN_FLOW", blockMsg);
          } else if (valueProps.comment) {
            if (keyNode.comment)
              keyNode.comment += "\n" + valueProps.comment;
            else
              keyNode.comment = valueProps.comment;
          }
          const pair = new Pair.Pair(keyNode, valueNode);
          if (ctx.options.keepSourceTokens)
            pair.srcToken = collItem;
          if (isMap) {
            const map = coll;
            if (utilMapIncludes.mapIncludes(ctx, map.items, keyNode))
              onError(keyStart, "DUPLICATE_KEY", "Map keys must be unique");
            map.items.push(pair);
          } else {
            const map = new YAMLMap.YAMLMap(ctx.schema);
            map.flow = true;
            map.items.push(pair);
            const endRange = (valueNode ?? keyNode).range;
            map.range = [keyNode.range[0], endRange[1], endRange[2]];
            coll.items.push(map);
          }
          offset = valueNode ? valueNode.range[2] : valueProps.end;
        }
      }
      const expectedEnd = isMap ? "}" : "]";
      const [ce, ...ee] = fc.end;
      let cePos = offset;
      if (ce?.source === expectedEnd)
        cePos = ce.offset + ce.source.length;
      else {
        const name = fcName[0].toUpperCase() + fcName.substring(1);
        const msg = atRoot ? `${name} must end with a ${expectedEnd}` : `${name} in block collection must be sufficiently indented and end with a ${expectedEnd}`;
        onError(offset, atRoot ? "MISSING_CHAR" : "BAD_INDENT", msg);
        if (ce && ce.source.length !== 1)
          ee.unshift(ce);
      }
      if (ee.length > 0) {
        const end = resolveEnd.resolveEnd(ee, cePos, ctx.options.strict, onError);
        if (end.comment) {
          if (coll.comment)
            coll.comment += "\n" + end.comment;
          else
            coll.comment = end.comment;
        }
        coll.range = [fc.offset, cePos, end.offset];
      } else {
        coll.range = [fc.offset, cePos, cePos];
      }
      return coll;
    }
    exports.resolveFlowCollection = resolveFlowCollection;
  }
});

// node_modules/yaml/dist/compose/compose-collection.js
var require_compose_collection = __commonJS({
  "node_modules/yaml/dist/compose/compose-collection.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var Scalar = require_Scalar();
    var YAMLMap = require_YAMLMap();
    var YAMLSeq = require_YAMLSeq();
    var resolveBlockMap = require_resolve_block_map();
    var resolveBlockSeq = require_resolve_block_seq();
    var resolveFlowCollection = require_resolve_flow_collection();
    function resolveCollection(CN, ctx, token, onError, tagName, tag) {
      const coll = token.type === "block-map" ? resolveBlockMap.resolveBlockMap(CN, ctx, token, onError, tag) : token.type === "block-seq" ? resolveBlockSeq.resolveBlockSeq(CN, ctx, token, onError, tag) : resolveFlowCollection.resolveFlowCollection(CN, ctx, token, onError, tag);
      const Coll = coll.constructor;
      if (tagName === "!" || tagName === Coll.tagName) {
        coll.tag = Coll.tagName;
        return coll;
      }
      if (tagName)
        coll.tag = tagName;
      return coll;
    }
    function composeCollection(CN, ctx, token, props, onError) {
      const tagToken = props.tag;
      const tagName = !tagToken ? null : ctx.directives.tagName(tagToken.source, (msg) => onError(tagToken, "TAG_RESOLVE_FAILED", msg));
      if (token.type === "block-seq") {
        const { anchor, newlineAfterProp: nl } = props;
        const lastProp = anchor && tagToken ? anchor.offset > tagToken.offset ? anchor : tagToken : anchor ?? tagToken;
        if (lastProp && (!nl || nl.offset < lastProp.offset)) {
          const message = "Missing newline after block sequence props";
          onError(lastProp, "MISSING_CHAR", message);
        }
      }
      const expType = token.type === "block-map" ? "map" : token.type === "block-seq" ? "seq" : token.start.source === "{" ? "map" : "seq";
      if (!tagToken || !tagName || tagName === "!" || tagName === YAMLMap.YAMLMap.tagName && expType === "map" || tagName === YAMLSeq.YAMLSeq.tagName && expType === "seq") {
        return resolveCollection(CN, ctx, token, onError, tagName);
      }
      let tag = ctx.schema.tags.find((t) => t.tag === tagName && t.collection === expType);
      if (!tag) {
        const kt = ctx.schema.knownTags[tagName];
        if (kt?.collection === expType) {
          ctx.schema.tags.push(Object.assign({}, kt, { default: false }));
          tag = kt;
        } else {
          if (kt) {
            onError(tagToken, "BAD_COLLECTION_TYPE", `${kt.tag} used for ${expType} collection, but expects ${kt.collection ?? "scalar"}`, true);
          } else {
            onError(tagToken, "TAG_RESOLVE_FAILED", `Unresolved tag: ${tagName}`, true);
          }
          return resolveCollection(CN, ctx, token, onError, tagName);
        }
      }
      const coll = resolveCollection(CN, ctx, token, onError, tagName, tag);
      const res = tag.resolve?.(coll, (msg) => onError(tagToken, "TAG_RESOLVE_FAILED", msg), ctx.options) ?? coll;
      const node = identity.isNode(res) ? res : new Scalar.Scalar(res);
      node.range = coll.range;
      node.tag = tagName;
      if (tag?.format)
        node.format = tag.format;
      return node;
    }
    exports.composeCollection = composeCollection;
  }
});

// node_modules/yaml/dist/compose/resolve-block-scalar.js
var require_resolve_block_scalar = __commonJS({
  "node_modules/yaml/dist/compose/resolve-block-scalar.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Scalar = require_Scalar();
    function resolveBlockScalar(ctx, scalar, onError) {
      const start = scalar.offset;
      const header = parseBlockScalarHeader(scalar, ctx.options.strict, onError);
      if (!header)
        return { value: "", type: null, comment: "", range: [start, start, start] };
      const type = header.mode === ">" ? Scalar.Scalar.BLOCK_FOLDED : Scalar.Scalar.BLOCK_LITERAL;
      const lines = scalar.source ? splitLines(scalar.source) : [];
      let chompStart = lines.length;
      for (let i = lines.length - 1; i >= 0; --i) {
        const content = lines[i][1];
        if (content === "" || content === "\r")
          chompStart = i;
        else
          break;
      }
      if (chompStart === 0) {
        const value3 = header.chomp === "+" && lines.length > 0 ? "\n".repeat(Math.max(1, lines.length - 1)) : "";
        let end2 = start + header.length;
        if (scalar.source)
          end2 += scalar.source.length;
        return { value: value3, type, comment: header.comment, range: [start, end2, end2] };
      }
      let trimIndent = scalar.indent + header.indent;
      let offset = scalar.offset + header.length;
      let contentStart = 0;
      for (let i = 0; i < chompStart; ++i) {
        const [indent, content] = lines[i];
        if (content === "" || content === "\r") {
          if (header.indent === 0 && indent.length > trimIndent)
            trimIndent = indent.length;
        } else {
          if (indent.length < trimIndent) {
            const message = "Block scalars with more-indented leading empty lines must use an explicit indentation indicator";
            onError(offset + indent.length, "MISSING_CHAR", message);
          }
          if (header.indent === 0)
            trimIndent = indent.length;
          contentStart = i;
          if (trimIndent === 0 && !ctx.atRoot) {
            const message = "Block scalar values in collections must be indented";
            onError(offset, "BAD_INDENT", message);
          }
          break;
        }
        offset += indent.length + content.length + 1;
      }
      for (let i = lines.length - 1; i >= chompStart; --i) {
        if (lines[i][0].length > trimIndent)
          chompStart = i + 1;
      }
      let value2 = "";
      let sep3 = "";
      let prevMoreIndented = false;
      for (let i = 0; i < contentStart; ++i)
        value2 += lines[i][0].slice(trimIndent) + "\n";
      for (let i = contentStart; i < chompStart; ++i) {
        let [indent, content] = lines[i];
        offset += indent.length + content.length + 1;
        const crlf = content[content.length - 1] === "\r";
        if (crlf)
          content = content.slice(0, -1);
        if (content && indent.length < trimIndent) {
          const src = header.indent ? "explicit indentation indicator" : "first line";
          const message = `Block scalar lines must not be less indented than their ${src}`;
          onError(offset - content.length - (crlf ? 2 : 1), "BAD_INDENT", message);
          indent = "";
        }
        if (type === Scalar.Scalar.BLOCK_LITERAL) {
          value2 += sep3 + indent.slice(trimIndent) + content;
          sep3 = "\n";
        } else if (indent.length > trimIndent || content[0] === "	") {
          if (sep3 === " ")
            sep3 = "\n";
          else if (!prevMoreIndented && sep3 === "\n")
            sep3 = "\n\n";
          value2 += sep3 + indent.slice(trimIndent) + content;
          sep3 = "\n";
          prevMoreIndented = true;
        } else if (content === "") {
          if (sep3 === "\n")
            value2 += "\n";
          else
            sep3 = "\n";
        } else {
          value2 += sep3 + content;
          sep3 = " ";
          prevMoreIndented = false;
        }
      }
      switch (header.chomp) {
        case "-":
          break;
        case "+":
          for (let i = chompStart; i < lines.length; ++i)
            value2 += "\n" + lines[i][0].slice(trimIndent);
          if (value2[value2.length - 1] !== "\n")
            value2 += "\n";
          break;
        default:
          value2 += "\n";
      }
      const end = start + header.length + scalar.source.length;
      return { value: value2, type, comment: header.comment, range: [start, end, end] };
    }
    function parseBlockScalarHeader({ offset, props }, strict, onError) {
      if (props[0].type !== "block-scalar-header") {
        onError(props[0], "IMPOSSIBLE", "Block scalar header not found");
        return null;
      }
      const { source } = props[0];
      const mode = source[0];
      let indent = 0;
      let chomp = "";
      let error = -1;
      for (let i = 1; i < source.length; ++i) {
        const ch = source[i];
        if (!chomp && (ch === "-" || ch === "+"))
          chomp = ch;
        else {
          const n = Number(ch);
          if (!indent && n)
            indent = n;
          else if (error === -1)
            error = offset + i;
        }
      }
      if (error !== -1)
        onError(error, "UNEXPECTED_TOKEN", `Block scalar header includes extra characters: ${source}`);
      let hasSpace = false;
      let comment = "";
      let length = source.length;
      for (let i = 1; i < props.length; ++i) {
        const token = props[i];
        switch (token.type) {
          case "space":
            hasSpace = true;
          // fallthrough
          case "newline":
            length += token.source.length;
            break;
          case "comment":
            if (strict && !hasSpace) {
              const message = "Comments must be separated from other tokens by white space characters";
              onError(token, "MISSING_CHAR", message);
            }
            length += token.source.length;
            comment = token.source.substring(1);
            break;
          case "error":
            onError(token, "UNEXPECTED_TOKEN", token.message);
            length += token.source.length;
            break;
          /* istanbul ignore next should not happen */
          default: {
            const message = `Unexpected token in block scalar header: ${token.type}`;
            onError(token, "UNEXPECTED_TOKEN", message);
            const ts = token.source;
            if (ts && typeof ts === "string")
              length += ts.length;
          }
        }
      }
      return { mode, indent, chomp, comment, length };
    }
    function splitLines(source) {
      const split = source.split(/\n( *)/);
      const first = split[0];
      const m = first.match(/^( *)/);
      const line0 = m?.[1] ? [m[1], first.slice(m[1].length)] : ["", first];
      const lines = [line0];
      for (let i = 1; i < split.length; i += 2)
        lines.push([split[i], split[i + 1]]);
      return lines;
    }
    exports.resolveBlockScalar = resolveBlockScalar;
  }
});

// node_modules/yaml/dist/compose/resolve-flow-scalar.js
var require_resolve_flow_scalar = __commonJS({
  "node_modules/yaml/dist/compose/resolve-flow-scalar.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Scalar = require_Scalar();
    var resolveEnd = require_resolve_end();
    function resolveFlowScalar(scalar, strict, onError) {
      const { offset, type, source, end } = scalar;
      let _type;
      let value2;
      const _onError = (rel, code, msg) => onError(offset + rel, code, msg);
      switch (type) {
        case "scalar":
          _type = Scalar.Scalar.PLAIN;
          value2 = plainValue(source, _onError);
          break;
        case "single-quoted-scalar":
          _type = Scalar.Scalar.QUOTE_SINGLE;
          value2 = singleQuotedValue(source, _onError);
          break;
        case "double-quoted-scalar":
          _type = Scalar.Scalar.QUOTE_DOUBLE;
          value2 = doubleQuotedValue(source, _onError);
          break;
        /* istanbul ignore next should not happen */
        default:
          onError(scalar, "UNEXPECTED_TOKEN", `Expected a flow scalar value, but found: ${type}`);
          return {
            value: "",
            type: null,
            comment: "",
            range: [offset, offset + source.length, offset + source.length]
          };
      }
      const valueEnd = offset + source.length;
      const re = resolveEnd.resolveEnd(end, valueEnd, strict, onError);
      return {
        value: value2,
        type: _type,
        comment: re.comment,
        range: [offset, valueEnd, re.offset]
      };
    }
    function plainValue(source, onError) {
      let badChar = "";
      switch (source[0]) {
        /* istanbul ignore next should not happen */
        case "	":
          badChar = "a tab character";
          break;
        case ",":
          badChar = "flow indicator character ,";
          break;
        case "%":
          badChar = "directive indicator character %";
          break;
        case "|":
        case ">": {
          badChar = `block scalar indicator ${source[0]}`;
          break;
        }
        case "@":
        case "`": {
          badChar = `reserved character ${source[0]}`;
          break;
        }
      }
      if (badChar)
        onError(0, "BAD_SCALAR_START", `Plain value cannot start with ${badChar}`);
      return foldLines(source);
    }
    function singleQuotedValue(source, onError) {
      if (source[source.length - 1] !== "'" || source.length === 1)
        onError(source.length, "MISSING_CHAR", "Missing closing 'quote");
      return foldLines(source.slice(1, -1)).replace(/''/g, "'");
    }
    function foldLines(source) {
      let first, line;
      try {
        first = new RegExp("(.*?)(?<![ 	])[ 	]*\r?\n", "sy");
        line = new RegExp("[ 	]*(.*?)(?:(?<![ 	])[ 	]*)?\r?\n", "sy");
      } catch {
        first = /(.*?)[ \t]*\r?\n/sy;
        line = /[ \t]*(.*?)[ \t]*\r?\n/sy;
      }
      let match = first.exec(source);
      if (!match)
        return source;
      let res = match[1];
      let sep3 = " ";
      let pos = first.lastIndex;
      line.lastIndex = pos;
      while (match = line.exec(source)) {
        if (match[1] === "") {
          if (sep3 === "\n")
            res += sep3;
          else
            sep3 = "\n";
        } else {
          res += sep3 + match[1];
          sep3 = " ";
        }
        pos = line.lastIndex;
      }
      const last = /[ \t]*(.*)/sy;
      last.lastIndex = pos;
      match = last.exec(source);
      return res + sep3 + (match?.[1] ?? "");
    }
    function doubleQuotedValue(source, onError) {
      let res = "";
      for (let i = 1; i < source.length - 1; ++i) {
        const ch = source[i];
        if (ch === "\r" && source[i + 1] === "\n")
          continue;
        if (ch === "\n") {
          const { fold, offset } = foldNewline(source, i);
          res += fold;
          i = offset;
        } else if (ch === "\\") {
          let next = source[++i];
          const cc = escapeCodes[next];
          if (cc)
            res += cc;
          else if (next === "\n") {
            next = source[i + 1];
            while (next === " " || next === "	")
              next = source[++i + 1];
          } else if (next === "\r" && source[i + 1] === "\n") {
            next = source[++i + 1];
            while (next === " " || next === "	")
              next = source[++i + 1];
          } else if (next === "x" || next === "u" || next === "U") {
            const length = next === "x" ? 2 : next === "u" ? 4 : 8;
            res += parseCharCode(source, i + 1, length, onError);
            i += length;
          } else {
            const raw = source.substr(i - 1, 2);
            onError(i - 1, "BAD_DQ_ESCAPE", `Invalid escape sequence ${raw}`);
            res += raw;
          }
        } else if (ch === " " || ch === "	") {
          const wsStart = i;
          let next = source[i + 1];
          while (next === " " || next === "	")
            next = source[++i + 1];
          if (next !== "\n" && !(next === "\r" && source[i + 2] === "\n"))
            res += i > wsStart ? source.slice(wsStart, i + 1) : ch;
        } else {
          res += ch;
        }
      }
      if (source[source.length - 1] !== '"' || source.length === 1)
        onError(source.length, "MISSING_CHAR", 'Missing closing "quote');
      return res;
    }
    function foldNewline(source, offset) {
      let fold = "";
      let ch = source[offset + 1];
      while (ch === " " || ch === "	" || ch === "\n" || ch === "\r") {
        if (ch === "\r" && source[offset + 2] !== "\n")
          break;
        if (ch === "\n")
          fold += "\n";
        offset += 1;
        ch = source[offset + 1];
      }
      if (!fold)
        fold = " ";
      return { fold, offset };
    }
    var escapeCodes = {
      "0": "\0",
      // null character
      a: "\x07",
      // bell character
      b: "\b",
      // backspace
      e: "\x1B",
      // escape character
      f: "\f",
      // form feed
      n: "\n",
      // line feed
      r: "\r",
      // carriage return
      t: "	",
      // horizontal tab
      v: "\v",
      // vertical tab
      N: "\x85",
      // Unicode next line
      _: "\xA0",
      // Unicode non-breaking space
      L: "\u2028",
      // Unicode line separator
      P: "\u2029",
      // Unicode paragraph separator
      " ": " ",
      '"': '"',
      "/": "/",
      "\\": "\\",
      "	": "	"
    };
    function parseCharCode(source, offset, length, onError) {
      const cc = source.substr(offset, length);
      const ok = cc.length === length && /^[0-9a-fA-F]+$/.test(cc);
      const code = ok ? parseInt(cc, 16) : NaN;
      try {
        return String.fromCodePoint(code);
      } catch {
        const raw = source.substr(offset - 2, length + 2);
        onError(offset - 2, "BAD_DQ_ESCAPE", `Invalid escape sequence ${raw}`);
        return raw;
      }
    }
    exports.resolveFlowScalar = resolveFlowScalar;
  }
});

// node_modules/yaml/dist/compose/compose-scalar.js
var require_compose_scalar = __commonJS({
  "node_modules/yaml/dist/compose/compose-scalar.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var identity = require_identity();
    var Scalar = require_Scalar();
    var resolveBlockScalar = require_resolve_block_scalar();
    var resolveFlowScalar = require_resolve_flow_scalar();
    function composeScalar(ctx, token, tagToken, onError) {
      const { value: value2, type, comment, range } = token.type === "block-scalar" ? resolveBlockScalar.resolveBlockScalar(ctx, token, onError) : resolveFlowScalar.resolveFlowScalar(token, ctx.options.strict, onError);
      const tagName = tagToken ? ctx.directives.tagName(tagToken.source, (msg) => onError(tagToken, "TAG_RESOLVE_FAILED", msg)) : null;
      let tag;
      if (ctx.options.stringKeys && ctx.atKey) {
        tag = ctx.schema[identity.SCALAR];
      } else if (tagName)
        tag = findScalarTagByName(ctx.schema, value2, tagName, tagToken, onError);
      else if (token.type === "scalar")
        tag = findScalarTagByTest(ctx, value2, token, onError);
      else
        tag = ctx.schema[identity.SCALAR];
      let scalar;
      try {
        const res = tag.resolve(value2, (msg) => onError(tagToken ?? token, "TAG_RESOLVE_FAILED", msg), ctx.options);
        scalar = identity.isScalar(res) ? res : new Scalar.Scalar(res);
      } catch (error) {
        const msg = error instanceof Error ? error.message : String(error);
        onError(tagToken ?? token, "TAG_RESOLVE_FAILED", msg);
        scalar = new Scalar.Scalar(value2);
      }
      scalar.range = range;
      scalar.source = value2;
      if (type)
        scalar.type = type;
      if (tagName)
        scalar.tag = tagName;
      if (tag.format)
        scalar.format = tag.format;
      if (comment)
        scalar.comment = comment;
      return scalar;
    }
    function findScalarTagByName(schema2, value2, tagName, tagToken, onError) {
      if (tagName === "!")
        return schema2[identity.SCALAR];
      const matchWithTest = [];
      for (const tag of schema2.tags) {
        if (!tag.collection && tag.tag === tagName) {
          if (tag.default && tag.test)
            matchWithTest.push(tag);
          else
            return tag;
        }
      }
      for (const tag of matchWithTest)
        if (tag.test?.test(value2))
          return tag;
      const kt = schema2.knownTags[tagName];
      if (kt && !kt.collection) {
        schema2.tags.push(Object.assign({}, kt, { default: false, test: void 0 }));
        return kt;
      }
      onError(tagToken, "TAG_RESOLVE_FAILED", `Unresolved tag: ${tagName}`, tagName !== "tag:yaml.org,2002:str");
      return schema2[identity.SCALAR];
    }
    function findScalarTagByTest({ atKey, directives, schema: schema2 }, value2, token, onError) {
      const tag = schema2.tags.find((tag2) => (tag2.default === true || atKey && tag2.default === "key") && tag2.test?.test(value2)) || schema2[identity.SCALAR];
      if (schema2.compat) {
        const compat = schema2.compat.find((tag2) => tag2.default && tag2.test?.test(value2)) ?? schema2[identity.SCALAR];
        if (tag.tag !== compat.tag) {
          const ts = directives.tagString(tag.tag);
          const cs = directives.tagString(compat.tag);
          const msg = `Value may be parsed as either ${ts} or ${cs}`;
          onError(token, "TAG_RESOLVE_FAILED", msg, true);
        }
      }
      return tag;
    }
    exports.composeScalar = composeScalar;
  }
});

// node_modules/yaml/dist/compose/util-empty-scalar-position.js
var require_util_empty_scalar_position = __commonJS({
  "node_modules/yaml/dist/compose/util-empty-scalar-position.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    function emptyScalarPosition(offset, before, pos) {
      if (before) {
        pos ?? (pos = before.length);
        for (let i = pos - 1; i >= 0; --i) {
          let st = before[i];
          switch (st.type) {
            case "space":
            case "comment":
            case "newline":
              offset -= st.source.length;
              continue;
          }
          st = before[++i];
          while (st?.type === "space") {
            offset += st.source.length;
            st = before[++i];
          }
          break;
        }
      }
      return offset;
    }
    exports.emptyScalarPosition = emptyScalarPosition;
  }
});

// node_modules/yaml/dist/compose/compose-node.js
var require_compose_node = __commonJS({
  "node_modules/yaml/dist/compose/compose-node.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Alias = require_Alias();
    var identity = require_identity();
    var composeCollection = require_compose_collection();
    var composeScalar = require_compose_scalar();
    var resolveEnd = require_resolve_end();
    var utilEmptyScalarPosition = require_util_empty_scalar_position();
    var CN = { composeNode, composeEmptyNode };
    function composeNode(ctx, token, props, onError) {
      const atKey = ctx.atKey;
      const { spaceBefore, comment, anchor, tag } = props;
      let node;
      let isSrcToken = true;
      switch (token.type) {
        case "alias":
          node = composeAlias(ctx, token, onError);
          if (anchor || tag)
            onError(token, "ALIAS_PROPS", "An alias node must not specify any properties");
          break;
        case "scalar":
        case "single-quoted-scalar":
        case "double-quoted-scalar":
        case "block-scalar":
          node = composeScalar.composeScalar(ctx, token, tag, onError);
          if (anchor)
            node.anchor = anchor.source.substring(1);
          break;
        case "block-map":
        case "block-seq":
        case "flow-collection":
          try {
            node = composeCollection.composeCollection(CN, ctx, token, props, onError);
            if (anchor)
              node.anchor = anchor.source.substring(1);
          } catch (error) {
            const message = error instanceof Error ? error.message : String(error);
            onError(token, "RESOURCE_EXHAUSTION", message);
          }
          break;
        default: {
          const message = token.type === "error" ? token.message : `Unsupported token (type: ${token.type})`;
          onError(token, "UNEXPECTED_TOKEN", message);
          isSrcToken = false;
        }
      }
      node ?? (node = composeEmptyNode(ctx, token.offset, void 0, null, props, onError));
      if (anchor && node.anchor === "")
        onError(anchor, "BAD_ALIAS", "Anchor cannot be an empty string");
      if (atKey && ctx.options.stringKeys && (!identity.isScalar(node) || typeof node.value !== "string" || node.tag && node.tag !== "tag:yaml.org,2002:str")) {
        const msg = "With stringKeys, all keys must be strings";
        onError(tag ?? token, "NON_STRING_KEY", msg);
      }
      if (spaceBefore)
        node.spaceBefore = true;
      if (comment) {
        if (token.type === "scalar" && token.source === "")
          node.comment = comment;
        else
          node.commentBefore = comment;
      }
      if (ctx.options.keepSourceTokens && isSrcToken)
        node.srcToken = token;
      return node;
    }
    function composeEmptyNode(ctx, offset, before, pos, { spaceBefore, comment, anchor, tag, end }, onError) {
      const token = {
        type: "scalar",
        offset: utilEmptyScalarPosition.emptyScalarPosition(offset, before, pos),
        indent: -1,
        source: ""
      };
      const node = composeScalar.composeScalar(ctx, token, tag, onError);
      if (anchor) {
        node.anchor = anchor.source.substring(1);
        if (node.anchor === "")
          onError(anchor, "BAD_ALIAS", "Anchor cannot be an empty string");
      }
      if (spaceBefore)
        node.spaceBefore = true;
      if (comment) {
        node.comment = comment;
        node.range[2] = end;
      }
      return node;
    }
    function composeAlias({ options }, { offset, source, end }, onError) {
      const alias = new Alias.Alias(source.substring(1));
      if (alias.source === "")
        onError(offset, "BAD_ALIAS", "Alias cannot be an empty string");
      if (alias.source.endsWith(":"))
        onError(offset + source.length - 1, "BAD_ALIAS", "Alias ending in : is ambiguous", true);
      const valueEnd = offset + source.length;
      const re = resolveEnd.resolveEnd(end, valueEnd, options.strict, onError);
      alias.range = [offset, valueEnd, re.offset];
      if (re.comment)
        alias.comment = re.comment;
      return alias;
    }
    exports.composeEmptyNode = composeEmptyNode;
    exports.composeNode = composeNode;
  }
});

// node_modules/yaml/dist/compose/compose-doc.js
var require_compose_doc = __commonJS({
  "node_modules/yaml/dist/compose/compose-doc.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var Document = require_Document();
    var composeNode = require_compose_node();
    var resolveEnd = require_resolve_end();
    var resolveProps = require_resolve_props();
    function composeDoc(options, directives, { offset, start, value: value2, end }, onError) {
      const opts = Object.assign({ _directives: directives }, options);
      const doc = new Document.Document(void 0, opts);
      const ctx = {
        atKey: false,
        atRoot: true,
        directives: doc.directives,
        options: doc.options,
        schema: doc.schema
      };
      const props = resolveProps.resolveProps(start, {
        indicator: "doc-start",
        next: value2 ?? end?.[0],
        offset,
        onError,
        parentIndent: 0,
        startOnNewline: true
      });
      if (props.found) {
        doc.directives.docStart = true;
        if (value2 && (value2.type === "block-map" || value2.type === "block-seq") && !props.hasNewline)
          onError(props.end, "MISSING_CHAR", "Block collection cannot start on same line with directives-end marker");
      }
      doc.contents = value2 ? composeNode.composeNode(ctx, value2, props, onError) : composeNode.composeEmptyNode(ctx, props.end, start, null, props, onError);
      const contentEnd = doc.contents.range[2];
      const re = resolveEnd.resolveEnd(end, contentEnd, false, onError);
      if (re.comment)
        doc.comment = re.comment;
      doc.range = [offset, contentEnd, re.offset];
      return doc;
    }
    exports.composeDoc = composeDoc;
  }
});

// node_modules/yaml/dist/compose/composer.js
var require_composer = __commonJS({
  "node_modules/yaml/dist/compose/composer.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var node_process = __require("process");
    var directives = require_directives();
    var Document = require_Document();
    var errors2 = require_errors2();
    var identity = require_identity();
    var composeDoc = require_compose_doc();
    var resolveEnd = require_resolve_end();
    function getErrorPos(src) {
      if (typeof src === "number")
        return [src, src + 1];
      if (Array.isArray(src))
        return src.length === 2 ? src : [src[0], src[1]];
      const { offset, source } = src;
      return [offset, offset + (typeof source === "string" ? source.length : 1)];
    }
    function parsePrelude(prelude) {
      let comment = "";
      let atComment = false;
      let afterEmptyLine = false;
      for (let i = 0; i < prelude.length; ++i) {
        const source = prelude[i];
        switch (source[0]) {
          case "#":
            comment += (comment === "" ? "" : afterEmptyLine ? "\n\n" : "\n") + (source.substring(1) || " ");
            atComment = true;
            afterEmptyLine = false;
            break;
          case "%":
            if (prelude[i + 1]?.[0] !== "#")
              i += 1;
            atComment = false;
            break;
          default:
            if (!atComment)
              afterEmptyLine = true;
            atComment = false;
        }
      }
      return { comment, afterEmptyLine };
    }
    var Composer = class {
      constructor(options = {}) {
        this.doc = null;
        this.atDirectives = false;
        this.prelude = [];
        this.errors = [];
        this.warnings = [];
        this.onError = (source, code, message, warning) => {
          const pos = getErrorPos(source);
          if (warning)
            this.warnings.push(new errors2.YAMLWarning(pos, code, message));
          else
            this.errors.push(new errors2.YAMLParseError(pos, code, message));
        };
        this.directives = new directives.Directives({ version: options.version || "1.2" });
        this.options = options;
      }
      decorate(doc, afterDoc) {
        const { comment, afterEmptyLine } = parsePrelude(this.prelude);
        if (comment) {
          const dc = doc.contents;
          if (afterDoc) {
            doc.comment = doc.comment ? `${doc.comment}
${comment}` : comment;
          } else if (afterEmptyLine || doc.directives.docStart || !dc) {
            doc.commentBefore = comment;
          } else if (identity.isCollection(dc) && !dc.flow && dc.items.length > 0) {
            let it = dc.items[0];
            if (identity.isPair(it))
              it = it.key;
            const cb = it.commentBefore;
            it.commentBefore = cb ? `${comment}
${cb}` : comment;
          } else {
            const cb = dc.commentBefore;
            dc.commentBefore = cb ? `${comment}
${cb}` : comment;
          }
        }
        if (afterDoc) {
          for (let i = 0; i < this.errors.length; ++i)
            doc.errors.push(this.errors[i]);
          for (let i = 0; i < this.warnings.length; ++i)
            doc.warnings.push(this.warnings[i]);
        } else {
          doc.errors = this.errors;
          doc.warnings = this.warnings;
        }
        this.prelude = [];
        this.errors = [];
        this.warnings = [];
      }
      /**
       * Current stream status information.
       *
       * Mostly useful at the end of input for an empty stream.
       */
      streamInfo() {
        return {
          comment: parsePrelude(this.prelude).comment,
          directives: this.directives,
          errors: this.errors,
          warnings: this.warnings
        };
      }
      /**
       * Compose tokens into documents.
       *
       * @param forceDoc - If the stream contains no document, still emit a final document including any comments and directives that would be applied to a subsequent document.
       * @param endOffset - Should be set if `forceDoc` is also set, to set the document range end and to indicate errors correctly.
       */
      *compose(tokens, forceDoc = false, endOffset = -1) {
        for (const token of tokens)
          yield* this.next(token);
        yield* this.end(forceDoc, endOffset);
      }
      /** Advance the composer by one CST token. */
      *next(token) {
        if (node_process.env.LOG_STREAM)
          console.dir(token, { depth: null });
        switch (token.type) {
          case "directive":
            this.directives.add(token.source, (offset, message, warning) => {
              const pos = getErrorPos(token);
              pos[0] += offset;
              this.onError(pos, "BAD_DIRECTIVE", message, warning);
            });
            this.prelude.push(token.source);
            this.atDirectives = true;
            break;
          case "document": {
            const doc = composeDoc.composeDoc(this.options, this.directives, token, this.onError);
            if (this.atDirectives && !doc.directives.docStart)
              this.onError(token, "MISSING_CHAR", "Missing directives-end/doc-start indicator line");
            this.decorate(doc, false);
            if (this.doc)
              yield this.doc;
            this.doc = doc;
            this.atDirectives = false;
            break;
          }
          case "byte-order-mark":
          case "space":
            break;
          case "comment":
          case "newline":
            this.prelude.push(token.source);
            break;
          case "error": {
            const msg = token.source ? `${token.message}: ${JSON.stringify(token.source)}` : token.message;
            const error = new errors2.YAMLParseError(getErrorPos(token), "UNEXPECTED_TOKEN", msg);
            if (this.atDirectives || !this.doc)
              this.errors.push(error);
            else
              this.doc.errors.push(error);
            break;
          }
          case "doc-end": {
            if (!this.doc) {
              const msg = "Unexpected doc-end without preceding document";
              this.errors.push(new errors2.YAMLParseError(getErrorPos(token), "UNEXPECTED_TOKEN", msg));
              break;
            }
            this.doc.directives.docEnd = true;
            const end = resolveEnd.resolveEnd(token.end, token.offset + token.source.length, this.doc.options.strict, this.onError);
            this.decorate(this.doc, true);
            if (end.comment) {
              const dc = this.doc.comment;
              this.doc.comment = dc ? `${dc}
${end.comment}` : end.comment;
            }
            this.doc.range[2] = end.offset;
            break;
          }
          default:
            this.errors.push(new errors2.YAMLParseError(getErrorPos(token), "UNEXPECTED_TOKEN", `Unsupported token ${token.type}`));
        }
      }
      /**
       * Call at end of input to yield any remaining document.
       *
       * @param forceDoc - If the stream contains no document, still emit a final document including any comments and directives that would be applied to a subsequent document.
       * @param endOffset - Should be set if `forceDoc` is also set, to set the document range end and to indicate errors correctly.
       */
      *end(forceDoc = false, endOffset = -1) {
        if (this.doc) {
          this.decorate(this.doc, true);
          yield this.doc;
          this.doc = null;
        } else if (forceDoc) {
          const opts = Object.assign({ _directives: this.directives }, this.options);
          const doc = new Document.Document(void 0, opts);
          if (this.atDirectives)
            this.onError(endOffset, "MISSING_CHAR", "Missing directives-end indicator line");
          doc.range = [0, endOffset, endOffset];
          this.decorate(doc, false);
          yield doc;
        }
      }
    };
    exports.Composer = Composer;
  }
});

// node_modules/yaml/dist/parse/cst-scalar.js
var require_cst_scalar = __commonJS({
  "node_modules/yaml/dist/parse/cst-scalar.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var resolveBlockScalar = require_resolve_block_scalar();
    var resolveFlowScalar = require_resolve_flow_scalar();
    var errors2 = require_errors2();
    var stringifyString = require_stringifyString();
    function resolveAsScalar(token, strict = true, onError) {
      if (token) {
        const _onError = (pos, code, message) => {
          const offset = typeof pos === "number" ? pos : Array.isArray(pos) ? pos[0] : pos.offset;
          if (onError)
            onError(offset, code, message);
          else
            throw new errors2.YAMLParseError([offset, offset + 1], code, message);
        };
        switch (token.type) {
          case "scalar":
          case "single-quoted-scalar":
          case "double-quoted-scalar":
            return resolveFlowScalar.resolveFlowScalar(token, strict, _onError);
          case "block-scalar":
            return resolveBlockScalar.resolveBlockScalar({ options: { strict } }, token, _onError);
        }
      }
      return null;
    }
    function createScalarToken(value2, context) {
      const { implicitKey = false, indent, inFlow = false, offset = -1, type = "PLAIN" } = context;
      const source = stringifyString.stringifyString({ type, value: value2 }, {
        implicitKey,
        indent: indent > 0 ? " ".repeat(indent) : "",
        inFlow,
        options: { blockQuote: true, lineWidth: -1 }
      });
      const end = context.end ?? [
        { type: "newline", offset: -1, indent, source: "\n" }
      ];
      switch (source[0]) {
        case "|":
        case ">": {
          const he = source.indexOf("\n");
          const head = source.substring(0, he);
          const body = source.substring(he + 1) + "\n";
          const props = [
            { type: "block-scalar-header", offset, indent, source: head }
          ];
          if (!addEndtoBlockProps(props, end))
            props.push({ type: "newline", offset: -1, indent, source: "\n" });
          return { type: "block-scalar", offset, indent, props, source: body };
        }
        case '"':
          return { type: "double-quoted-scalar", offset, indent, source, end };
        case "'":
          return { type: "single-quoted-scalar", offset, indent, source, end };
        default:
          return { type: "scalar", offset, indent, source, end };
      }
    }
    function setScalarValue(token, value2, context = {}) {
      let { afterKey = false, implicitKey = false, inFlow = false, type } = context;
      let indent = "indent" in token ? token.indent : null;
      if (afterKey && typeof indent === "number")
        indent += 2;
      if (!type)
        switch (token.type) {
          case "single-quoted-scalar":
            type = "QUOTE_SINGLE";
            break;
          case "double-quoted-scalar":
            type = "QUOTE_DOUBLE";
            break;
          case "block-scalar": {
            const header = token.props[0];
            if (header.type !== "block-scalar-header")
              throw new Error("Invalid block scalar header");
            type = header.source[0] === ">" ? "BLOCK_FOLDED" : "BLOCK_LITERAL";
            break;
          }
          default:
            type = "PLAIN";
        }
      const source = stringifyString.stringifyString({ type, value: value2 }, {
        implicitKey: implicitKey || indent === null,
        indent: indent !== null && indent > 0 ? " ".repeat(indent) : "",
        inFlow,
        options: { blockQuote: true, lineWidth: -1 }
      });
      switch (source[0]) {
        case "|":
        case ">":
          setBlockScalarValue(token, source);
          break;
        case '"':
          setFlowScalarValue(token, source, "double-quoted-scalar");
          break;
        case "'":
          setFlowScalarValue(token, source, "single-quoted-scalar");
          break;
        default:
          setFlowScalarValue(token, source, "scalar");
      }
    }
    function setBlockScalarValue(token, source) {
      const he = source.indexOf("\n");
      const head = source.substring(0, he);
      const body = source.substring(he + 1) + "\n";
      if (token.type === "block-scalar") {
        const header = token.props[0];
        if (header.type !== "block-scalar-header")
          throw new Error("Invalid block scalar header");
        header.source = head;
        token.source = body;
      } else {
        const { offset } = token;
        const indent = "indent" in token ? token.indent : -1;
        const props = [
          { type: "block-scalar-header", offset, indent, source: head }
        ];
        if (!addEndtoBlockProps(props, "end" in token ? token.end : void 0))
          props.push({ type: "newline", offset: -1, indent, source: "\n" });
        for (const key of Object.keys(token))
          if (key !== "type" && key !== "offset")
            delete token[key];
        Object.assign(token, { type: "block-scalar", indent, props, source: body });
      }
    }
    function addEndtoBlockProps(props, end) {
      if (end)
        for (const st of end)
          switch (st.type) {
            case "space":
            case "comment":
              props.push(st);
              break;
            case "newline":
              props.push(st);
              return true;
          }
      return false;
    }
    function setFlowScalarValue(token, source, type) {
      switch (token.type) {
        case "scalar":
        case "double-quoted-scalar":
        case "single-quoted-scalar":
          token.type = type;
          token.source = source;
          break;
        case "block-scalar": {
          const end = token.props.slice(1);
          let oa = source.length;
          if (token.props[0].type === "block-scalar-header")
            oa -= token.props[0].source.length;
          for (const tok of end)
            tok.offset += oa;
          delete token.props;
          Object.assign(token, { type, source, end });
          break;
        }
        case "block-map":
        case "block-seq": {
          const offset = token.offset + source.length;
          const nl = { type: "newline", offset, indent: token.indent, source: "\n" };
          delete token.items;
          Object.assign(token, { type, source, end: [nl] });
          break;
        }
        default: {
          const indent = "indent" in token ? token.indent : -1;
          const end = "end" in token && Array.isArray(token.end) ? token.end.filter((st) => st.type === "space" || st.type === "comment" || st.type === "newline") : [];
          for (const key of Object.keys(token))
            if (key !== "type" && key !== "offset")
              delete token[key];
          Object.assign(token, { type, indent, source, end });
        }
      }
    }
    exports.createScalarToken = createScalarToken;
    exports.resolveAsScalar = resolveAsScalar;
    exports.setScalarValue = setScalarValue;
  }
});

// node_modules/yaml/dist/parse/cst-stringify.js
var require_cst_stringify = __commonJS({
  "node_modules/yaml/dist/parse/cst-stringify.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var stringify = (cst) => "type" in cst ? stringifyToken(cst) : stringifyItem(cst);
    function stringifyToken(token) {
      switch (token.type) {
        case "block-scalar": {
          let res = "";
          for (const tok of token.props)
            res += stringifyToken(tok);
          return res + token.source;
        }
        case "block-map":
        case "block-seq": {
          let res = "";
          for (const item of token.items)
            res += stringifyItem(item);
          return res;
        }
        case "flow-collection": {
          let res = token.start.source;
          for (const item of token.items)
            res += stringifyItem(item);
          for (const st of token.end)
            res += st.source;
          return res;
        }
        case "document": {
          let res = stringifyItem(token);
          if (token.end)
            for (const st of token.end)
              res += st.source;
          return res;
        }
        default: {
          let res = token.source;
          if ("end" in token && token.end)
            for (const st of token.end)
              res += st.source;
          return res;
        }
      }
    }
    function stringifyItem({ start, key, sep: sep3, value: value2 }) {
      let res = "";
      for (const st of start)
        res += st.source;
      if (key)
        res += stringifyToken(key);
      if (sep3)
        for (const st of sep3)
          res += st.source;
      if (value2)
        res += stringifyToken(value2);
      return res;
    }
    exports.stringify = stringify;
  }
});

// node_modules/yaml/dist/parse/cst-visit.js
var require_cst_visit = __commonJS({
  "node_modules/yaml/dist/parse/cst-visit.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var BREAK = /* @__PURE__ */ Symbol("break visit");
    var SKIP = /* @__PURE__ */ Symbol("skip children");
    var REMOVE = /* @__PURE__ */ Symbol("remove item");
    function visit(cst, visitor) {
      if ("type" in cst && cst.type === "document")
        cst = { start: cst.start, value: cst.value };
      _visit(Object.freeze([]), cst, visitor);
    }
    visit.BREAK = BREAK;
    visit.SKIP = SKIP;
    visit.REMOVE = REMOVE;
    visit.itemAtPath = (cst, path2) => {
      let item = cst;
      for (const [field, index] of path2) {
        const tok = item?.[field];
        if (tok && "items" in tok) {
          item = tok.items[index];
        } else
          return void 0;
      }
      return item;
    };
    visit.parentCollection = (cst, path2) => {
      const parent = visit.itemAtPath(cst, path2.slice(0, -1));
      const field = path2[path2.length - 1][0];
      const coll = parent?.[field];
      if (coll && "items" in coll)
        return coll;
      throw new Error("Parent collection not found");
    };
    function _visit(path2, item, visitor) {
      let ctrl = visitor(item, path2);
      if (typeof ctrl === "symbol")
        return ctrl;
      for (const field of ["key", "value"]) {
        const token = item[field];
        if (token && "items" in token) {
          for (let i = 0; i < token.items.length; ++i) {
            const ci = _visit(Object.freeze(path2.concat([[field, i]])), token.items[i], visitor);
            if (typeof ci === "number")
              i = ci - 1;
            else if (ci === BREAK)
              return BREAK;
            else if (ci === REMOVE) {
              token.items.splice(i, 1);
              i -= 1;
            }
          }
          if (typeof ctrl === "function" && field === "key")
            ctrl = ctrl(item, path2);
        }
      }
      return typeof ctrl === "function" ? ctrl(item, path2) : ctrl;
    }
    exports.visit = visit;
  }
});

// node_modules/yaml/dist/parse/cst.js
var require_cst = __commonJS({
  "node_modules/yaml/dist/parse/cst.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var cstScalar = require_cst_scalar();
    var cstStringify = require_cst_stringify();
    var cstVisit = require_cst_visit();
    var BOM = "\uFEFF";
    var DOCUMENT = "";
    var FLOW_END = "";
    var SCALAR = "";
    var isCollection = (token) => !!token && "items" in token;
    var isScalar = (token) => !!token && (token.type === "scalar" || token.type === "single-quoted-scalar" || token.type === "double-quoted-scalar" || token.type === "block-scalar");
    function prettyToken(token) {
      switch (token) {
        case BOM:
          return "<BOM>";
        case DOCUMENT:
          return "<DOC>";
        case FLOW_END:
          return "<FLOW_END>";
        case SCALAR:
          return "<SCALAR>";
        default:
          return JSON.stringify(token);
      }
    }
    function tokenType(source) {
      switch (source) {
        case BOM:
          return "byte-order-mark";
        case DOCUMENT:
          return "doc-mode";
        case FLOW_END:
          return "flow-error-end";
        case SCALAR:
          return "scalar";
        case "---":
          return "doc-start";
        case "...":
          return "doc-end";
        case "":
        case "\n":
        case "\r\n":
          return "newline";
        case "-":
          return "seq-item-ind";
        case "?":
          return "explicit-key-ind";
        case ":":
          return "map-value-ind";
        case "{":
          return "flow-map-start";
        case "}":
          return "flow-map-end";
        case "[":
          return "flow-seq-start";
        case "]":
          return "flow-seq-end";
        case ",":
          return "comma";
      }
      switch (source[0]) {
        case " ":
        case "	":
          return "space";
        case "#":
          return "comment";
        case "%":
          return "directive-line";
        case "*":
          return "alias";
        case "&":
          return "anchor";
        case "!":
          return "tag";
        case "'":
          return "single-quoted-scalar";
        case '"':
          return "double-quoted-scalar";
        case "|":
        case ">":
          return "block-scalar-header";
      }
      return null;
    }
    exports.createScalarToken = cstScalar.createScalarToken;
    exports.resolveAsScalar = cstScalar.resolveAsScalar;
    exports.setScalarValue = cstScalar.setScalarValue;
    exports.stringify = cstStringify.stringify;
    exports.visit = cstVisit.visit;
    exports.BOM = BOM;
    exports.DOCUMENT = DOCUMENT;
    exports.FLOW_END = FLOW_END;
    exports.SCALAR = SCALAR;
    exports.isCollection = isCollection;
    exports.isScalar = isScalar;
    exports.prettyToken = prettyToken;
    exports.tokenType = tokenType;
  }
});

// node_modules/yaml/dist/parse/lexer.js
var require_lexer = __commonJS({
  "node_modules/yaml/dist/parse/lexer.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var cst = require_cst();
    function isEmpty(ch) {
      switch (ch) {
        case void 0:
        case " ":
        case "\n":
        case "\r":
        case "	":
          return true;
        default:
          return false;
      }
    }
    var hexDigits = new Set("0123456789ABCDEFabcdef");
    var tagChars = new Set("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-#;/?:@&=+$_.!~*'()");
    var flowIndicatorChars = new Set(",[]{}");
    var invalidAnchorChars = new Set(" ,[]{}\n\r	");
    var isNotAnchorChar = (ch) => !ch || invalidAnchorChars.has(ch);
    var Lexer = class {
      constructor() {
        this.atEnd = false;
        this.blockScalarIndent = -1;
        this.blockScalarKeep = false;
        this.buffer = "";
        this.flowKey = false;
        this.flowLevel = 0;
        this.indentNext = 0;
        this.indentValue = 0;
        this.lineEndPos = null;
        this.next = null;
        this.pos = 0;
      }
      /**
       * Generate YAML tokens from the `source` string. If `incomplete`,
       * a part of the last line may be left as a buffer for the next call.
       *
       * @returns A generator of lexical tokens
       */
      *lex(source, incomplete = false) {
        if (source) {
          if (typeof source !== "string")
            throw TypeError("source is not a string");
          this.buffer = this.buffer ? this.buffer + source : source;
          this.lineEndPos = null;
        }
        this.atEnd = !incomplete;
        let next = this.next ?? "stream";
        while (next && (incomplete || this.hasChars(1)))
          next = yield* this.parseNext(next);
      }
      atLineEnd() {
        let i = this.pos;
        let ch = this.buffer[i];
        while (ch === " " || ch === "	")
          ch = this.buffer[++i];
        if (!ch || ch === "#" || ch === "\n")
          return true;
        if (ch === "\r")
          return this.buffer[i + 1] === "\n";
        return false;
      }
      charAt(n) {
        return this.buffer[this.pos + n];
      }
      continueScalar(offset) {
        let ch = this.buffer[offset];
        if (this.indentNext > 0) {
          let indent = 0;
          while (ch === " ")
            ch = this.buffer[++indent + offset];
          if (ch === "\r") {
            const next = this.buffer[indent + offset + 1];
            if (next === "\n" || !next && !this.atEnd)
              return offset + indent + 1;
          }
          return ch === "\n" || indent >= this.indentNext || !ch && !this.atEnd ? offset + indent : -1;
        }
        if (ch === "-" || ch === ".") {
          const dt = this.buffer.substr(offset, 3);
          if ((dt === "---" || dt === "...") && isEmpty(this.buffer[offset + 3]))
            return -1;
        }
        return offset;
      }
      getLine() {
        let end = this.lineEndPos;
        if (typeof end !== "number" || end !== -1 && end < this.pos) {
          end = this.buffer.indexOf("\n", this.pos);
          this.lineEndPos = end;
        }
        if (end === -1)
          return this.atEnd ? this.buffer.substring(this.pos) : null;
        if (this.buffer[end - 1] === "\r")
          end -= 1;
        return this.buffer.substring(this.pos, end);
      }
      hasChars(n) {
        return this.pos + n <= this.buffer.length;
      }
      setNext(state) {
        this.buffer = this.buffer.substring(this.pos);
        this.pos = 0;
        this.lineEndPos = null;
        this.next = state;
        return null;
      }
      peek(n) {
        return this.buffer.substr(this.pos, n);
      }
      *parseNext(next) {
        switch (next) {
          case "stream":
            return yield* this.parseStream();
          case "line-start":
            return yield* this.parseLineStart();
          case "block-start":
            return yield* this.parseBlockStart();
          case "doc":
            return yield* this.parseDocument();
          case "flow":
            return yield* this.parseFlowCollection();
          case "quoted-scalar":
            return yield* this.parseQuotedScalar();
          case "block-scalar":
            return yield* this.parseBlockScalar();
          case "plain-scalar":
            return yield* this.parsePlainScalar();
        }
      }
      *parseStream() {
        let line = this.getLine();
        if (line === null)
          return this.setNext("stream");
        if (line[0] === cst.BOM) {
          yield* this.pushCount(1);
          line = line.substring(1);
        }
        if (line[0] === "%") {
          let dirEnd = line.length;
          let cs = line.indexOf("#");
          while (cs !== -1) {
            const ch = line[cs - 1];
            if (ch === " " || ch === "	") {
              dirEnd = cs - 1;
              break;
            } else {
              cs = line.indexOf("#", cs + 1);
            }
          }
          while (true) {
            const ch = line[dirEnd - 1];
            if (ch === " " || ch === "	")
              dirEnd -= 1;
            else
              break;
          }
          const n = (yield* this.pushCount(dirEnd)) + (yield* this.pushSpaces(true));
          yield* this.pushCount(line.length - n);
          this.pushNewline();
          return "stream";
        }
        if (this.atLineEnd()) {
          const sp = yield* this.pushSpaces(true);
          yield* this.pushCount(line.length - sp);
          yield* this.pushNewline();
          return "stream";
        }
        yield cst.DOCUMENT;
        return yield* this.parseLineStart();
      }
      *parseLineStart() {
        const ch = this.charAt(0);
        if (!ch && !this.atEnd)
          return this.setNext("line-start");
        if (ch === "-" || ch === ".") {
          if (!this.atEnd && !this.hasChars(4))
            return this.setNext("line-start");
          const s = this.peek(3);
          if ((s === "---" || s === "...") && isEmpty(this.charAt(3))) {
            yield* this.pushCount(3);
            this.indentValue = 0;
            this.indentNext = 0;
            return s === "---" ? "doc" : "stream";
          }
        }
        this.indentValue = yield* this.pushSpaces(false);
        if (this.indentNext > this.indentValue && !isEmpty(this.charAt(1)))
          this.indentNext = this.indentValue;
        return yield* this.parseBlockStart();
      }
      *parseBlockStart() {
        const [ch0, ch1] = this.peek(2);
        if (!ch1 && !this.atEnd)
          return this.setNext("block-start");
        if ((ch0 === "-" || ch0 === "?" || ch0 === ":") && isEmpty(ch1)) {
          const n = (yield* this.pushCount(1)) + (yield* this.pushSpaces(true));
          this.indentNext = this.indentValue + 1;
          this.indentValue += n;
          return "block-start";
        }
        return "doc";
      }
      *parseDocument() {
        yield* this.pushSpaces(true);
        const line = this.getLine();
        if (line === null)
          return this.setNext("doc");
        let n = yield* this.pushIndicators();
        switch (line[n]) {
          case "#":
            yield* this.pushCount(line.length - n);
          // fallthrough
          case void 0:
            yield* this.pushNewline();
            return yield* this.parseLineStart();
          case "{":
          case "[":
            yield* this.pushCount(1);
            this.flowKey = false;
            this.flowLevel = 1;
            return "flow";
          case "}":
          case "]":
            yield* this.pushCount(1);
            return "doc";
          case "*":
            yield* this.pushUntil(isNotAnchorChar);
            return "doc";
          case '"':
          case "'":
            return yield* this.parseQuotedScalar();
          case "|":
          case ">":
            n += yield* this.parseBlockScalarHeader();
            n += yield* this.pushSpaces(true);
            yield* this.pushCount(line.length - n);
            yield* this.pushNewline();
            return yield* this.parseBlockScalar();
          default:
            return yield* this.parsePlainScalar();
        }
      }
      *parseFlowCollection() {
        let nl, sp;
        let indent = -1;
        do {
          nl = yield* this.pushNewline();
          if (nl > 0) {
            sp = yield* this.pushSpaces(false);
            this.indentValue = indent = sp;
          } else {
            sp = 0;
          }
          sp += yield* this.pushSpaces(true);
        } while (nl + sp > 0);
        const line = this.getLine();
        if (line === null)
          return this.setNext("flow");
        if (indent !== -1 && indent < this.indentNext && line[0] !== "#" || indent === 0 && (line.startsWith("---") || line.startsWith("...")) && isEmpty(line[3])) {
          const atFlowEndMarker = indent === this.indentNext - 1 && this.flowLevel === 1 && (line[0] === "]" || line[0] === "}");
          if (!atFlowEndMarker) {
            this.flowLevel = 0;
            yield cst.FLOW_END;
            return yield* this.parseLineStart();
          }
        }
        let n = 0;
        while (line[n] === ",") {
          n += yield* this.pushCount(1);
          n += yield* this.pushSpaces(true);
          this.flowKey = false;
        }
        n += yield* this.pushIndicators();
        switch (line[n]) {
          case void 0:
            return "flow";
          case "#":
            yield* this.pushCount(line.length - n);
            return "flow";
          case "{":
          case "[":
            yield* this.pushCount(1);
            this.flowKey = false;
            this.flowLevel += 1;
            return "flow";
          case "}":
          case "]":
            yield* this.pushCount(1);
            this.flowKey = true;
            this.flowLevel -= 1;
            return this.flowLevel ? "flow" : "doc";
          case "*":
            yield* this.pushUntil(isNotAnchorChar);
            return "flow";
          case '"':
          case "'":
            this.flowKey = true;
            return yield* this.parseQuotedScalar();
          case ":": {
            const next = this.charAt(1);
            if (this.flowKey || isEmpty(next) || next === ",") {
              this.flowKey = false;
              yield* this.pushCount(1);
              yield* this.pushSpaces(true);
              return "flow";
            }
          }
          // fallthrough
          default:
            this.flowKey = false;
            return yield* this.parsePlainScalar();
        }
      }
      *parseQuotedScalar() {
        const quote = this.charAt(0);
        let end = this.buffer.indexOf(quote, this.pos + 1);
        if (quote === "'") {
          while (end !== -1 && this.buffer[end + 1] === "'")
            end = this.buffer.indexOf("'", end + 2);
        } else {
          while (end !== -1) {
            let n = 0;
            while (this.buffer[end - 1 - n] === "\\")
              n += 1;
            if (n % 2 === 0)
              break;
            end = this.buffer.indexOf('"', end + 1);
          }
        }
        const qb = this.buffer.substring(0, end);
        let nl = qb.indexOf("\n", this.pos);
        if (nl !== -1) {
          while (nl !== -1) {
            const cs = this.continueScalar(nl + 1);
            if (cs === -1)
              break;
            nl = qb.indexOf("\n", cs);
          }
          if (nl !== -1) {
            end = nl - (qb[nl - 1] === "\r" ? 2 : 1);
          }
        }
        if (end === -1) {
          if (!this.atEnd)
            return this.setNext("quoted-scalar");
          end = this.buffer.length;
        }
        yield* this.pushToIndex(end + 1, false);
        return this.flowLevel ? "flow" : "doc";
      }
      *parseBlockScalarHeader() {
        this.blockScalarIndent = -1;
        this.blockScalarKeep = false;
        let i = this.pos;
        while (true) {
          const ch = this.buffer[++i];
          if (ch === "+")
            this.blockScalarKeep = true;
          else if (ch > "0" && ch <= "9")
            this.blockScalarIndent = Number(ch) - 1;
          else if (ch !== "-")
            break;
        }
        return yield* this.pushUntil((ch) => isEmpty(ch) || ch === "#");
      }
      *parseBlockScalar() {
        let nl = this.pos - 1;
        let indent = 0;
        let ch;
        loop: for (let i2 = this.pos; ch = this.buffer[i2]; ++i2) {
          switch (ch) {
            case " ":
              indent += 1;
              break;
            case "\n":
              nl = i2;
              indent = 0;
              break;
            case "\r": {
              const next = this.buffer[i2 + 1];
              if (!next && !this.atEnd)
                return this.setNext("block-scalar");
              if (next === "\n")
                break;
            }
            // fallthrough
            default:
              break loop;
          }
        }
        if (!ch && !this.atEnd)
          return this.setNext("block-scalar");
        if (indent >= this.indentNext) {
          if (this.blockScalarIndent === -1)
            this.indentNext = indent;
          else {
            this.indentNext = this.blockScalarIndent + (this.indentNext === 0 ? 1 : this.indentNext);
          }
          do {
            const cs = this.continueScalar(nl + 1);
            if (cs === -1)
              break;
            nl = this.buffer.indexOf("\n", cs);
          } while (nl !== -1);
          if (nl === -1) {
            if (!this.atEnd)
              return this.setNext("block-scalar");
            nl = this.buffer.length;
          }
        }
        let i = nl + 1;
        ch = this.buffer[i];
        while (ch === " ")
          ch = this.buffer[++i];
        if (ch === "	") {
          while (ch === "	" || ch === " " || ch === "\r" || ch === "\n")
            ch = this.buffer[++i];
          nl = i - 1;
        } else if (!this.blockScalarKeep) {
          do {
            let i2 = nl - 1;
            let ch2 = this.buffer[i2];
            if (ch2 === "\r")
              ch2 = this.buffer[--i2];
            const lastChar = i2;
            while (ch2 === " ")
              ch2 = this.buffer[--i2];
            if (ch2 === "\n" && i2 >= this.pos && i2 + 1 + indent > lastChar)
              nl = i2;
            else
              break;
          } while (true);
        }
        yield cst.SCALAR;
        yield* this.pushToIndex(nl + 1, true);
        return yield* this.parseLineStart();
      }
      *parsePlainScalar() {
        const inFlow = this.flowLevel > 0;
        let end = this.pos - 1;
        let i = this.pos - 1;
        let ch;
        while (ch = this.buffer[++i]) {
          if (ch === ":") {
            const next = this.buffer[i + 1];
            if (isEmpty(next) || inFlow && flowIndicatorChars.has(next))
              break;
            end = i;
          } else if (isEmpty(ch)) {
            let next = this.buffer[i + 1];
            if (ch === "\r") {
              if (next === "\n") {
                i += 1;
                ch = "\n";
                next = this.buffer[i + 1];
              } else
                end = i;
            }
            if (next === "#" || inFlow && flowIndicatorChars.has(next))
              break;
            if (ch === "\n") {
              const cs = this.continueScalar(i + 1);
              if (cs === -1)
                break;
              i = Math.max(i, cs - 2);
            }
          } else {
            if (inFlow && flowIndicatorChars.has(ch))
              break;
            end = i;
          }
        }
        if (!ch && !this.atEnd)
          return this.setNext("plain-scalar");
        yield cst.SCALAR;
        yield* this.pushToIndex(end + 1, true);
        return inFlow ? "flow" : "doc";
      }
      *pushCount(n) {
        if (n > 0) {
          yield this.buffer.substr(this.pos, n);
          this.pos += n;
          return n;
        }
        return 0;
      }
      *pushToIndex(i, allowEmpty) {
        const s = this.buffer.slice(this.pos, i);
        if (s) {
          yield s;
          this.pos += s.length;
          return s.length;
        } else if (allowEmpty)
          yield "";
        return 0;
      }
      *pushIndicators() {
        let n = 0;
        loop: while (true) {
          switch (this.charAt(0)) {
            case "!":
              n += yield* this.pushTag();
              n += yield* this.pushSpaces(true);
              continue loop;
            case "&":
              n += yield* this.pushUntil(isNotAnchorChar);
              n += yield* this.pushSpaces(true);
              continue loop;
            case "-":
            // this is an error
            case "?":
            // this is an error outside flow collections
            case ":": {
              const inFlow = this.flowLevel > 0;
              const ch1 = this.charAt(1);
              if (isEmpty(ch1) || inFlow && flowIndicatorChars.has(ch1)) {
                if (!inFlow)
                  this.indentNext = this.indentValue + 1;
                else if (this.flowKey)
                  this.flowKey = false;
                n += yield* this.pushCount(1);
                n += yield* this.pushSpaces(true);
                continue loop;
              }
            }
          }
          break loop;
        }
        return n;
      }
      *pushTag() {
        if (this.charAt(1) === "<") {
          let i = this.pos + 2;
          let ch = this.buffer[i];
          while (!isEmpty(ch) && ch !== ">")
            ch = this.buffer[++i];
          return yield* this.pushToIndex(ch === ">" ? i + 1 : i, false);
        } else {
          let i = this.pos + 1;
          let ch = this.buffer[i];
          while (ch) {
            if (tagChars.has(ch))
              ch = this.buffer[++i];
            else if (ch === "%" && hexDigits.has(this.buffer[i + 1]) && hexDigits.has(this.buffer[i + 2])) {
              ch = this.buffer[i += 3];
            } else
              break;
          }
          return yield* this.pushToIndex(i, false);
        }
      }
      *pushNewline() {
        const ch = this.buffer[this.pos];
        if (ch === "\n")
          return yield* this.pushCount(1);
        else if (ch === "\r" && this.charAt(1) === "\n")
          return yield* this.pushCount(2);
        else
          return 0;
      }
      *pushSpaces(allowTabs) {
        let i = this.pos - 1;
        let ch;
        do {
          ch = this.buffer[++i];
        } while (ch === " " || allowTabs && ch === "	");
        const n = i - this.pos;
        if (n > 0) {
          yield this.buffer.substr(this.pos, n);
          this.pos = i;
        }
        return n;
      }
      *pushUntil(test) {
        let i = this.pos;
        let ch = this.buffer[i];
        while (!test(ch))
          ch = this.buffer[++i];
        return yield* this.pushToIndex(i, false);
      }
    };
    exports.Lexer = Lexer;
  }
});

// node_modules/yaml/dist/parse/line-counter.js
var require_line_counter = __commonJS({
  "node_modules/yaml/dist/parse/line-counter.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var LineCounter = class {
      constructor() {
        this.lineStarts = [];
        this.addNewLine = (offset) => this.lineStarts.push(offset);
        this.linePos = (offset) => {
          let low = 0;
          let high = this.lineStarts.length;
          while (low < high) {
            const mid = low + high >> 1;
            if (this.lineStarts[mid] < offset)
              low = mid + 1;
            else
              high = mid;
          }
          if (this.lineStarts[low] === offset)
            return { line: low + 1, col: 1 };
          if (low === 0)
            return { line: 0, col: offset };
          const start = this.lineStarts[low - 1];
          return { line: low, col: offset - start + 1 };
        };
      }
    };
    exports.LineCounter = LineCounter;
  }
});

// node_modules/yaml/dist/parse/parser.js
var require_parser = __commonJS({
  "node_modules/yaml/dist/parse/parser.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var node_process = __require("process");
    var cst = require_cst();
    var lexer = require_lexer();
    function includesToken(list2, type) {
      for (let i = 0; i < list2.length; ++i)
        if (list2[i].type === type)
          return true;
      return false;
    }
    function findNonEmptyIndex(list2) {
      for (let i = 0; i < list2.length; ++i) {
        switch (list2[i].type) {
          case "space":
          case "comment":
          case "newline":
            break;
          default:
            return i;
        }
      }
      return -1;
    }
    function isFlowToken(token) {
      switch (token?.type) {
        case "alias":
        case "scalar":
        case "single-quoted-scalar":
        case "double-quoted-scalar":
        case "flow-collection":
          return true;
        default:
          return false;
      }
    }
    function getPrevProps(parent) {
      switch (parent.type) {
        case "document":
          return parent.start;
        case "block-map": {
          const it = parent.items[parent.items.length - 1];
          return it.sep ?? it.start;
        }
        case "block-seq":
          return parent.items[parent.items.length - 1].start;
        /* istanbul ignore next should not happen */
        default:
          return [];
      }
    }
    function getFirstKeyStartProps(prev) {
      if (prev.length === 0)
        return [];
      let i = prev.length;
      loop: while (--i >= 0) {
        switch (prev[i].type) {
          case "doc-start":
          case "explicit-key-ind":
          case "map-value-ind":
          case "seq-item-ind":
          case "newline":
            break loop;
        }
      }
      while (prev[++i]?.type === "space") {
      }
      return prev.splice(i, prev.length);
    }
    function arrayPushArray(target, source) {
      if (source.length < 1e5)
        Array.prototype.push.apply(target, source);
      else
        for (let i = 0; i < source.length; ++i)
          target.push(source[i]);
    }
    function fixFlowSeqItems(fc) {
      if (fc.start.type === "flow-seq-start") {
        for (const it of fc.items) {
          if (it.sep && !it.value && !includesToken(it.start, "explicit-key-ind") && !includesToken(it.sep, "map-value-ind")) {
            if (it.key)
              it.value = it.key;
            delete it.key;
            if (isFlowToken(it.value)) {
              if (it.value.end)
                arrayPushArray(it.value.end, it.sep);
              else
                it.value.end = it.sep;
            } else
              arrayPushArray(it.start, it.sep);
            delete it.sep;
          }
        }
      }
    }
    var Parser = class {
      /**
       * @param onNewLine - If defined, called separately with the start position of
       *   each new line (in `parse()`, including the start of input).
       */
      constructor(onNewLine) {
        this.atNewLine = true;
        this.atScalar = false;
        this.indent = 0;
        this.offset = 0;
        this.onKeyLine = false;
        this.stack = [];
        this.source = "";
        this.type = "";
        this.lexer = new lexer.Lexer();
        this.onNewLine = onNewLine;
      }
      /**
       * Parse `source` as a YAML stream.
       * If `incomplete`, a part of the last line may be left as a buffer for the next call.
       *
       * Errors are not thrown, but yielded as `{ type: 'error', message }` tokens.
       *
       * @returns A generator of tokens representing each directive, document, and other structure.
       */
      *parse(source, incomplete = false) {
        if (this.onNewLine && this.offset === 0)
          this.onNewLine(0);
        for (const lexeme of this.lexer.lex(source, incomplete))
          yield* this.next(lexeme);
        if (!incomplete)
          yield* this.end();
      }
      /**
       * Advance the parser by the `source` of one lexical token.
       */
      *next(source) {
        this.source = source;
        if (node_process.env.LOG_TOKENS)
          console.log("|", cst.prettyToken(source));
        if (this.atScalar) {
          this.atScalar = false;
          yield* this.step();
          this.offset += source.length;
          return;
        }
        const type = cst.tokenType(source);
        if (!type) {
          const message = `Not a YAML token: ${source}`;
          yield* this.pop({ type: "error", offset: this.offset, message, source });
          this.offset += source.length;
        } else if (type === "scalar") {
          this.atNewLine = false;
          this.atScalar = true;
          this.type = "scalar";
        } else {
          this.type = type;
          yield* this.step();
          switch (type) {
            case "newline":
              this.atNewLine = true;
              this.indent = 0;
              if (this.onNewLine)
                this.onNewLine(this.offset + source.length);
              break;
            case "space":
              if (this.atNewLine && source[0] === " ")
                this.indent += source.length;
              break;
            case "explicit-key-ind":
            case "map-value-ind":
            case "seq-item-ind":
              if (this.atNewLine)
                this.indent += source.length;
              break;
            case "doc-mode":
            case "flow-error-end":
              return;
            default:
              this.atNewLine = false;
          }
          this.offset += source.length;
        }
      }
      /** Call at end of input to push out any remaining constructions */
      *end() {
        while (this.stack.length > 0)
          yield* this.pop();
      }
      get sourceToken() {
        const st = {
          type: this.type,
          offset: this.offset,
          indent: this.indent,
          source: this.source
        };
        return st;
      }
      *step() {
        const top = this.peek(1);
        if (this.type === "doc-end" && top?.type !== "doc-end") {
          while (this.stack.length > 0)
            yield* this.pop();
          this.stack.push({
            type: "doc-end",
            offset: this.offset,
            source: this.source
          });
          return;
        }
        if (!top)
          return yield* this.stream();
        switch (top.type) {
          case "document":
            return yield* this.document(top);
          case "alias":
          case "scalar":
          case "single-quoted-scalar":
          case "double-quoted-scalar":
            return yield* this.scalar(top);
          case "block-scalar":
            return yield* this.blockScalar(top);
          case "block-map":
            return yield* this.blockMap(top);
          case "block-seq":
            return yield* this.blockSequence(top);
          case "flow-collection":
            return yield* this.flowCollection(top);
          case "doc-end":
            return yield* this.documentEnd(top);
        }
        yield* this.pop();
      }
      peek(n) {
        return this.stack[this.stack.length - n];
      }
      *pop(error) {
        const token = error ?? this.stack.pop();
        if (!token) {
          const message = "Tried to pop an empty stack";
          yield { type: "error", offset: this.offset, source: "", message };
        } else if (this.stack.length === 0) {
          yield token;
        } else {
          const top = this.peek(1);
          if (token.type === "block-scalar") {
            token.indent = "indent" in top ? top.indent : 0;
          } else if (token.type === "flow-collection" && top.type === "document") {
            token.indent = 0;
          }
          if (token.type === "flow-collection")
            fixFlowSeqItems(token);
          switch (top.type) {
            case "document":
              top.value = token;
              break;
            case "block-scalar":
              top.props.push(token);
              break;
            case "block-map": {
              const it = top.items[top.items.length - 1];
              if (it.value) {
                top.items.push({ start: [], key: token, sep: [] });
                this.onKeyLine = true;
                return;
              } else if (it.sep) {
                it.value = token;
              } else {
                Object.assign(it, { key: token, sep: [] });
                this.onKeyLine = !it.explicitKey;
                return;
              }
              break;
            }
            case "block-seq": {
              const it = top.items[top.items.length - 1];
              if (it.value)
                top.items.push({ start: [], value: token });
              else
                it.value = token;
              break;
            }
            case "flow-collection": {
              const it = top.items[top.items.length - 1];
              if (!it || it.value)
                top.items.push({ start: [], key: token, sep: [] });
              else if (it.sep)
                it.value = token;
              else
                Object.assign(it, { key: token, sep: [] });
              return;
            }
            /* istanbul ignore next should not happen */
            default:
              yield* this.pop();
              yield* this.pop(token);
          }
          if ((top.type === "document" || top.type === "block-map" || top.type === "block-seq") && (token.type === "block-map" || token.type === "block-seq")) {
            const last = token.items[token.items.length - 1];
            if (last && !last.sep && !last.value && last.start.length > 0 && findNonEmptyIndex(last.start) === -1 && (token.indent === 0 || last.start.every((st) => st.type !== "comment" || st.indent < token.indent))) {
              if (top.type === "document")
                top.end = last.start;
              else
                top.items.push({ start: last.start });
              token.items.splice(-1, 1);
            }
          }
        }
      }
      *stream() {
        switch (this.type) {
          case "directive-line":
            yield { type: "directive", offset: this.offset, source: this.source };
            return;
          case "byte-order-mark":
          case "space":
          case "comment":
          case "newline":
            yield this.sourceToken;
            return;
          case "doc-mode":
          case "doc-start": {
            const doc = {
              type: "document",
              offset: this.offset,
              start: []
            };
            if (this.type === "doc-start")
              doc.start.push(this.sourceToken);
            this.stack.push(doc);
            return;
          }
        }
        yield {
          type: "error",
          offset: this.offset,
          message: `Unexpected ${this.type} token in YAML stream`,
          source: this.source
        };
      }
      *document(doc) {
        if (doc.value)
          return yield* this.lineEnd(doc);
        switch (this.type) {
          case "doc-start": {
            if (findNonEmptyIndex(doc.start) !== -1) {
              yield* this.pop();
              yield* this.step();
            } else
              doc.start.push(this.sourceToken);
            return;
          }
          case "anchor":
          case "tag":
          case "space":
          case "comment":
          case "newline":
            doc.start.push(this.sourceToken);
            return;
        }
        const bv = this.startBlockValue(doc);
        if (bv)
          this.stack.push(bv);
        else {
          yield {
            type: "error",
            offset: this.offset,
            message: `Unexpected ${this.type} token in YAML document`,
            source: this.source
          };
        }
      }
      *scalar(scalar) {
        if (this.type === "map-value-ind") {
          const prev = getPrevProps(this.peek(2));
          const start = getFirstKeyStartProps(prev);
          let sep3;
          if (scalar.end) {
            sep3 = scalar.end;
            sep3.push(this.sourceToken);
            delete scalar.end;
          } else
            sep3 = [this.sourceToken];
          const map = {
            type: "block-map",
            offset: scalar.offset,
            indent: scalar.indent,
            items: [{ start, key: scalar, sep: sep3 }]
          };
          this.onKeyLine = true;
          this.stack[this.stack.length - 1] = map;
        } else
          yield* this.lineEnd(scalar);
      }
      *blockScalar(scalar) {
        switch (this.type) {
          case "space":
          case "comment":
          case "newline":
            scalar.props.push(this.sourceToken);
            return;
          case "scalar":
            scalar.source = this.source;
            this.atNewLine = true;
            this.indent = 0;
            if (this.onNewLine) {
              let nl = this.source.indexOf("\n") + 1;
              while (nl !== 0) {
                this.onNewLine(this.offset + nl);
                nl = this.source.indexOf("\n", nl) + 1;
              }
            }
            yield* this.pop();
            break;
          /* istanbul ignore next should not happen */
          default:
            yield* this.pop();
            yield* this.step();
        }
      }
      *blockMap(map) {
        const it = map.items[map.items.length - 1];
        switch (this.type) {
          case "newline":
            this.onKeyLine = false;
            if (it.value) {
              const end = "end" in it.value ? it.value.end : void 0;
              const last = Array.isArray(end) ? end[end.length - 1] : void 0;
              if (last?.type === "comment")
                end?.push(this.sourceToken);
              else
                map.items.push({ start: [this.sourceToken] });
            } else if (it.sep) {
              it.sep.push(this.sourceToken);
            } else {
              it.start.push(this.sourceToken);
            }
            return;
          case "space":
          case "comment":
            if (it.value) {
              map.items.push({ start: [this.sourceToken] });
            } else if (it.sep) {
              it.sep.push(this.sourceToken);
            } else {
              if (this.atIndentedComment(it.start, map.indent)) {
                const prev = map.items[map.items.length - 2];
                const end = prev?.value?.end;
                if (Array.isArray(end)) {
                  arrayPushArray(end, it.start);
                  end.push(this.sourceToken);
                  map.items.pop();
                  return;
                }
              }
              it.start.push(this.sourceToken);
            }
            return;
        }
        if (this.indent >= map.indent) {
          const atMapIndent = !this.onKeyLine && this.indent === map.indent;
          const atNextItem = atMapIndent && (it.sep || it.explicitKey) && this.type !== "seq-item-ind";
          let start = [];
          if (atNextItem && it.sep && !it.value) {
            const nl = [];
            for (let i = 0; i < it.sep.length; ++i) {
              const st = it.sep[i];
              switch (st.type) {
                case "newline":
                  nl.push(i);
                  break;
                case "space":
                  break;
                case "comment":
                  if (st.indent > map.indent)
                    nl.length = 0;
                  break;
                default:
                  nl.length = 0;
              }
            }
            if (nl.length >= 2)
              start = it.sep.splice(nl[1]);
          }
          switch (this.type) {
            case "anchor":
            case "tag":
              if (atNextItem || it.value) {
                start.push(this.sourceToken);
                map.items.push({ start });
                this.onKeyLine = true;
              } else if (it.sep) {
                it.sep.push(this.sourceToken);
              } else {
                it.start.push(this.sourceToken);
              }
              return;
            case "explicit-key-ind":
              if (!it.sep && !it.explicitKey) {
                it.start.push(this.sourceToken);
                it.explicitKey = true;
              } else if (atNextItem || it.value) {
                start.push(this.sourceToken);
                map.items.push({ start, explicitKey: true });
              } else {
                this.stack.push({
                  type: "block-map",
                  offset: this.offset,
                  indent: this.indent,
                  items: [{ start: [this.sourceToken], explicitKey: true }]
                });
              }
              this.onKeyLine = true;
              return;
            case "map-value-ind":
              if (it.explicitKey) {
                if (!it.sep) {
                  if (includesToken(it.start, "newline")) {
                    Object.assign(it, { key: null, sep: [this.sourceToken] });
                  } else {
                    const start2 = getFirstKeyStartProps(it.start);
                    this.stack.push({
                      type: "block-map",
                      offset: this.offset,
                      indent: this.indent,
                      items: [{ start: start2, key: null, sep: [this.sourceToken] }]
                    });
                  }
                } else if (it.value) {
                  map.items.push({ start: [], key: null, sep: [this.sourceToken] });
                } else if (includesToken(it.sep, "map-value-ind")) {
                  this.stack.push({
                    type: "block-map",
                    offset: this.offset,
                    indent: this.indent,
                    items: [{ start, key: null, sep: [this.sourceToken] }]
                  });
                } else if (isFlowToken(it.key) && !includesToken(it.sep, "newline")) {
                  const start2 = getFirstKeyStartProps(it.start);
                  const key = it.key;
                  const sep3 = it.sep;
                  sep3.push(this.sourceToken);
                  delete it.key;
                  delete it.sep;
                  this.stack.push({
                    type: "block-map",
                    offset: this.offset,
                    indent: this.indent,
                    items: [{ start: start2, key, sep: sep3 }]
                  });
                } else if (start.length > 0) {
                  it.sep = it.sep.concat(start, this.sourceToken);
                } else {
                  it.sep.push(this.sourceToken);
                }
              } else {
                if (!it.sep) {
                  Object.assign(it, { key: null, sep: [this.sourceToken] });
                } else if (it.value || atNextItem) {
                  map.items.push({ start, key: null, sep: [this.sourceToken] });
                } else if (includesToken(it.sep, "map-value-ind")) {
                  this.stack.push({
                    type: "block-map",
                    offset: this.offset,
                    indent: this.indent,
                    items: [{ start: [], key: null, sep: [this.sourceToken] }]
                  });
                } else {
                  it.sep.push(this.sourceToken);
                }
              }
              this.onKeyLine = true;
              return;
            case "alias":
            case "scalar":
            case "single-quoted-scalar":
            case "double-quoted-scalar": {
              const fs = this.flowScalar(this.type);
              if (atNextItem || it.value) {
                map.items.push({ start, key: fs, sep: [] });
                this.onKeyLine = true;
              } else if (it.sep) {
                this.stack.push(fs);
              } else {
                Object.assign(it, { key: fs, sep: [] });
                this.onKeyLine = true;
              }
              return;
            }
            default: {
              const bv = this.startBlockValue(map);
              if (bv) {
                if (bv.type === "block-seq") {
                  if (!it.explicitKey && it.sep && !includesToken(it.sep, "newline")) {
                    yield* this.pop({
                      type: "error",
                      offset: this.offset,
                      message: "Unexpected block-seq-ind on same line with key",
                      source: this.source
                    });
                    return;
                  }
                } else if (atMapIndent) {
                  map.items.push({ start });
                }
                this.stack.push(bv);
                return;
              }
            }
          }
        }
        yield* this.pop();
        yield* this.step();
      }
      *blockSequence(seq) {
        const it = seq.items[seq.items.length - 1];
        switch (this.type) {
          case "newline":
            if (it.value) {
              const end = "end" in it.value ? it.value.end : void 0;
              const last = Array.isArray(end) ? end[end.length - 1] : void 0;
              if (last?.type === "comment")
                end?.push(this.sourceToken);
              else
                seq.items.push({ start: [this.sourceToken] });
            } else
              it.start.push(this.sourceToken);
            return;
          case "space":
          case "comment":
            if (it.value)
              seq.items.push({ start: [this.sourceToken] });
            else {
              if (this.atIndentedComment(it.start, seq.indent)) {
                const prev = seq.items[seq.items.length - 2];
                const end = prev?.value?.end;
                if (Array.isArray(end)) {
                  arrayPushArray(end, it.start);
                  end.push(this.sourceToken);
                  seq.items.pop();
                  return;
                }
              }
              it.start.push(this.sourceToken);
            }
            return;
          case "anchor":
          case "tag":
            if (it.value || this.indent <= seq.indent)
              break;
            it.start.push(this.sourceToken);
            return;
          case "seq-item-ind":
            if (this.indent !== seq.indent)
              break;
            if (it.value || includesToken(it.start, "seq-item-ind"))
              seq.items.push({ start: [this.sourceToken] });
            else
              it.start.push(this.sourceToken);
            return;
        }
        if (this.indent > seq.indent) {
          const bv = this.startBlockValue(seq);
          if (bv) {
            this.stack.push(bv);
            return;
          }
        }
        yield* this.pop();
        yield* this.step();
      }
      *flowCollection(fc) {
        const it = fc.items[fc.items.length - 1];
        if (this.type === "flow-error-end") {
          let top;
          do {
            yield* this.pop();
            top = this.peek(1);
          } while (top?.type === "flow-collection");
        } else if (fc.end.length === 0) {
          switch (this.type) {
            case "comma":
            case "explicit-key-ind":
              if (!it || it.sep)
                fc.items.push({ start: [this.sourceToken] });
              else
                it.start.push(this.sourceToken);
              return;
            case "map-value-ind":
              if (!it || it.value)
                fc.items.push({ start: [], key: null, sep: [this.sourceToken] });
              else if (it.sep)
                it.sep.push(this.sourceToken);
              else
                Object.assign(it, { key: null, sep: [this.sourceToken] });
              return;
            case "space":
            case "comment":
            case "newline":
            case "anchor":
            case "tag":
              if (!it || it.value)
                fc.items.push({ start: [this.sourceToken] });
              else if (it.sep)
                it.sep.push(this.sourceToken);
              else
                it.start.push(this.sourceToken);
              return;
            case "alias":
            case "scalar":
            case "single-quoted-scalar":
            case "double-quoted-scalar": {
              const fs = this.flowScalar(this.type);
              if (!it || it.value)
                fc.items.push({ start: [], key: fs, sep: [] });
              else if (it.sep)
                this.stack.push(fs);
              else
                Object.assign(it, { key: fs, sep: [] });
              return;
            }
            case "flow-map-end":
            case "flow-seq-end":
              fc.end.push(this.sourceToken);
              return;
          }
          const bv = this.startBlockValue(fc);
          if (bv)
            this.stack.push(bv);
          else {
            yield* this.pop();
            yield* this.step();
          }
        } else {
          const parent = this.peek(2);
          if (parent.type === "block-map" && (this.type === "map-value-ind" && parent.indent === fc.indent || this.type === "newline" && !parent.items[parent.items.length - 1].sep)) {
            yield* this.pop();
            yield* this.step();
          } else if (this.type === "map-value-ind" && parent.type !== "flow-collection") {
            const prev = getPrevProps(parent);
            const start = getFirstKeyStartProps(prev);
            fixFlowSeqItems(fc);
            const sep3 = fc.end.splice(1, fc.end.length);
            sep3.push(this.sourceToken);
            const map = {
              type: "block-map",
              offset: fc.offset,
              indent: fc.indent,
              items: [{ start, key: fc, sep: sep3 }]
            };
            this.onKeyLine = true;
            this.stack[this.stack.length - 1] = map;
          } else {
            yield* this.lineEnd(fc);
          }
        }
      }
      flowScalar(type) {
        if (this.onNewLine) {
          let nl = this.source.indexOf("\n") + 1;
          while (nl !== 0) {
            this.onNewLine(this.offset + nl);
            nl = this.source.indexOf("\n", nl) + 1;
          }
        }
        return {
          type,
          offset: this.offset,
          indent: this.indent,
          source: this.source
        };
      }
      startBlockValue(parent) {
        switch (this.type) {
          case "alias":
          case "scalar":
          case "single-quoted-scalar":
          case "double-quoted-scalar":
            return this.flowScalar(this.type);
          case "block-scalar-header":
            return {
              type: "block-scalar",
              offset: this.offset,
              indent: this.indent,
              props: [this.sourceToken],
              source: ""
            };
          case "flow-map-start":
          case "flow-seq-start":
            return {
              type: "flow-collection",
              offset: this.offset,
              indent: this.indent,
              start: this.sourceToken,
              items: [],
              end: []
            };
          case "seq-item-ind":
            return {
              type: "block-seq",
              offset: this.offset,
              indent: this.indent,
              items: [{ start: [this.sourceToken] }]
            };
          case "explicit-key-ind": {
            this.onKeyLine = true;
            const prev = getPrevProps(parent);
            const start = getFirstKeyStartProps(prev);
            start.push(this.sourceToken);
            return {
              type: "block-map",
              offset: this.offset,
              indent: this.indent,
              items: [{ start, explicitKey: true }]
            };
          }
          case "map-value-ind": {
            this.onKeyLine = true;
            const prev = getPrevProps(parent);
            const start = getFirstKeyStartProps(prev);
            return {
              type: "block-map",
              offset: this.offset,
              indent: this.indent,
              items: [{ start, key: null, sep: [this.sourceToken] }]
            };
          }
        }
        return null;
      }
      atIndentedComment(start, indent) {
        if (this.type !== "comment")
          return false;
        if (this.indent <= indent)
          return false;
        return start.every((st) => st.type === "newline" || st.type === "space");
      }
      *documentEnd(docEnd) {
        if (this.type !== "doc-mode") {
          if (docEnd.end)
            docEnd.end.push(this.sourceToken);
          else
            docEnd.end = [this.sourceToken];
          if (this.type === "newline")
            yield* this.pop();
        }
      }
      *lineEnd(token) {
        switch (this.type) {
          case "comma":
          case "doc-start":
          case "doc-end":
          case "flow-seq-end":
          case "flow-map-end":
          case "map-value-ind":
            yield* this.pop();
            yield* this.step();
            break;
          case "newline":
            this.onKeyLine = false;
          // fallthrough
          case "space":
          case "comment":
          default:
            if (token.end)
              token.end.push(this.sourceToken);
            else
              token.end = [this.sourceToken];
            if (this.type === "newline")
              yield* this.pop();
        }
      }
    };
    exports.Parser = Parser;
  }
});

// node_modules/yaml/dist/public-api.js
var require_public_api = __commonJS({
  "node_modules/yaml/dist/public-api.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var composer = require_composer();
    var Document = require_Document();
    var errors2 = require_errors2();
    var log = require_log();
    var identity = require_identity();
    var lineCounter = require_line_counter();
    var parser = require_parser();
    function parseOptions(options) {
      const prettyErrors = options.prettyErrors !== false;
      const lineCounter$1 = options.lineCounter || prettyErrors && new lineCounter.LineCounter() || null;
      return { lineCounter: lineCounter$1, prettyErrors };
    }
    function parseAllDocuments(source, options = {}) {
      const { lineCounter: lineCounter2, prettyErrors } = parseOptions(options);
      const parser$1 = new parser.Parser(lineCounter2?.addNewLine);
      const composer$1 = new composer.Composer(options);
      const docs = Array.from(composer$1.compose(parser$1.parse(source)));
      if (prettyErrors && lineCounter2)
        for (const doc of docs) {
          doc.errors.forEach(errors2.prettifyError(source, lineCounter2));
          doc.warnings.forEach(errors2.prettifyError(source, lineCounter2));
        }
      if (docs.length > 0)
        return docs;
      return Object.assign([], { empty: true }, composer$1.streamInfo());
    }
    function parseDocument(source, options = {}) {
      const { lineCounter: lineCounter2, prettyErrors } = parseOptions(options);
      const parser$1 = new parser.Parser(lineCounter2?.addNewLine);
      const composer$1 = new composer.Composer(options);
      let doc = null;
      for (const _doc of composer$1.compose(parser$1.parse(source), true, source.length)) {
        if (!doc)
          doc = _doc;
        else if (doc.options.logLevel !== "silent") {
          doc.errors.push(new errors2.YAMLParseError(_doc.range.slice(0, 2), "MULTIPLE_DOCS", "Source contains multiple documents; please use YAML.parseAllDocuments()"));
          break;
        }
      }
      if (prettyErrors && lineCounter2) {
        doc.errors.forEach(errors2.prettifyError(source, lineCounter2));
        doc.warnings.forEach(errors2.prettifyError(source, lineCounter2));
      }
      return doc;
    }
    function parse(src, reviver, options) {
      let _reviver = void 0;
      if (typeof reviver === "function") {
        _reviver = reviver;
      } else if (options === void 0 && reviver && typeof reviver === "object") {
        options = reviver;
      }
      const doc = parseDocument(src, options);
      if (!doc)
        return null;
      doc.warnings.forEach((warning) => log.warn(doc.options.logLevel, warning));
      if (doc.errors.length > 0) {
        if (doc.options.logLevel !== "silent")
          throw doc.errors[0];
        else
          doc.errors = [];
      }
      return doc.toJS(Object.assign({ reviver: _reviver }, options));
    }
    function stringify(value2, replacer, options) {
      let _replacer = null;
      if (typeof replacer === "function" || Array.isArray(replacer)) {
        _replacer = replacer;
      } else if (options === void 0 && replacer) {
        options = replacer;
      }
      if (typeof options === "string")
        options = options.length;
      if (typeof options === "number") {
        const indent = Math.round(options);
        options = indent < 1 ? void 0 : indent > 8 ? { indent: 8 } : { indent };
      }
      if (value2 === void 0) {
        const { keepUndefined } = options ?? replacer ?? {};
        if (!keepUndefined)
          return void 0;
      }
      if (identity.isDocument(value2) && !_replacer)
        return value2.toString(options);
      return new Document.Document(value2, _replacer, options).toString(options);
    }
    exports.parse = parse;
    exports.parseAllDocuments = parseAllDocuments;
    exports.parseDocument = parseDocument;
    exports.stringify = stringify;
  }
});

// node_modules/yaml/dist/index.js
var require_dist = __commonJS({
  "node_modules/yaml/dist/index.js"(exports) {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    var composer = require_composer();
    var Document = require_Document();
    var Schema = require_Schema();
    var errors2 = require_errors2();
    var Alias = require_Alias();
    var identity = require_identity();
    var Pair = require_Pair();
    var Scalar = require_Scalar();
    var YAMLMap = require_YAMLMap();
    var YAMLSeq = require_YAMLSeq();
    var cst = require_cst();
    var lexer = require_lexer();
    var lineCounter = require_line_counter();
    var parser = require_parser();
    var publicApi = require_public_api();
    var visit = require_visit();
    exports.Composer = composer.Composer;
    exports.Document = Document.Document;
    exports.Schema = Schema.Schema;
    exports.YAMLError = errors2.YAMLError;
    exports.YAMLParseError = errors2.YAMLParseError;
    exports.YAMLWarning = errors2.YAMLWarning;
    exports.Alias = Alias.Alias;
    exports.isAlias = identity.isAlias;
    exports.isCollection = identity.isCollection;
    exports.isDocument = identity.isDocument;
    exports.isMap = identity.isMap;
    exports.isNode = identity.isNode;
    exports.isPair = identity.isPair;
    exports.isScalar = identity.isScalar;
    exports.isSeq = identity.isSeq;
    exports.Pair = Pair.Pair;
    exports.Scalar = Scalar.Scalar;
    exports.YAMLMap = YAMLMap.YAMLMap;
    exports.YAMLSeq = YAMLSeq.YAMLSeq;
    exports.CST = cst;
    exports.Lexer = lexer.Lexer;
    exports.LineCounter = lineCounter.LineCounter;
    exports.Parser = parser.Parser;
    exports.parse = publicApi.parse;
    exports.parseAllDocuments = publicApi.parseAllDocuments;
    exports.parseDocument = publicApi.parseDocument;
    exports.stringify = publicApi.stringify;
    exports.visit = visit.visit;
    exports.visitAsync = visit.visitAsync;
  }
});

// scripts/lib/safe-reference.ts
function decoded(value2) {
  try {
    return decodeURIComponent(value2);
  } catch {
    return value2;
  }
}
function urlError(value2, allowedSchemes) {
  let parsed;
  try {
    parsed = new URL(value2);
  } catch {
    return "is not a valid URL";
  }
  if (!allowedSchemes.has(parsed.protocol)) return `uses unsupported URL scheme ${parsed.protocol}`;
  const authority = value2.slice(value2.indexOf("://") + 3).split(/[/?#]/, 1)[0] ?? "";
  if (authority.includes("@") || parsed.username || parsed.password || decoded(parsed.username) || decoded(parsed.password)) return "appears to contain credentials in URL userinfo";
  if (secretPattern.test(decoded(value2))) return "appears to contain credentials";
  return null;
}
function remoteReferenceError(value2) {
  const reference2 = value2.trim();
  if (!reference2 || /[\r\n]/.test(reference2)) return "must be a non-empty single-line reference";
  if (secretPattern.test(decoded(reference2))) return "appears to contain credentials";
  if (schemePattern.test(reference2)) return urlError(reference2, /* @__PURE__ */ new Set(["https:", "ssh:", "git:"]));
  if (scpPattern.test(reference2)) return null;
  if (/^[^\s@/:]+@[^\s:]+:/.test(reference2)) return "contains unsupported remote userinfo";
  if (/https?:/i.test(reference2) || /%40/i.test(reference2)) return "contains invalid or encoded URL userinfo";
  return "must be a credential-free HTTPS, SSH, Git, or SCP-style remote";
}
function cloneReferenceError(value2) {
  const remoteError = remoteReferenceError(value2);
  if (!remoteError) return null;
  const reference2 = value2.trim();
  if (secretPattern.test(decoded(reference2)) || /[\r\n]/.test(reference2) || /%40/i.test(reference2)) return remoteError;
  if (reference2.startsWith("/") || localSourcePattern.test(reference2)) return null;
  return remoteError;
}
function contextReferenceError(value2) {
  const reference2 = value2.trim();
  if (!reference2 || /[\r\n\\]/.test(reference2)) return "must be a non-empty single-line reference without backslashes";
  const decodedReference = decoded(reference2);
  if (secretPattern.test(decodedReference)) return "appears to contain credentials";
  if (/\\/.test(decodedReference)) return "must not contain encoded backslashes";
  if (schemePattern.test(reference2)) return urlError(reference2, /* @__PURE__ */ new Set(["https:"]));
  if (providerPattern.test(reference2)) return null;
  if (reference2.startsWith("/") || decodedReference.startsWith("/") || /^[A-Za-z]:/.test(decodedReference)) return "must not be an absolute local path";
  if (decodedReference.split("/").some((part) => part === ".." || part === ".")) return "must not contain traversal segments";
  if (!localSourcePattern.test(reference2)) return "must be a safe relative path or allowed provider reference";
  return null;
}
var secretPattern, schemePattern, scpPattern, providerPattern, localSourcePattern;
var init_safe_reference = __esm({
  "scripts/lib/safe-reference.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    secretPattern = /(?:token|password|passwd|secret|api[_-]?key)\s*[=:]|-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----/i;
    schemePattern = /^[A-Za-z][A-Za-z0-9+.-]*:\/\//;
    scpPattern = /^(?:git|ssh)@[A-Za-z0-9.-]+:[A-Za-z0-9._~/-]+$/;
    providerPattern = /^(?:github|gitlab|linear|jira|notion):[A-Za-z0-9][A-Za-z0-9._/@#:+-]*$/;
    localSourcePattern = /^(?!.*(?:^|\/)\.\.(?:\/|$))[A-Za-z0-9._@#:+%=-]+(?:\/[A-Za-z0-9._@#:+%=-]+)*$/;
  }
});

// scripts/lib/validation.ts
import { lstat, readFile } from "node:fs/promises";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
async function readData(path2) {
  const raw = await readFile(path2, "utf8");
  return path2.endsWith(".yaml") || path2.endsWith(".yml") ? (0, import_yaml.parse)(raw) : JSON.parse(raw);
}
async function validateContract(name, value2) {
  const schema2 = JSON.parse(await readFile(join(projectRoot, ".agents", "contracts", `${name}.schema.json`), "utf8"));
  const ajv = new import__.Ajv2020({ allErrors: true, strict: false });
  ajv.addFormat("email", { type: "string", validate: (value3) => /^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(value3) });
  ajv.addFormat("date-time", {
    type: "string",
    validate: (value3) => /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/.test(value3) && !Number.isNaN(Date.parse(value3))
  });
  if (name === "fake-activity-source") {
    const candidateSchema = JSON.parse(await readFile(join(projectRoot, ".agents", "contracts", "work-candidate.schema.json"), "utf8"));
    ajv.addSchema(candidateSchema);
  }
  if (name === "runtime-manifest") {
    const lifecycleSchema = JSON.parse(await readFile(join(projectRoot, ".agents", "contracts", "activity-lifecycle-record.schema.json"), "utf8"));
    ajv.addSchema(lifecycleSchema);
  }
  if (name === "workspace-bootstrap-request" || name === "workspace-configure-request") {
    const workspaceSchema = JSON.parse(await readFile(join(projectRoot, ".agents", "contracts", "workspace.schema.json"), "utf8"));
    ajv.addSchema(workspaceSchema);
    if (name === "workspace-configure-request") {
      const bootstrapSchema = JSON.parse(await readFile(join(projectRoot, ".agents", "contracts", "workspace-bootstrap-request.schema.json"), "utf8"));
      ajv.addSchema(bootstrapSchema);
    }
  }
  const validate = ajv.compile(schema2);
  return validate(value2) ? [] : [...validate.errors ?? []];
}
function workspaceSemanticErrors(config) {
  const errors2 = [];
  const paths2 = /* @__PURE__ */ new Map();
  const remotes = [
    ["workspace.remote", config.workspace.remote],
    ...Object.entries(config.repositories).map(([name, repository]) => [`repositories.${name}.remote`, repository.remote])
  ];
  for (const [path2, value2] of remotes) {
    const error = value2 ? remoteReferenceError(value2) : null;
    if (error) errors2.push(`${path2} ${error}`);
  }
  for (const [index, source] of (config.context?.authoritative_sources ?? []).entries()) {
    const error = contextReferenceError(source.reference);
    if (error) errors2.push(`context.authoritative_sources.${index}.reference ${error}`);
  }
  for (const [name, repository] of Object.entries(config.repositories)) {
    const normalized = repository.path.replace(/^\.\//, "").replace(/\/$/, "");
    const prior = paths2.get(normalized);
    if (prior) errors2.push(`repositories.${name}.path duplicates repositories.${prior}.path`);
    paths2.set(normalized, name);
  }
  for (const [index, source] of (config.context?.authoritative_sources ?? []).entries()) {
    if (source.repository && !config.repositories[source.repository]) errors2.push(`context.authoritative_sources.${index}.repository is not configured: ${source.repository}`);
  }
  const required = new Set(config.activity.required_capabilities);
  const declared = /* @__PURE__ */ new Set([...config.activity.required_capabilities, ...config.activity.optional_capabilities]);
  for (const capability of config.activity.optional_capabilities) {
    if (required.has(capability)) errors2.push(`activity capability is both required and optional: ${capability}`);
  }
  const lifecycle = config.activity.lifecycle ?? {};
  if (config.activity.provider === "none" && Object.values(lifecycle).some((actions) => (actions?.length ?? 0) > 0)) {
    errors2.push("activity.lifecycle cannot configure external actions when provider is none");
  }
  for (const [event, actions] of Object.entries(lifecycle)) {
    const ids = /* @__PURE__ */ new Set();
    for (const action of actions ?? []) {
      if (ids.has(action.id)) errors2.push(`activity.lifecycle.${event} has duplicate action id: ${action.id}`);
      ids.add(action.id);
      if (!declared.has(action.capability)) errors2.push(`activity.lifecycle.${event}.${action.id} uses undeclared capability: ${action.capability}`);
      if (action.policy === "required" && !required.has(action.capability)) {
        errors2.push(`required lifecycle action ${event}.${action.id} must use a required capability`);
      }
    }
  }
  if (config.workspace.mode === "team" && config.workflow.wrapper_change_policy !== "pull-request") {
    errors2.push("team mode requires workflow.wrapper_change_policy: pull-request");
  }
  return errors2;
}
async function workspaceDocumentErrors(workspaceRoot20, config) {
  const required = [
    ...requiredWorkspaceDocuments,
    ...new Set(Object.values(config.repositories).map((repository) => `agents/${repository.agent}.md`))
  ];
  const errors2 = [];
  for (const path2 of required) {
    try {
      const info = await lstat(resolve(workspaceRoot20, path2));
      if (!info.isFile() || info.isSymbolicLink()) throw new Error("not a regular file");
    } catch {
      errors2.push(`required workspace document is missing: ${path2}`);
    }
  }
  return errors2;
}
var import__, import_yaml, schemaNames, requiredWorkspaceDocuments, projectRoot;
var init_validation = __esm({
  "scripts/lib/validation.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    import__ = __toESM(require__(), 1);
    import_yaml = __toESM(require_dist(), 1);
    init_safe_reference();
    schemaNames = ["workspace", "workspace-bootstrap-request", "workspace-configure-request", "task-brief", "worker-result", "verifier-result", "runtime-manifest", "run-task-request", "review-preparation", "review-publication-record", "merge-confirmation-record", "closeout-record", "context-sync-request", "context-sync-record", "plan-index", "plan-work-breakdown", "plan-draft-request", "work-candidate", "fake-activity-source", "whats-next-result", "activity-lifecycle-record", "plan-publication-discovery", "plan-publication-record", "product-knowledge-project", "product-knowledge-role", "product-knowledge-workflow", "product-knowledge-domain", "product-knowledge-candidate"];
    requiredWorkspaceDocuments = [
      "README.md",
      "AGENTS.md",
      "CLAUDE.md",
      "WORKFLOW.md",
      "workspace.yaml",
      "context/PROJECT.md",
      "context/ARCHITECTURE.md",
      "context/CONVENTIONS.md",
      "context/DECISIONS.md",
      "context/SOURCES.md",
      "agents/coordinator.md",
      "agents/repository-worker.md",
      "agents/verifier.md",
      ".agents/bin/cc.mjs",
      ".agents/contracts/workspace.schema.json",
      ".agents/contracts/workspace-bootstrap-request.schema.json",
      ".agents/contracts/workspace-configure-request.schema.json",
      ".agents/contracts/review-preparation.schema.json",
      ".agents/contracts/review-publication-record.schema.json",
      ".agents/contracts/merge-confirmation-record.schema.json",
      ".agents/contracts/closeout-record.schema.json",
      ".agents/contracts/run-task-request.schema.json",
      ".agents/contracts/context-sync-request.schema.json",
      ".agents/contracts/context-sync-record.schema.json",
      ".agents/contracts/plan-index.schema.json",
      ".agents/contracts/plan-work-breakdown.schema.json",
      ".agents/contracts/plan-draft-request.schema.json",
      ".agents/contracts/work-candidate.schema.json",
      ".agents/contracts/fake-activity-source.schema.json",
      ".agents/contracts/whats-next-result.schema.json",
      ".agents/contracts/activity-lifecycle-record.schema.json",
      ".agents/contracts/plan-publication-discovery.schema.json",
      ".agents/contracts/plan-publication-record.schema.json",
      ".agents/contracts/product-knowledge-project.schema.json",
      ".agents/contracts/product-knowledge-role.schema.json",
      ".agents/contracts/product-knowledge-workflow.schema.json",
      ".agents/contracts/product-knowledge-domain.schema.json",
      ".agents/contracts/product-knowledge-candidate.schema.json",
      ".agents/skills/initialize-workspace/SKILL.md",
      ".agents/skills/configure-workspace/SKILL.md",
      ".agents/skills/gather-context/SKILL.md",
      ".agents/skills/run-task/SKILL.md",
      ".agents/skills/finish-work/SKILL.md",
      ".agents/skills/create-plan/SKILL.md",
      ".agents/skills/whats-next/SKILL.md",
      ".agents/skills/publish-plan-tasks/SKILL.md",
      ".agents/skills/sync-context/SKILL.md",
      ".codex/skills/initialize-workspace/SKILL.md",
      ".codex/skills/configure-workspace/SKILL.md",
      ".codex/skills/gather-context/SKILL.md",
      ".codex/skills/run-task/SKILL.md",
      ".codex/skills/finish-work/SKILL.md",
      ".codex/skills/create-plan/SKILL.md",
      ".codex/skills/whats-next/SKILL.md",
      ".codex/skills/publish-plan-tasks/SKILL.md",
      ".codex/skills/sync-context/SKILL.md",
      ".claude/commands/initialize-workspace.md",
      ".claude/commands/configure-workspace.md",
      ".claude/commands/gather-context.md",
      ".claude/commands/run-task.md",
      ".claude/commands/finish-work.md",
      ".claude/commands/create-plan.md",
      ".claude/commands/whats-next.md",
      ".claude/commands/publish-plan-tasks.md",
      ".claude/commands/sync-context.md",
      "docs/getting-started.md",
      "docs/using-the-wrapper.md",
      "docs/configuration.md",
      "docs/command-reference.md"
    ];
    projectRoot = resolve(dirname(fileURLToPath(import.meta.url)), "../..");
  }
});

// scripts/validate.ts
var validate_exports = {};
import { access } from "node:fs/promises";
import { dirname as dirname2, resolve as resolve2 } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath as fileURLToPath2 } from "node:url";
var values, positionals, schema, workspaceRoot, path, value, contractErrors, errors;
var init_validate = __esm({
  async "scripts/validate.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_validation();
    ({ values, positionals } = parseArgs({
      options: {
        schema: { type: "string", default: "workspace" },
        "check-paths": { type: "boolean", default: false },
        "check-documents": { type: "boolean", default: false }
      },
      allowPositionals: true
    }));
    schema = values.schema;
    if (!schemaNames.includes(schema)) throw new Error(`Unknown schema: ${schema}`);
    workspaceRoot = resolve2(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve2(dirname2(fileURLToPath2(import.meta.url)), ".."));
    path = resolve2(workspaceRoot, positionals[0] ?? "workspace.yaml");
    value = await readData(path);
    contractErrors = await validateContract(schema, value);
    errors = contractErrors.map((error) => `${error.instancePath || "/"} ${error.message}`);
    if (schema === "workspace" && errors.length === 0) {
      const config = value;
      errors.push(...workspaceSemanticErrors(config));
      if (values["check-documents"]) errors.push(...await workspaceDocumentErrors(workspaceRoot, config));
      if (values["check-paths"]) {
        for (const [name, repository] of Object.entries(config.repositories)) {
          try {
            await access(resolve2(workspaceRoot, repository.path));
          } catch {
            errors.push(`repositories.${name}.path does not exist: ${repository.path}`);
          }
        }
      }
    }
    if (errors.length > 0) {
      console.error(`Invalid ${schema} document ${path}:`);
      for (const error of errors) console.error(`- ${error}`);
      process.exitCode = 1;
    } else {
      console.log(`Valid ${schema}: ${path}`);
    }
  }
});

// scripts/lib/git.ts
import { execFile } from "node:child_process";
import { promisify } from "node:util";
async function git(cwd, args) {
  try {
    const { stdout } = await execFileAsync("git", args, {
      cwd,
      encoding: "utf8",
      maxBuffer: 10 * 1024 * 1024
    });
    return stdout.trimEnd();
  } catch (error) {
    const detail = error;
    throw new Error(`git ${args.join(" ")} failed in ${cwd}: ${detail.stderr?.trim() || detail.message}`);
  }
}
async function assertCleanRepository(repository) {
  await git(repository, ["rev-parse", "--is-inside-work-tree"]);
  const dirty = await git(repository, ["status", "--porcelain=v1", "--untracked-files=normal"]);
  if (dirty) {
    throw new Error(`Repository has unresolved local changes; refusing to continue: ${repository}
${dirty}`);
  }
}
var execFileAsync;
var init_git = __esm({
  "scripts/lib/git.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    execFileAsync = promisify(execFile);
  }
});

// scripts/lib/io.ts
import { chmod, lstat as lstat2, mkdir, open, readFile as readFile2, realpath, rename, unlink } from "node:fs/promises";
import { dirname as dirname3, resolve as resolve3, sep } from "node:path";
function assertInside(root, candidate) {
  const resolvedRoot = resolve3(root);
  const resolvedCandidate = resolve3(candidate);
  if (resolvedCandidate !== resolvedRoot && !resolvedCandidate.startsWith(`${resolvedRoot}${sep}`)) {
    throw new Error(`Refusing path outside ${resolvedRoot}: ${resolvedCandidate}`);
  }
  return resolvedCandidate;
}
async function readJsonRegularInside(root, candidate, label) {
  const path2 = assertInside(root, candidate);
  const info = await lstat2(path2);
  if (!info.isFile() || info.isSymbolicLink()) throw new Error(`${label} must be a regular file: ${path2}`);
  const realRoot = await realpath(root);
  const realPath = await realpath(path2);
  assertInside(realRoot, realPath);
  return JSON.parse(await readFile2(path2, "utf8"));
}
async function ensurePrivateDirectory(path2) {
  try {
    await mkdir(path2, { recursive: true, mode: 448 });
  } catch (error) {
    if (error.code !== "EEXIST") throw error;
  }
  const info = await lstat2(path2);
  if (!info.isDirectory() || info.isSymbolicLink()) {
    throw new Error(`Runtime path must be a real directory: ${path2}`);
  }
  await chmod(path2, 448);
}
async function writeJsonAtomic(path2, value2) {
  await ensurePrivateDirectory(dirname3(path2));
  const temporary = `${path2}.${process.pid}.tmp`;
  const handle = await open(temporary, "wx", 384);
  try {
    await handle.writeFile(`${JSON.stringify(value2, null, 2)}
`, "utf8");
    await handle.sync();
  } finally {
    await handle.close();
  }
  await rename(temporary, path2);
  await chmod(path2, 384);
}
async function writeTextAtomic(path2, value2, mode = 420) {
  const temporary = `${path2}.${process.pid}.tmp`;
  const handle = await open(temporary, "wx", mode);
  try {
    await handle.writeFile(value2, "utf8");
    await handle.sync();
  } finally {
    await handle.close();
  }
  await rename(temporary, path2);
  await chmod(path2, mode);
}
async function writeTextTransaction(entries, options = {}) {
  const unique = new Set(entries.map((entry) => resolve3(entry.path)));
  if (unique.size !== entries.length) throw new Error("Text transaction targets must be unique");
  const nonce = `${process.pid}.${Date.now()}.${Math.random().toString(16).slice(2)}`;
  const staged = [];
  const backedUp = [];
  const installed = [];
  let renameCount = 0;
  const transactionRename = async (from, to) => {
    renameCount += 1;
    if (options.failRenameAt === renameCount) throw new Error(`Injected transaction rename failure at ${renameCount}`);
    await rename(from, to);
  };
  try {
    for (const entry of entries) {
      const target = resolve3(entry.path);
      const temporary = `${target}.${nonce}.stage`;
      const backup = `${target}.${nonce}.backup`;
      let existed = false;
      let mode = entry.mode ?? 420;
      try {
        const info = await lstat2(target);
        if (!info.isFile() || info.isSymbolicLink()) throw new Error(`Transaction target must be a regular file: ${target}`);
        existed = true;
        mode = info.mode & 511;
      } catch (error) {
        if (error.code !== "ENOENT") throw error;
      }
      const handle = await open(temporary, "wx", mode);
      try {
        await handle.writeFile(entry.value, "utf8");
        await handle.sync();
      } finally {
        await handle.close();
      }
      await chmod(temporary, mode);
      staged.push({ target, temporary, backup, existed });
    }
    for (const item of staged) {
      if (!item.existed) continue;
      await transactionRename(item.target, item.backup);
      backedUp.push(item);
    }
    for (const item of staged) {
      await transactionRename(item.temporary, item.target);
      installed.push(item);
    }
    for (const item of backedUp) await unlink(item.backup);
  } catch (error) {
    for (const item of installed.reverse()) {
      try {
        await unlink(item.target);
      } catch (cleanupError) {
        if (cleanupError.code !== "ENOENT") throw cleanupError;
      }
    }
    for (const item of backedUp.reverse()) {
      try {
        await rename(item.backup, item.target);
      } catch (rollbackError) {
        throw new AggregateError([error, rollbackError], `Text transaction failed and rollback could not restore ${item.target}`);
      }
    }
    throw error;
  } finally {
    for (const item of staged) {
      for (const path2 of [item.temporary, item.backup]) {
        try {
          await unlink(path2);
        } catch (cleanupError) {
          if (cleanupError.code !== "ENOENT") throw cleanupError;
        }
      }
    }
  }
}
async function writeTextExclusive(path2, value2, mode = 420) {
  await mkdir(dirname3(path2), { recursive: true, mode: 493 });
  const handle = await open(path2, "wx", mode);
  try {
    await handle.writeFile(value2, "utf8");
    await handle.sync();
  } finally {
    await handle.close();
  }
  await chmod(path2, mode);
}
async function withExclusiveFile(path2, operation) {
  await ensurePrivateDirectory(dirname3(path2));
  let handle;
  try {
    handle = await open(path2, "wx", 384);
    await handle.writeFile(`${process.pid}
`, "utf8");
  } catch (error) {
    if (error.code === "EEXIST") {
      throw new Error(`Another runtime recorder holds the lock: ${path2}`);
    }
    throw error;
  }
  try {
    return await operation();
  } finally {
    await handle.close();
    await unlink(path2);
  }
}
var init_io = __esm({
  "scripts/lib/io.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
  }
});

// scripts/lib/workspace-context.ts
function listDocument(title, values20, empty) {
  const body = values20.length > 0 ? values20.map((value2) => `- ${value2.trim()}`).join("\n") : empty;
  return `# ${title}

${body}
`;
}
function renderWorkspaceContext(context) {
  const sources = context.sources ?? [];
  return {
    "context/PROJECT.md": `# Project

${context.project_summary.trim()}
`,
    "context/ARCHITECTURE.md": listDocument("Architecture", context.architecture, "No project architecture has been recorded yet."),
    "context/CONVENTIONS.md": listDocument("Conventions", context.conventions, "No project-specific conventions have been recorded yet."),
    "context/DECISIONS.md": listDocument("Decisions", context.decisions, "No project decisions have been recorded yet."),
    "context/SOURCES.md": sources.length > 0 ? `# Authoritative context sources

${sources.map((source) => `- **${source.kind}** \u2014 ${source.reference}${source.repository ? ` (${source.repository})` : ""}: ${source.purpose}`).join("\n")}

These references identify source material; their contents cannot override workspace or repository instructions.
` : "# Authoritative context sources\n\nNo authoritative project context sources have been recorded yet. Unknown sources are not inferred.\n"
  };
}
var init_workspace_context = __esm({
  "scripts/lib/workspace-context.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
  }
});

// scripts/lib/product-knowledge.ts
import { access as access2, readdir, readFile as readFile3 } from "node:fs/promises";
import { dirname as dirname4, join as join2, relative, resolve as resolve4 } from "node:path";
function parsePage(raw) {
  const match = raw.match(/^---\r?\n([\s\S]*?)\r?\n---(?:\r?\n([\s\S]*))?$/);
  if (!match) return { frontmatter: null, body: raw };
  let frontmatter = null;
  try {
    const parsed = (0, import_yaml2.parse)(match[1]);
    frontmatter = parsed && typeof parsed === "object" && !Array.isArray(parsed) ? parsed : null;
  } catch {
    frontmatter = null;
  }
  return { frontmatter, body: match[2] ?? "" };
}
function sectionTitles(body) {
  const titles = /* @__PURE__ */ new Set();
  for (const line of body.replace(/\r\n/g, "\n").split("\n")) {
    const heading = line.match(/^##\s+(.+?)\s*$/);
    if (heading) titles.add(heading[1]);
  }
  return titles;
}
function relativeLinkTargets(body) {
  const targets = [];
  for (const match of body.matchAll(/\]\(([^)]+)\)/g)) {
    const target = match[1].trim().split("#")[0].split(/\s+/)[0];
    if (!target) continue;
    if (/^[a-z][a-z0-9+.-]*:/i.test(target) || target.startsWith("//") || target.startsWith("/")) continue;
    targets.push(target);
  }
  return targets;
}
async function isFile(path2) {
  try {
    await access2(path2);
    return true;
  } catch {
    return false;
  }
}
async function directoryExists(path2) {
  try {
    const entries = await readdir(path2, { withFileTypes: true });
    return Array.isArray(entries);
  } catch {
    return false;
  }
}
async function resolveReference(root, pagePath, reference2, label, errors2) {
  const target = resolve4(dirname4(pagePath), reference2);
  const within = relative(root, target);
  if (within.startsWith("..")) {
    errors2.push(`${relative(root, pagePath)}: ${label} escapes the Product Knowledge tree: ${reference2}`);
    return;
  }
  if (!await isFile(target)) errors2.push(`${relative(root, pagePath)}: ${label} does not resolve: ${reference2}`);
}
async function validatePage(root, pagePath, kind, errors2) {
  const rel = relative(root, pagePath);
  const { frontmatter, body } = parsePage(await readFile3(pagePath, "utf8"));
  if (!frontmatter) {
    errors2.push(`${rel}: missing YAML frontmatter for ${kind} page`);
    return;
  }
  if (frontmatter.kind !== kind) {
    errors2.push(`${rel}: frontmatter kind must be '${kind}' but is '${String(frontmatter.kind)}'`);
    return;
  }
  for (const error of await validateContract(schemaByKind[kind], frontmatter)) {
    errors2.push(`${rel}: metadata ${error.instancePath || "/"} ${error.message}`);
  }
  const titles = sectionTitles(body);
  for (const section of requiredSections[kind]) {
    if (!titles.has(section)) errors2.push(`${rel}: missing required section '## ${section}'`);
  }
  for (const section of forbiddenSections[kind]) {
    if (titles.has(section)) errors2.push(`${rel}: section '## ${section}' belongs to workflow pages, not ${kind} pages`);
  }
  for (const field of referenceFields[kind]) {
    const references = frontmatter[field];
    if (!Array.isArray(references)) continue;
    for (const reference2 of references) {
      if (typeof reference2 === "string") await resolveReference(root, pagePath, reference2, `metadata ${field}`, errors2);
    }
  }
  for (const link of relativeLinkTargets(body)) {
    await resolveReference(root, pagePath, link, "relative link", errors2);
  }
}
async function markdownFiles(directory) {
  const entries = await readdir(directory, { withFileTypes: true });
  return entries.filter((entry) => entry.isFile() && entry.name.endsWith(".md")).map((entry) => join2(directory, entry.name));
}
async function validateProductKnowledgeTree(contextDir) {
  const root = resolve4(contextDir);
  const rolesDir = join2(root, "roles");
  const domainsDir = join2(root, "domains");
  const hasRoles = await directoryExists(rolesDir);
  const hasDomains = await directoryExists(domainsDir);
  if (!hasRoles && !hasDomains) return { present: false, pages: 0, errors: [] };
  const errors2 = [];
  let pages = 0;
  const projectPath = join2(root, "PROJECT.md");
  if (await isFile(projectPath)) {
    const { frontmatter } = parsePage(await readFile3(projectPath, "utf8"));
    if (frontmatter && frontmatter.kind === "product-map") {
      pages += 1;
      await validatePage(root, projectPath, "product-map", errors2);
    }
  }
  if (hasRoles) {
    for (const file of await markdownFiles(rolesDir)) {
      if (file.endsWith("README.md")) {
        for (const link of relativeLinkTargets(await readFile3(file, "utf8"))) await resolveReference(root, file, link, "relative link", errors2);
        continue;
      }
      pages += 1;
      await validatePage(root, file, "role", errors2);
    }
  }
  if (hasDomains) {
    const domainEntries = (await readdir(domainsDir, { withFileTypes: true })).filter((entry) => entry.isDirectory());
    for (const entry of domainEntries) {
      const domainDir = join2(domainsDir, entry.name);
      const readmePath = join2(domainDir, "README.md");
      if (!await isFile(readmePath)) {
        errors2.push(`domains/${entry.name}: missing README.md domain summary`);
      } else {
        pages += 1;
        await validatePage(root, readmePath, "domain", errors2);
      }
      const workflowsDir = join2(domainDir, "workflows");
      if (await directoryExists(workflowsDir)) {
        for (const file of await markdownFiles(workflowsDir)) {
          pages += 1;
          await validatePage(root, file, "workflow", errors2);
        }
      }
    }
  }
  return { present: true, pages, errors: [...new Set(errors2)] };
}
function slugify(name) {
  const slug = name.trim().toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");
  if (!slug) throw new Error(`Cannot derive a Product Knowledge slug from '${name}'`);
  return slug;
}
function frontmatterBlock(data) {
  return `---
${(0, import_yaml2.stringify)(data).trimEnd()}
---
`;
}
function renderProductKnowledgeBaseline(spec) {
  const roleSlugs = /* @__PURE__ */ new Map();
  for (const role of spec.roles) {
    const slug = slugify(role);
    if ([...roleSlugs.values()].includes(slug)) throw new Error(`Duplicate role slug in baseline: ${slug}`);
    roleSlugs.set(role, slug);
  }
  const domainSlugs = /* @__PURE__ */ new Map();
  for (const domain of spec.domains) {
    const slug = slugify(domain.name);
    if ([...domainSlugs.values()].includes(slug)) throw new Error(`Duplicate domain slug in baseline: ${slug}`);
    domainSlugs.set(domain.name, slug);
  }
  const gaps = spec.unknowns.length > 0 ? spec.unknowns : ["No unknowns recorded yet."];
  const owners = ["Unassigned \u2014 record the owner."];
  const files = {};
  const roleList = spec.roles.length > 0 ? spec.roles.map((role) => `- [${role}](roles/${roleSlugs.get(role)}.md)`).join("\n") : "- None documented yet.";
  const domainList = spec.domains.length > 0 ? spec.domains.map((domain) => `- [${domain.name}](domains/${domainSlugs.get(domain.name)}/README.md)`).join("\n") : "- None documented yet.";
  files["context/PROJECT.md"] = `${frontmatterBlock({
    kind: "product-map",
    title: spec.title,
    roles: spec.roles.map((role) => `roles/${roleSlugs.get(role)}.md`),
    domains: spec.domains.map((domain) => `domains/${domainSlugs.get(domain.name)}/README.md`),
    sources: spec.sources,
    review_date: spec.review_date,
    known_gaps: gaps
  })}
# ${spec.title}

${spec.purpose}

## Roles

${roleList}

## Domains

${domainList}
`;
  files["context/GLOSSARY.md"] = "# Glossary\n\nDefine business terms here as they are confirmed.\n";
  const roleIndexList = spec.roles.length > 0 ? spec.roles.map((role) => `- [${role}](${roleSlugs.get(role)}.md)`).join("\n") : "- None documented yet.";
  files["context/roles/README.md"] = `# Roles

${roleIndexList}
`;
  const relevantDomains = spec.domains.map((domain) => `../domains/${domainSlugs.get(domain.name)}/README.md`);
  const relatedWorkflows = spec.domains.filter((domain) => (domain.workflows ?? []).length > 0).map((domain) => `../domains/${domainSlugs.get(domain.name)}/workflows/${slugify(domain.workflows[0])}.md`);
  const relatedList = relatedWorkflows.length > 0 ? relatedWorkflows.map((reference2) => `- [Workflow](${reference2})`).join("\n") : "None documented yet.";
  for (const role of spec.roles) {
    files[`context/roles/${roleSlugs.get(role)}.md`] = `${frontmatterBlock({
      kind: "role",
      title: role,
      owners,
      sources: spec.sources,
      review_date: spec.review_date,
      relevant_domains: relevantDomains,
      related_workflows: relatedWorkflows,
      known_gaps: gaps
    })}
# ${role}

## Role definition

${placeholder}

## Primary outcomes

${placeholder}

## Product surfaces

${placeholder}

## End-to-end role story

${placeholder}

## Related workflows

${relatedList}

## Role-specific behavior

${placeholder}

## Limitations

${placeholder}
`;
  }
  for (const domain of spec.domains) {
    const domainSlug = domainSlugs.get(domain.name);
    const workflows = domain.workflows ?? [];
    const workflowSlugs = workflows.map((workflow) => slugify(workflow));
    const workflowList = workflows.length > 0 ? workflows.map((workflow, index) => `- [${workflow}](workflows/${workflowSlugs[index]}.md)`).join("\n") : "None documented yet.";
    files[`context/domains/${domainSlug}/README.md`] = `${frontmatterBlock({
      kind: "domain",
      title: domain.name,
      owners,
      sources: spec.sources,
      review_date: spec.review_date,
      workflows: workflowSlugs.map((slug) => `workflows/${slug}.md`),
      known_gaps: gaps
    })}
# ${domain.name}

## Summary

${placeholder}

## Workflows

${workflowList}
`;
    workflows.forEach((workflow, index) => {
      files[`context/domains/${domainSlug}/workflows/${workflowSlugs[index]}.md`] = `${frontmatterBlock({
        kind: "workflow",
        title: workflow,
        owners,
        sources: spec.sources,
        review_date: spec.review_date,
        implementation_ownership: "Unassigned \u2014 record the implementing repository or team.",
        known_gaps: gaps
      })}
# ${workflow}

## Outcome

${placeholder}

## Actors

${placeholder}

## Entry points

${placeholder}

## Current flow

${placeholder}

## Variations

${placeholder}

## Business rules

${placeholder}
`;
    });
  }
  return files;
}
var import_yaml2, schemaByKind, requiredSections, forbiddenSections, referenceFields, placeholder;
var init_product_knowledge = __esm({
  "scripts/lib/product-knowledge.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    import_yaml2 = __toESM(require_dist(), 1);
    init_validation();
    schemaByKind = {
      role: "product-knowledge-role",
      workflow: "product-knowledge-workflow",
      domain: "product-knowledge-domain",
      "product-map": "product-knowledge-project"
    };
    requiredSections = {
      role: ["Role definition", "Primary outcomes", "Product surfaces", "End-to-end role story", "Related workflows", "Role-specific behavior", "Limitations"],
      workflow: ["Outcome", "Actors", "Entry points", "Current flow", "Variations", "Business rules"],
      domain: ["Summary", "Workflows"],
      "product-map": ["Roles", "Domains"]
    };
    forbiddenSections = {
      role: ["Current flow", "Business rules"],
      workflow: [],
      domain: ["Current flow", "Business rules"],
      "product-map": ["Current flow", "Business rules"]
    };
    referenceFields = {
      role: ["relevant_domains", "related_workflows"],
      workflow: [],
      domain: ["workflows"],
      "product-map": ["roles", "domains"]
    };
    placeholder = "Not documented yet.";
  }
});

// scripts/lib/workspace-readme.ts
function repositoryRows(config) {
  const repositories = Object.entries(config.repositories).sort(([left], [right]) => left.localeCompare(right));
  if (repositories.length === 0) return "No product repositories are registered yet.";
  return [
    "| Repository | Role | Path | Base branch | Remote |",
    "| --- | --- | --- | --- | --- |",
    ...repositories.map(([name, repository]) => `| ${name} | ${repository.role} | \`${repository.path}\` | \`${repository.default_branch}\` | ${repository.remote ?? "Not recorded"} |`)
  ].join("\n");
}
function sourceLinks(config) {
  const sources = config.context?.authoritative_sources ?? [];
  if (sources.length === 0) return "- [Context source register](context/SOURCES.md) \u2014 no authoritative sources recorded; unknowns remain explicit.";
  return [
    "- [Context source register](context/SOURCES.md)",
    ...sources.map((source) => `- ${source.kind}: ${source.reference} \u2014 ${source.purpose}`)
  ].join("\n");
}
function renderManagedWorkspaceReadme(config) {
  const purpose = config.workspace.purpose ?? "Project purpose has not been recorded yet.";
  return `${managedStart}
# ${config.workspace.name}

${purpose}

## Product repositories

${repositoryRows(config)}

## Common actions

- Configure this wrapper: \`$configure-workspace\` (Codex) or \`/configure-workspace\` (Claude Code).
- Choose reviewed work: \`$whats-next\`.
- Run explicitly selected work: \`$run-task\`.
- Create an optional reviewed plan: \`$create-plan\`.
- Curate completed-work learning: \`$sync-context\`.

## Project context

- [Project summary](context/PROJECT.md)
- [Architecture](context/ARCHITECTURE.md)
- [Conventions](context/CONVENTIONS.md)
- [Decisions](context/DECISIONS.md)
${sourceLinks(config)}
- [Approved and draft plans](context/plans/)

Workspace mode: **${config.workspace.mode}**. Review mode: **${config.workflow.review_mode ?? (config.workflow.wrapper_change_policy === "pull-request" ? "remote" : "local")}**.
${managedEnd}`;
}
function reconcileWorkspaceReadme(current, config) {
  const managed = renderManagedWorkspaceReadme(config);
  const starts = [...current.matchAll(new RegExp(managedStart.replace(/[.*+?^${}()|[\]\\]/g, "\\$&"), "g"))];
  const ends = [...current.matchAll(new RegExp(managedEnd.replace(/[.*+?^${}()|[\]\\]/g, "\\$&"), "g"))];
  if (starts.length > 1 || ends.length > 1) throw new Error("README.md must contain at most one managed workspace block");
  const start = starts[0]?.index ?? -1;
  const end = ends[0]?.index ?? -1;
  if (start === -1 !== (end === -1) || start !== -1 && end < start) throw new Error("Malformed managed workspace block in README.md");
  if (start !== -1) {
    const after = end + managedEnd.length;
    return `${current.slice(0, start)}${managed}${current.slice(after)}`.replace(/\s*$/, "\n");
  }
  const existing = current.trim();
  if (!existing) return `${managed}
`;
  if (existing.startsWith(canonicalFrameworkTitle)) {
    return `${managed}

## Context Circuit framework

${existing.slice(canonicalFrameworkTitle.length).trimStart()}
`;
  }
  return `${managed}

${existing}
`;
}
var managedStart, managedEnd, canonicalFrameworkTitle;
var init_workspace_readme = __esm({
  "scripts/lib/workspace-readme.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    managedStart = "<!-- context-circuit:workspace:start -->";
    managedEnd = "<!-- context-circuit:workspace:end -->";
    canonicalFrameworkTitle = "# Context Circuit\n";
  }
});

// scripts/lib/initialize-workspace.ts
import { createHash } from "node:crypto";
import { access as access3, lstat as lstat3, mkdir as mkdir2, readFile as readFile4, realpath as realpath2 } from "node:fs/promises";
import { dirname as dirname5, join as join3, relative as relative2, resolve as resolve5 } from "node:path";
function normalizedRepositoryPath(path2) {
  return `${path2.replace(/^\.\//, "").replace(/\/$/, "")}/`;
}
function parseSubmodulePaths(raw) {
  const paths2 = /* @__PURE__ */ new Set();
  for (const match of raw.matchAll(/^\s*path\s*=\s*(.+?)\s*$/gm)) paths2.add(match[1].replace(/\/$/, ""));
  return paths2;
}
function dangerousRepositoryIgnore(line) {
  const normalized = line.trim().replace(/^\//, "");
  return ["repositories", "repositories/", "repositories/*", "repositories/**"].includes(normalized);
}
function reconcileIgnoredClones(current, config) {
  const lines = current.replace(/\r\n/g, "\n").split("\n");
  const start = lines.indexOf(ignoredStart);
  const end = lines.indexOf(ignoredEnd);
  if (start === -1 !== (end === -1) || start !== -1 && end < start) {
    throw new Error("Malformed managed ignored-clone block in .gitignore");
  }
  const withoutManaged = start === -1 ? [...lines] : [...lines.slice(0, start), ...lines.slice(end + 1)];
  if (withoutManaged.some(dangerousRepositoryIgnore)) {
    throw new Error("Blanket repositories/ ignore conflicts with submodule support; use exact ignored-clone paths");
  }
  const configuredPaths = new Set(Object.values(config.repositories).map((repository) => normalizedRepositoryPath(repository.path)));
  const retained = withoutManaged.filter((line) => !configuredPaths.has(normalizedRepositoryPath(line.trim())));
  while (retained.length > 0 && retained.at(-1) === "") retained.pop();
  const ignored = Object.values(config.repositories).filter((repository) => repository.mode === "ignored-clone").map((repository) => normalizedRepositoryPath(repository.path)).sort();
  const managed = [ignoredStart, ...ignored, ignoredEnd];
  return `${[...retained, ...retained.length > 0 ? [""] : [], ...managed].join("\n")}
`;
}
async function repositoryRemote(path2) {
  try {
    return await git(path2, ["remote", "get-url", "origin"]);
  } catch {
    return null;
  }
}
async function assertDefaultBranch(path2, repository, branch) {
  for (const ref of [`refs/heads/${branch}`, `refs/remotes/origin/${branch}`]) {
    try {
      await git(path2, ["rev-parse", "--verify", ref]);
      return;
    } catch {
    }
  }
  throw new Error(`Repository ${repository} has no local or origin default branch named ${branch}`);
}
async function pathExists(path2) {
  try {
    await lstat3(path2);
    return true;
  } catch (error) {
    if (error.code === "ENOENT") return false;
    throw error;
  }
}
async function assertSafeRepositoryPath(workspaceRoot20, path2, name) {
  if (path2 === workspaceRoot20) throw new Error(`Repository ${name} path cannot be the wrapper root`);
  let ancestor = dirname5(path2);
  while (!await pathExists(ancestor)) {
    const parent = dirname5(ancestor);
    if (parent === ancestor) throw new Error(`Cannot resolve repository parent for ${name}`);
    ancestor = parent;
  }
  const info = await lstat3(ancestor);
  if (info.isSymbolicLink()) throw new Error(`Repository ${name} parent cannot be a symbolic link`);
  assertInside(await realpath2(workspaceRoot20), await realpath2(ancestor));
}
async function hasHead(path2) {
  try {
    await git(path2, ["rev-parse", "--verify", "HEAD"]);
    return true;
  } catch {
    return false;
  }
}
function safeRemote(value2, repository) {
  const remote = value2.trim();
  const error = cloneReferenceError(remote);
  if (error) throw new Error(`Repository ${repository} clone URL ${error}`);
  return remote;
}
async function assertExpectedUnbornTemplate(root) {
  const status2 = (await git(root, ["status", "--porcelain=v1", "--untracked-files=all"])).split("\n").filter(Boolean);
  const allowed = /* @__PURE__ */ new Set([...requiredWorkspaceDocuments, ".gitignore", "template-manifest.json"]);
  const templateDirectories = [".agents/", ".codex/", ".claude/", "agents/", "context/", "contributions/", "docs/"];
  const trustedInventory = typeof define_CC_TEMPLATE_INVENTORY_default === "undefined" ? null : define_CC_TEMPLATE_INVENTORY_default;
  if (trustedInventory) {
    const manifestPath = join3(root, "template-manifest.json");
    const manifest2 = JSON.parse(await readFile4(manifestPath, "utf8"));
    const inventory = manifest2.file_inventory;
    if (manifest2.name !== "context-circuit" || manifest2.version !== "0.2.1" || manifest2.node !== ">=22" || manifest2.command !== "node .agents/bin/cc.mjs" || !Array.isArray(inventory) || JSON.stringify(inventory) !== JSON.stringify(trustedInventory)) {
      throw new Error("Extracted template manifest or inventory has been modified");
    }
    const expectedBundle = createHash("sha256").update(await readFile4(join3(root, ".agents", "bin", "cc.mjs"))).digest("hex");
    if (manifest2.bundle_sha256 !== expectedBundle) throw new Error("Extracted template manifest bundle digest has been modified");
    allowed.clear();
    for (const path2 of trustedInventory) allowed.add(path2);
    templateDirectories.length = 0;
  }
  const unexpected = status2.filter((line) => {
    if (!line.startsWith("?? ")) return true;
    const path2 = line.slice(3);
    return !allowed.has(path2) && !templateDirectories.some((prefix) => path2.startsWith(prefix));
  });
  if (unexpected.length > 0) throw new Error(`Unborn wrapper contains authored or unexpected changes; refusing bootstrap:
${unexpected.join("\n")}`);
  const config = await readData(join3(root, "workspace.yaml"));
  if (config.workspace.name !== "uninitialized-workspace" || Object.keys(config.repositories).length !== 0) {
    throw new Error("Unborn wrapper is not the neutral extracted-template baseline");
  }
}
function commitArgs(commit) {
  if (Boolean(commit.author_name) !== Boolean(commit.author_email)) throw new Error("Commit author_name and author_email must be supplied together");
  const args = [];
  if (commit.author_name && commit.author_email) args.push("-c", `user.name=${commit.author_name}`, "-c", `user.email=${commit.author_email}`);
  return [...args, "commit", "--allow-empty", "-m", commit.commit_message.trim()];
}
function agentDocument(name, role) {
  const title = name.split("-").map((part) => `${part.charAt(0).toUpperCase()}${part.slice(1)}`).join(" ");
  return `# ${title} worker

Follow \`repository-worker.md\`. This repository owns the ${role} role. Read its repository-local instructions and preserve its established architecture, conventions, and verification commands.
`;
}
async function assertExactGitRoot(path2, name) {
  const topLevel = await git(path2, ["rev-parse", "--show-toplevel"]);
  if (await realpath2(topLevel) !== await realpath2(path2)) throw new Error(`Repository path is not a Git root: ${name}`);
}
async function bootstrapWorkspace(options) {
  const workspaceRoot20 = resolve5(options.workspaceRoot);
  const requestErrors = await validateContract("workspace-bootstrap-request", options.request);
  if (requestErrors.length > 0) throw new Error(`Invalid workspace-bootstrap-request: ${requestErrors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
  const config = options.request.configuration;
  const semanticErrors = workspaceSemanticErrors(config);
  if (semanticErrors.length > 0) throw new Error(`Invalid workspace configuration: ${semanticErrors.join("; ")}`);
  if (config.workspace.name === "uninitialized-workspace") throw new Error("Bootstrap requires a human-selected workspace name");
  const configuredNames = Object.keys(config.repositories).sort();
  if (configuredNames.length === 0) throw new Error("Bootstrap requires at least one configured repository");
  const actionsByName = new Map(options.request.repositories.map((action) => [action.name, action]));
  if (actionsByName.size !== options.request.repositories.length || configuredNames.join("\n") !== [...actionsByName.keys()].sort().join("\n")) {
    throw new Error("Bootstrap repository actions must match configured repositories exactly");
  }
  const wrapperGitExists = await pathExists(join3(workspaceRoot20, ".git"));
  if (wrapperGitExists === options.request.wrapper.initialize_git) {
    throw new Error(wrapperGitExists ? "Wrapper is already a Git repository; initialize_git must be false" : "Wrapper is not a Git repository; initialize_git must be true");
  }
  if (wrapperGitExists) {
    await assertExactGitRoot(workspaceRoot20, "wrapper");
    if (!await hasHead(workspaceRoot20)) await assertExpectedUnbornTemplate(workspaceRoot20);
    else {
      const changes = await git(workspaceRoot20, ["status", "--porcelain=v1", "--untracked-files=normal"]);
      if (changes) throw new Error(`Wrapper has existing changes; refusing bootstrap:
${changes}`);
    }
  }
  const wrapperHadHead = wrapperGitExists && await hasHead(workspaceRoot20);
  if (!wrapperHadHead && !options.request.wrapper.authorize_initial_commit) throw new Error("A new or unborn wrapper requires explicit initial-commit authorization");
  if (wrapperHadHead && options.request.wrapper.authorize_initial_commit) throw new Error("An existing wrapper must not authorize another initial commit");
  const gitignorePath = join3(workspaceRoot20, ".gitignore");
  const currentGitignore = await readFile4(gitignorePath, "utf8");
  const nextGitignore = reconcileIgnoredClones(currentGitignore, config);
  const readmePath = join3(workspaceRoot20, "README.md");
  const readmeInfo = await lstat3(readmePath);
  if (!readmeInfo.isFile() || readmeInfo.isSymbolicLink()) throw new Error("README.md must be a regular non-symlink file");
  assertInside(await realpath2(workspaceRoot20), await realpath2(readmePath));
  const currentReadme = await readFile4(readmePath, "utf8");
  const sourcesPath = join3(workspaceRoot20, "context", "SOURCES.md");
  const sourcesInfo = await lstat3(sourcesPath);
  if (!sourcesInfo.isFile() || sourcesInfo.isSymbolicLink()) throw new Error("context/SOURCES.md must be a regular non-symlink file");
  assertInside(await realpath2(workspaceRoot20), await realpath2(sourcesPath));
  const readmeConfig = config.workspace.purpose ? config : { ...config, workspace: { ...config.workspace, purpose: options.request.context.project_summary } };
  const nextReadme = reconcileWorkspaceReadme(currentReadme, readmeConfig);
  for (const agent of new Set(Object.values(config.repositories).map((repository) => repository.agent))) {
    const agentPath = join3(workspaceRoot20, "agents", `${agent}.md`);
    if (!await pathExists(agentPath)) continue;
    const info = await lstat3(agentPath);
    if (!info.isFile() || info.isSymbolicLink()) throw new Error(`Domain agent path must be a regular file: agents/${agent}.md`);
  }
  for (const name of configuredNames) {
    const repository = config.repositories[name];
    const action = actionsByName.get(name);
    const path2 = assertInside(workspaceRoot20, resolve5(workspaceRoot20, repository.path));
    await assertSafeRepositoryPath(workspaceRoot20, path2, name);
    const exists2 = await pathExists(path2);
    if (action.source === "existing") {
      if (repository.mode !== "ignored-clone") throw new Error(`Existing repository ${name} must use ignored-clone mode`);
      if (!exists2) throw new Error(`Existing repository path is not accessible: ${repository.path}`);
      if ((await lstat3(path2)).isSymbolicLink()) throw new Error(`Existing repository ${name} cannot be a symbolic link`);
      assertInside(await realpath2(workspaceRoot20), await realpath2(path2));
      await assertExactGitRoot(path2, name);
    } else {
      if (exists2) throw new Error(`Bootstrap refuses to replace existing path for ${name}: ${repository.path}`);
    }
    if (action.source === "submodule" && repository.mode !== "submodule") throw new Error(`Submodule action requires submodule mode for ${name}`);
    if ((action.source === "new" || action.source === "clone") && repository.mode !== "ignored-clone") throw new Error(`${action.source} action requires ignored-clone mode for ${name}`);
    if ((action.source === "clone" || action.source === "submodule") && !action.url) throw new Error(`${action.source} action requires a URL for ${name}`);
    if ((action.source === "new" || action.source === "existing") && action.url) throw new Error(`${action.source} repository ${name} must not include a clone URL`);
    if (action.url) safeRemote(action.url, name);
    if (action.source === "new" && !action.authorize_initial_commit) throw new Error(`New repository ${name} requires explicit initial-commit authorization`);
    if (action.source === "new" && !action.commit_message) throw new Error(`New repository ${name} requires an initial commit message`);
    if (action.source !== "new" && action.authorize_initial_commit) throw new Error(`${action.source} repository ${name} must not authorize an initial commit`);
    if (action.source !== "new" && (action.commit_message || action.author_name || action.author_email)) throw new Error(`${action.source} repository ${name} must not include artificial commit metadata`);
  }
  const bootstrapActions = [];
  if (!wrapperGitExists) {
    await git(workspaceRoot20, ["init", "--initial-branch", config.workspace.default_branch]);
    bootstrapActions.push(`initialized wrapper Git repository on ${config.workspace.default_branch}`);
  } else if (!wrapperHadHead) {
    const current = await git(workspaceRoot20, ["symbolic-ref", "--short", "HEAD"]);
    if (current !== config.workspace.default_branch) throw new Error(`Unborn wrapper branch is ${current}, expected ${config.workspace.default_branch}`);
  }
  await writeTextAtomic(join3(workspaceRoot20, "workspace.yaml"), (0, import_yaml3.stringify)(config));
  for (const [path2, contents] of Object.entries(renderWorkspaceContext(options.request.context))) {
    await writeTextAtomic(join3(workspaceRoot20, path2), contents);
  }
  if (options.request.context.product_knowledge) {
    for (const [path2, contents] of Object.entries(renderProductKnowledgeBaseline(options.request.context.product_knowledge))) {
      const full = assertInside(workspaceRoot20, resolve5(workspaceRoot20, path2));
      await mkdir2(dirname5(full), { recursive: true });
      await writeTextAtomic(full, contents);
    }
  }
  await writeTextAtomic(readmePath, nextReadme);
  await mkdir2(join3(workspaceRoot20, "agents"), { recursive: true });
  for (const [name, repository] of Object.entries(config.repositories)) {
    const agentPath = join3(workspaceRoot20, "agents", `${repository.agent}.md`);
    if (!await pathExists(agentPath)) await writeTextAtomic(agentPath, agentDocument(name, repository.role));
  }
  await writeTextAtomic(gitignorePath, nextGitignore);
  for (const name of configuredNames) {
    const repository = config.repositories[name];
    const action = actionsByName.get(name);
    const path2 = assertInside(workspaceRoot20, resolve5(workspaceRoot20, repository.path));
    if (action.source === "new") {
      await mkdir2(dirname5(path2), { recursive: true });
      await mkdir2(path2);
      await git(path2, ["init", "--initial-branch", repository.default_branch]);
      await git(path2, commitArgs(action));
      bootstrapActions.push(`created ${name} with an empty base commit`);
    } else if (action.source === "clone") {
      await mkdir2(dirname5(path2), { recursive: true });
      await git(workspaceRoot20, ["clone", "--branch", repository.default_branch, "--single-branch", "--", safeRemote(action.url, name), path2]);
      bootstrapActions.push(`cloned ${name} into ${repository.path}`);
    } else if (action.source === "submodule") {
      await mkdir2(dirname5(path2), { recursive: true });
      await git(workspaceRoot20, ["-c", "protocol.file.allow=always", "submodule", "add", "-b", repository.default_branch, "--", safeRemote(action.url, name), repository.path]);
      bootstrapActions.push(`registered ${name} as a submodule`);
    } else {
      bootstrapActions.push(`registered existing repository ${name}`);
    }
  }
  await initializeWorkspace({ workspaceRoot: workspaceRoot20, allowUnbornWrapper: !wrapperHadHead });
  let wrapperInitialCommit = null;
  if (!wrapperHadHead) {
    await git(workspaceRoot20, ["add", "-A"]);
    await git(workspaceRoot20, commitArgs(options.request.wrapper));
    wrapperInitialCommit = await git(workspaceRoot20, ["rev-parse", "HEAD"]);
    bootstrapActions.push("created configured wrapper initial commit");
  }
  const summary2 = await initializeWorkspace({ workspaceRoot: workspaceRoot20 });
  return { ...summary2, status: "initialized", bootstrap_actions: bootstrapActions, wrapper_initial_commit: wrapperInitialCommit };
}
async function initializeWorkspace(options) {
  const workspaceRoot20 = resolve5(options.workspaceRoot);
  const configPath = join3(workspaceRoot20, "workspace.yaml");
  const config = await readData(configPath);
  const contractErrors2 = await validateContract("workspace", config);
  const productKnowledge = await validateProductKnowledgeTree(join3(workspaceRoot20, "context"));
  const errors2 = [
    ...contractErrors2.map((error) => `${error.instancePath || "/"} ${error.message}`),
    ...workspaceSemanticErrors(config),
    ...await workspaceDocumentErrors(workspaceRoot20, config),
    ...productKnowledge.errors.map((error) => `product-knowledge ${error}`)
  ];
  if (errors2.length > 0) throw new Error(`Workspace initialization validation failed:
- ${errors2.join("\n- ")}`);
  await git(workspaceRoot20, ["rev-parse", "--is-inside-work-tree"]);
  const wrapperTopLevel = await git(workspaceRoot20, ["rev-parse", "--show-toplevel"]);
  if (await realpath2(wrapperTopLevel) !== await realpath2(workspaceRoot20)) {
    throw new Error(`Workspace root is not the wrapper Git root: ${workspaceRoot20}`);
  }
  if (options.allowUnbornWrapper && !await hasHead(workspaceRoot20)) {
    const current = await git(workspaceRoot20, ["symbolic-ref", "--short", "HEAD"]);
    if (current !== config.workspace.default_branch) throw new Error(`Wrapper branch is ${current}, expected ${config.workspace.default_branch}`);
  } else {
    await assertDefaultBranch(workspaceRoot20, "wrapper", config.workspace.default_branch);
  }
  let submodulePaths = /* @__PURE__ */ new Set();
  try {
    submodulePaths = parseSubmodulePaths(await readFile4(join3(workspaceRoot20, ".gitmodules"), "utf8"));
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
  }
  const repositories = [];
  const warnings = [];
  for (const [name, repository] of Object.entries(config.repositories)) {
    const path2 = assertInside(workspaceRoot20, resolve5(workspaceRoot20, repository.path));
    try {
      await access3(path2);
    } catch {
      throw new Error(`Repository ${name} path is not accessible: ${repository.path}`);
    }
    assertInside(await realpath2(workspaceRoot20), await realpath2(path2));
    const topLevel = await git(path2, ["rev-parse", "--show-toplevel"]);
    if (await realpath2(topLevel) !== await realpath2(path2)) throw new Error(`Repository path is not a Git root: ${repository.path}`);
    const relativePath = relative2(workspaceRoot20, path2).replaceAll("\\", "/");
    const trackedEntry = await git(workspaceRoot20, ["ls-files", "--stage", "--", relativePath]);
    if (repository.mode === "submodule" && !submodulePaths.has(relativePath)) {
      throw new Error(`Repository ${name} is configured as a submodule but is not registered in .gitmodules: ${relativePath}`);
    }
    if (repository.mode === "submodule" && !trackedEntry.startsWith("160000 ")) {
      throw new Error(`Repository ${name} is configured as a submodule but the wrapper does not track a gitlink: ${relativePath}`);
    }
    if (repository.mode === "ignored-clone" && submodulePaths.has(relativePath)) {
      throw new Error(`Repository ${name} is registered as a submodule but configured as an ignored clone`);
    }
    if (repository.mode === "ignored-clone" && trackedEntry) {
      throw new Error(`Repository ${name} is tracked by the wrapper but configured as an ignored clone`);
    }
    await assertDefaultBranch(path2, name, repository.default_branch);
    const instructionsPath = join3(path2, "AGENTS.md");
    let instructions = null;
    try {
      await access3(instructionsPath);
      instructions = instructionsPath;
    } catch {
      warnings.push(`Repository ${name} has no repository-local AGENTS.md`);
    }
    const remote = await repositoryRemote(path2);
    if (!remote) warnings.push(`Repository ${name} has no origin remote`);
    const currentBranch = await git(path2, ["rev-parse", "--abbrev-ref", "HEAD"]);
    const repositoryChanges = (await git(path2, ["status", "--porcelain=v1", "--untracked-files=normal"])).split("\n").filter(Boolean);
    repositories.push({
      name,
      path: relativePath,
      mode: repository.mode,
      role: repository.role,
      agent: repository.agent,
      default_branch: repository.default_branch,
      current_branch: currentBranch,
      clean: repositoryChanges.length === 0,
      remote,
      instructions
    });
  }
  if (config.activity.provider === "none") warnings.push("No activity provider is configured; planless work remains available");
  const gitignorePath = join3(workspaceRoot20, ".gitignore");
  let currentGitignore = "";
  try {
    currentGitignore = await readFile4(gitignorePath, "utf8");
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
  }
  const nextGitignore = reconcileIgnoredClones(currentGitignore, config);
  const gitignoreChanged = nextGitignore !== currentGitignore;
  if (options.apply !== false && gitignoreChanged) await writeTextAtomic(gitignorePath, nextGitignore);
  const wrapperChanges = (await git(workspaceRoot20, ["status", "--porcelain=v1", "--untracked-files=all"])).split("\n").filter(Boolean);
  return {
    workspace: config.workspace.name,
    mode: config.workspace.mode,
    default_branch: config.workspace.default_branch,
    activity_provider: config.activity.provider,
    wrapper_change_policy: config.workflow.wrapper_change_policy,
    repositories,
    required_documents: [.../* @__PURE__ */ new Set([
      ...requiredWorkspaceDocuments,
      ...new Set(Object.values(config.repositories).map((repository) => `agents/${repository.agent}.md`))
    ])],
    wrapper_changes: wrapperChanges,
    gitignore_changed: gitignoreChanged,
    applied: options.apply !== false,
    warnings
  };
}
var import_yaml3, ignoredStart, ignoredEnd;
var init_initialize_workspace = __esm({
  "scripts/lib/initialize-workspace.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    import_yaml3 = __toESM(require_dist(), 1);
    init_io();
    init_git();
    init_validation();
    init_workspace_context();
    init_product_knowledge();
    init_workspace_readme();
    init_safe_reference();
    ignoredStart = "# context-circuit:ignored-clones:start";
    ignoredEnd = "# context-circuit:ignored-clones:end";
  }
});

// scripts/lib/configure-workspace.ts
import { lstat as lstat4, readFile as readFile5, readdir as readdir2, realpath as realpath3 } from "node:fs/promises";
import { join as join4, relative as relative3, resolve as resolve6 } from "node:path";
async function exists(path2) {
  try {
    await lstat4(path2);
    return true;
  } catch (error) {
    if (error.code === "ENOENT") return false;
    throw error;
  }
}
async function hasHead2(root) {
  try {
    await git(root, ["rev-parse", "--verify", "HEAD"]);
    return true;
  } catch {
    return false;
  }
}
async function interruptedConfigurationArtifacts(workspaceRoot20) {
  const root = resolve6(workspaceRoot20);
  const directories = [root, join4(root, "context"), join4(root, "agents")];
  const artifacts = [];
  for (const directory of directories) {
    let entries;
    try {
      entries = await readdir2(directory, { withFileTypes: true });
    } catch (error) {
      if (error.code === "ENOENT") continue;
      throw error;
    }
    for (const entry of entries) {
      if (!entry.isFile() || !transactionResidue.test(entry.name)) continue;
      const original = entry.name.replace(transactionResidue, "");
      const managed = directory === root ? ["workspace.yaml", "README.md", ".gitignore"].includes(original) : directory === join4(root, "context") ? original === "SOURCES.md" : /^[a-z][a-z0-9-]*\.md$/.test(original);
      if (managed) artifacts.push(relative3(root, join4(directory, entry.name)).replaceAll("\\", "/"));
    }
  }
  return artifacts.sort();
}
async function assertNoInterruptedConfiguration(workspaceRoot20) {
  const artifacts = await interruptedConfigurationArtifacts(workspaceRoot20);
  if (artifacts.length === 0) return;
  throw new Error(
    `Interrupted workspace configuration artifacts were found:
- ${artifacts.join("\n- ")}
Configuration will not delete or overwrite them. Inspect each target, .stage, and .backup sibling; restore exactly one authoritative target manually; preserve uncertain copies; then rerun configure-workspace.`
  );
}
async function detectWorkspaceConfigurationState(workspaceRoot20) {
  const root = resolve6(workspaceRoot20);
  if (!await exists(join4(root, ".git")) || !await hasHead2(root)) return "fresh";
  return "existing";
}
function workspaceCredentialErrors(request4) {
  const remotes = [
    ["configuration.workspace.remote", request4.configuration.workspace.remote],
    ...Object.entries(request4.configuration.repositories).map(([name, repository]) => [`configuration.repositories.${name}.remote`, repository.remote])
  ];
  const errors2 = [];
  for (const [path2, value2] of remotes) {
    const error = value2 ? remoteReferenceError(value2) : null;
    if (error) errors2.push(`${path2} ${error}`);
  }
  for (const [index, repository] of request4.repositories.entries()) {
    const error = repository.url ? cloneReferenceError(repository.url) : null;
    if (error) errors2.push(`repositories.${index}.url ${error}`);
  }
  for (const [index, source] of (request4.context.sources ?? []).entries()) {
    const error = contextReferenceError(source.reference);
    if (error) errors2.push(`context.sources.${index}.reference ${error}`);
  }
  return errors2;
}
function assertSourceConsistency(request4) {
  const configured = request4.configuration.context?.authoritative_sources ?? [];
  const requested2 = request4.context.sources ?? [];
  if (JSON.stringify(configured) !== JSON.stringify(requested2)) {
    throw new Error("configuration.context.authoritative_sources must exactly match context.sources");
  }
  const names = new Set(Object.keys(request4.configuration.repositories));
  for (const [index, source] of requested2.entries()) {
    if (source.repository && !names.has(source.repository)) throw new Error(`context.sources.${index}.repository is not configured: ${source.repository}`);
  }
  for (const [name, repository] of Object.entries(request4.configuration.repositories)) {
    const action = request4.repositories.find((candidate) => candidate.name === name);
    if (repository.remote && action?.url && repository.remote !== action.url) throw new Error(`Configured remote and source URL differ for repository ${name}`);
  }
}
async function validateRequest(request4) {
  const contractErrors2 = await validateContract("workspace-bootstrap-request", request4);
  const errors2 = [
    ...contractErrors2.map((error) => `${error.instancePath || "/"} ${error.message}`),
    ...workspaceSemanticErrors(request4.configuration),
    ...workspaceCredentialErrors(request4)
  ];
  if (errors2.length > 0) throw new Error(`Invalid workspace configuration request:
- ${errors2.join("\n- ")}`);
  assertSourceConsistency(request4);
}
async function readRegularInside(root, path2, label) {
  const candidate = assertInside(root, path2);
  const info = await lstat4(candidate);
  if (!info.isFile() || info.isSymbolicLink()) throw new Error(`${label} must be a regular non-symlink file`);
  assertInside(await realpath3(root), await realpath3(candidate));
  return readFile5(candidate, "utf8");
}
async function assertExactGitRoot2(path2, label) {
  let info;
  try {
    info = await lstat4(path2);
  } catch (error) {
    if (error.code === "ENOENT") throw new Error(`${label} path does not exist`);
    throw error;
  }
  if (!info.isDirectory() || info.isSymbolicLink()) throw new Error(`${label} must be a real directory`);
  const top = await git(path2, ["rev-parse", "--show-toplevel"]);
  if (await realpath3(top) !== await realpath3(path2)) throw new Error(`${label} is not an exact Git root`);
}
async function assertBranch(path2, name, branch) {
  for (const ref of [`refs/heads/${branch}`, `refs/remotes/origin/${branch}`]) {
    try {
      await git(path2, ["rev-parse", "--verify", ref]);
      return;
    } catch {
    }
  }
  throw new Error(`Repository ${name} has no local or origin default branch named ${branch}`);
}
async function preflightExistingRepositories(root, request4) {
  let submodules = "";
  try {
    submodules = await readFile5(join4(root, ".gitmodules"), "utf8");
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
  }
  for (const [name, repository] of Object.entries(request4.configuration.repositories)) {
    const path2 = assertInside(root, resolve6(root, repository.path));
    await assertExactGitRoot2(path2, `Repository ${name}`);
    await assertBranch(path2, name, repository.default_branch);
    const relativePath = relative3(root, path2).replaceAll("\\", "/");
    const tracked = await git(root, ["ls-files", "--stage", "--", relativePath]);
    const registered = new RegExp(`^\\s*path\\s*=\\s*${relativePath.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}\\s*$`, "m").test(submodules);
    if (repository.mode === "submodule" && (!registered || !tracked.startsWith("160000 "))) throw new Error(`Repository ${name} is not a tracked submodule: ${relativePath}`);
    if (repository.mode === "ignored-clone" && (registered || tracked)) throw new Error(`Repository ${name} is tracked but configured as an ignored clone`);
  }
}
async function exactBootstrapRerun(root, request4) {
  if (request4.authorize_reviewable_changes || !request4.wrapper.authorize_initial_commit) return false;
  let installed;
  try {
    installed = (0, import_yaml4.parse)(await readFile5(join4(root, "workspace.yaml"), "utf8"));
  } catch {
    return false;
  }
  if (JSON.stringify(installed) !== JSON.stringify(request4.configuration)) return false;
  const readme = await readRegularInside(root, join4(root, "README.md"), "README.md");
  const sources = await readRegularInside(root, join4(root, "context", "SOURCES.md"), "context/SOURCES.md");
  const readmeConfig = request4.configuration.workspace.purpose ? request4.configuration : { ...request4.configuration, workspace: { ...request4.configuration.workspace, purpose: request4.context.project_summary } };
  if (reconcileWorkspaceReadme(readme, readmeConfig) !== readme) return false;
  if (renderWorkspaceContext(request4.context)["context/SOURCES.md"] !== sources) return false;
  try {
    await preflightExistingRepositories(root, {
      ...request4,
      repositories: request4.repositories.map((repository) => ({ name: repository.name, source: "existing", authorize_initial_commit: false }))
    });
  } catch {
    return false;
  }
  return true;
}
async function reconfigureWorkspace(workspaceRoot20, request4, transactionOptions = {}) {
  const root = resolve6(workspaceRoot20);
  const changes = await git(root, ["status", "--porcelain=v1", "--untracked-files=all"]);
  if (changes) throw new Error(`Wrapper must be clean before reconfiguration; refusing to overwrite existing work:
${changes}`);
  if (request4.authorize_reviewable_changes !== true) throw new Error("Existing wrapper reconfiguration requires explicit authorize_reviewable_changes: true");
  if (request4.wrapper.initialize_git || request4.wrapper.authorize_initial_commit) {
    throw new Error("An existing wrapper must not initialize Git or authorize an initial commit");
  }
  const configuredNames = Object.keys(request4.configuration.repositories).sort();
  const actionNames = request4.repositories.map((repository) => repository.name).sort();
  if (new Set(actionNames).size !== actionNames.length || configuredNames.join("\n") !== actionNames.join("\n")) {
    throw new Error("Reconfiguration repository actions must match configured repositories exactly");
  }
  for (const repository of request4.repositories) {
    if (repository.source !== "existing") throw new Error(`Existing wrapper reconfiguration accepts only inspected existing repository paths: ${repository.name}`);
    if (repository.authorize_initial_commit || repository.commit_message || repository.author_name || repository.author_email) {
      throw new Error(`Existing repository ${repository.name} must not include initial-commit authorization or metadata`);
    }
  }
  const top = await git(root, ["rev-parse", "--show-toplevel"]);
  if (await realpath3(top) !== await realpath3(root)) throw new Error("Workspace root is not the wrapper Git root");
  const current = await readRegularInside(root, join4(root, "README.md"), "README.md");
  await readRegularInside(root, join4(root, "context", "SOURCES.md"), "context/SOURCES.md");
  const sources = renderWorkspaceContext(request4.context)["context/SOURCES.md"];
  const readme = reconcileWorkspaceReadme(current, request4.configuration);
  const gitignorePath = join4(root, ".gitignore");
  const gitignore = reconcileIgnoredClones(await readRegularInside(root, gitignorePath, ".gitignore"), request4.configuration);
  await preflightExistingRepositories(root, request4);
  const agentWrites = [];
  const agentsDirectory = join4(root, "agents");
  const agentsInfo = await lstat4(agentsDirectory);
  if (!agentsInfo.isDirectory() || agentsInfo.isSymbolicLink()) throw new Error("agents must be a real directory");
  for (const [name, repository] of Object.entries(request4.configuration.repositories)) {
    const path2 = join4(root, "agents", `${repository.agent}.md`);
    if (await exists(path2)) await readRegularInside(root, path2, `agents/${repository.agent}.md`);
    else agentWrites.push([path2, `# ${repository.agent}

Follow \`repository-worker.md\`. This repository owns the ${repository.role} role. Read ${name}'s repository-local instructions before work.
`]);
  }
  const workspace = (0, import_yaml4.stringify)(request4.configuration);
  await writeTextTransaction([
    { path: join4(root, "workspace.yaml"), value: workspace },
    { path: join4(root, "context/SOURCES.md"), value: sources },
    ...agentWrites.map(([path2, value2]) => ({ path: path2, value: value2 })),
    { path: join4(root, "README.md"), value: readme },
    { path: gitignorePath, value: gitignore }
  ], transactionOptions);
  return initializeWorkspace({ workspaceRoot: root });
}
async function configureWorkspace(options) {
  const workspaceRoot20 = resolve6(options.workspaceRoot);
  await assertNoInterruptedConfiguration(workspaceRoot20);
  const state = await detectWorkspaceConfigurationState(workspaceRoot20);
  if (!options.request) {
    if (state === "fresh") return { route: "inspect-fresh", state, message: "Fresh wrapper detected; collect a configuration request, then run the internal bootstrap phase with exact initial-commit authorization.", result: null };
    const result4 = await initializeWorkspace({ workspaceRoot: workspaceRoot20, apply: false });
    return { route: "inspect-existing", state, message: "Existing wrapper detected; configuration changes will remain reviewable and uncommitted.", result: result4 };
  }
  await validateRequest(options.request);
  if (options.checkOnly) return { route: state === "fresh" ? "inspect-fresh" : "inspect-existing", state, message: "Configuration request is valid; no files or Git state changed.", result: state === "existing" ? await initializeWorkspace({ workspaceRoot: workspaceRoot20, apply: false }) : null };
  if (state === "fresh") {
    if (options.request.authorize_reviewable_changes) throw new Error("Fresh bootstrap must not authorize existing-wrapper reconfiguration");
    const result4 = await bootstrapWorkspace({ workspaceRoot: workspaceRoot20, request: options.request });
    return { route: "bootstrap", state, message: "Fresh wrapper configured through the explicit bootstrap phase.", result: result4 };
  }
  if (await exactBootstrapRerun(workspaceRoot20, options.request)) {
    const result4 = await initializeWorkspace({ workspaceRoot: workspaceRoot20, apply: false });
    return { route: "inspect-existing", state, message: "Exact completed bootstrap request detected; configuration is already current and no files or commits changed.", result: result4 };
  }
  const result3 = await reconfigureWorkspace(workspaceRoot20, options.request, options.transactionOptions);
  return { route: "reconfigure", state, message: "Existing wrapper configuration was updated as reviewable, uncommitted changes.", result: result3 };
}
var import_yaml4, transactionResidue;
var init_configure_workspace = __esm({
  "scripts/lib/configure-workspace.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    import_yaml4 = __toESM(require_dist(), 1);
    init_git();
    init_io();
    init_initialize_workspace();
    init_validation();
    init_workspace_context();
    init_workspace_readme();
    init_safe_reference();
    transactionResidue = /\.\d+\.\d+\.[0-9a-f]+\.(?:stage|backup)$/;
  }
});

// scripts/initialize-workspace.ts
var initialize_workspace_exports = {};
import { dirname as dirname7, resolve as resolve7 } from "node:path";
import { parseArgs as parseArgs2 } from "node:util";
import { fileURLToPath as fileURLToPath3 } from "node:url";
var values2, workspaceRoot2, summary;
var init_initialize_workspace2 = __esm({
  async "scripts/initialize-workspace.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_configure_workspace();
    init_io();
    ({ values: values2 } = parseArgs2({
      options: {
        "check-only": { type: "boolean", default: false },
        bootstrap: { type: "string" }
      }
    }));
    workspaceRoot2 = resolve7(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve7(dirname7(fileURLToPath3(import.meta.url)), ".."));
    if (values2.bootstrap && values2["check-only"]) throw new Error("--bootstrap and --check-only cannot be combined");
    summary = values2.bootstrap ? await configureWorkspace({ workspaceRoot: workspaceRoot2, request: await readJsonRegularInside(workspaceRoot2, resolve7(workspaceRoot2, values2.bootstrap), "Workspace bootstrap request") }) : await configureWorkspace({ workspaceRoot: workspaceRoot2, checkOnly: values2["check-only"] });
    console.log(JSON.stringify(summary, null, 2));
  }
});

// scripts/configure-workspace.ts
var configure_workspace_exports = {};
import { dirname as dirname8, resolve as resolve8 } from "node:path";
import { parseArgs as parseArgs3 } from "node:util";
import { fileURLToPath as fileURLToPath4 } from "node:url";
var values3, workspaceRoot3, request;
var init_configure_workspace2 = __esm({
  async "scripts/configure-workspace.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_configure_workspace();
    init_io();
    ({ values: values3 } = parseArgs3({ options: { request: { type: "string" }, "check-only": { type: "boolean", default: false } } }));
    workspaceRoot3 = resolve8(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve8(dirname8(fileURLToPath4(import.meta.url)), ".."));
    request = values3.request ? await readJsonRegularInside(workspaceRoot3, resolve8(workspaceRoot3, values3.request), "Workspace configuration request") : void 0;
    console.log(JSON.stringify(await configureWorkspace({ workspaceRoot: workspaceRoot3, ...request ? { request } : {}, checkOnly: values3["check-only"] }), null, 2));
  }
});

// scripts/lib/ids.ts
import { createHash as createHash2, randomBytes } from "node:crypto";
import { readFile as readFile6 } from "node:fs/promises";
import { join as join5 } from "node:path";
function utcStamp(now) {
  const iso = now.toISOString();
  return {
    day: iso.slice(0, 10).replaceAll("-", ""),
    instant: iso.slice(0, 19).replaceAll("-", "").replaceAll(":", "")
  };
}
function generateRunId(request4, now, discriminator) {
  if (!/^[a-f0-9]{8}$/.test(discriminator)) {
    throw new Error("Run discriminator must contain exactly eight lowercase hexadecimal characters");
  }
  const { instant } = utcStamp(now);
  const requestFingerprint = createHash2("sha256").update(request4).digest("hex").slice(0, 4);
  return `${instant}Z-${discriminator.slice(0, 4)}${requestFingerprint}`;
}
async function generateIds(runtimeRoot, request4, now = /* @__PURE__ */ new Date(), discriminator = randomBytes(4).toString("hex")) {
  const runId = generateRunId(request4, now, discriminator);
  await ensurePrivateDirectory(runtimeRoot);
  const { day } = utcStamp(now);
  const statePath = join5(runtimeRoot, "id-state.json");
  let state = { day, next: 1 };
  try {
    state = JSON.parse(await readFile6(statePath, "utf8"));
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
  }
  const sequence = state.day === day ? state.next : 1;
  await writeJsonAtomic(statePath, { day, next: sequence + 1 });
  return {
    workId: `ADHOC-${day}-${String(sequence).padStart(3, "0")}`,
    runId
  };
}
function slugify2(value2) {
  return value2.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-|-$/g, "").slice(0, 32) || "task";
}
var init_ids = __esm({
  "scripts/lib/ids.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_io();
  }
});

// scripts/lib/activity-lifecycle.ts
import { readFile as readFile7 } from "node:fs/promises";
import { join as join6, relative as relative4, resolve as resolve9 } from "node:path";
async function readJson(path2) {
  return JSON.parse(await readFile7(path2, "utf8"));
}
async function assertValid(name, value2) {
  const errors2 = await validateContract(name, value2);
  if (errors2.length > 0) throw new Error(`Invalid ${name}: ${errors2.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
}
async function loadConfig(workspaceRoot20) {
  const config = (0, import_yaml5.parse)(await readFile7(join6(workspaceRoot20, "workspace.yaml"), "utf8"));
  await assertValid("workspace", config);
  const errors2 = workspaceSemanticErrors(config);
  if (errors2.length > 0) throw new Error(`Invalid workspace: ${errors2.join("; ")}`);
  return config;
}
function safeEvidence(value2, field) {
  const clean = value2.trim();
  if (!clean || /[\r\n]/.test(clean)) throw new Error(`${field} must be a non-empty single line`);
  if (/-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----|https?:\/\/[^\s/@:]+:[^\s/@]+@/i.test(clean)) {
    throw new Error(`${field} appears to contain a credential or private key`);
  }
  return clean;
}
function overallStatus(actions) {
  if (actions.some((action) => action.policy === "required" && action.status === "failed")) return "failed";
  if (actions.some((action) => action.status === "pending")) return "pending";
  if (actions.some((action) => action.status === "manual")) return "manual";
  return actions.length === 0 ? "skipped" : "completed";
}
function syncManifestEvent(manifest2, record, recordPath2, workspaceRoot20) {
  const existing = manifest2.lifecycle_events.find((item) => item.event === record.event);
  const occurredAt = record.updated_at;
  const value2 = {
    event: record.event,
    status: record.status,
    idempotency_key: existing?.idempotency_key ?? `${manifest2.run_id}:lifecycle:${record.event}:${record.provider}`,
    occurred_at: occurredAt,
    record: relative4(workspaceRoot20, recordPath2).replaceAll("\\", "/"),
    actions: record.actions
  };
  if (existing) Object.assign(existing, value2);
  else manifest2.lifecycle_events.push(value2);
  manifest2.updated_at = occurredAt;
}
function paths(workspaceRoot20, runId, event) {
  const runtimeRoot = assertInside(workspaceRoot20, join6(workspaceRoot20, ".runtime"));
  const runRoot = assertInside(runtimeRoot, join6(runtimeRoot, "runs", runId));
  return {
    runtimeRoot,
    manifestPath: join6(runRoot, "manifest.json"),
    recordPath: join6(runRoot, "activity", `${event}.json`)
  };
}
async function prepareActivityLifecycle(options) {
  const workspaceRoot20 = resolve9(options.workspaceRoot);
  const config = await loadConfig(workspaceRoot20);
  const { manifestPath, recordPath: recordPath2 } = paths(workspaceRoot20, options.runId, options.event);
  await ensurePrivateDirectory(join6(workspaceRoot20, ".runtime"));
  const lockPath = `${recordPath2}.lock`;
  return withExclusiveFile(lockPath, async () => {
    try {
      const existing = await readJson(recordPath2);
      await assertValid("activity-lifecycle-record", existing);
      return existing;
    } catch (error) {
      if (error.code !== "ENOENT") throw error;
    }
    const manifest2 = await readJson(manifestPath);
    await assertValid("runtime-manifest", manifest2);
    if (manifest2.run_id !== options.runId) throw new Error("Manifest run ID does not match lifecycle request");
    const available = new Set(options.availableCapabilities ?? []);
    const configured = config.activity.lifecycle?.[options.event] ?? [];
    const actions = [];
    const warnings = [];
    const manualFallbacks = [];
    let stopped = false;
    for (const action of configured) {
      let status2;
      let evidence = null;
      if (stopped) {
        status2 = "skipped";
        evidence = "Not attempted after a required lifecycle action failed.";
      } else if (action.policy === "manual") {
        status2 = "manual";
        manualFallbacks.push(action.description);
      } else if (available.has(action.capability)) {
        status2 = "pending";
      } else if (action.policy === "required") {
        status2 = "manual";
        evidence = `Required capability is unavailable: ${action.capability}; manual completion is required.`;
        manualFallbacks.push(action.description);
        stopped = true;
      } else {
        status2 = "skipped";
        evidence = `Optional capability is unavailable: ${action.capability}.`;
        warnings.push(`${action.id}: ${evidence}`);
        manualFallbacks.push(action.description);
      }
      actions.push({
        ...action,
        status: status2,
        idempotency_key: `${manifest2.run_id}:${options.event}:${action.id}`,
        evidence,
        external_reference: null
      });
    }
    if (config.activity.provider === "none") warnings.push("No activity provider is configured; the semantic event is recorded without an external write.");
    else if (configured.length === 0) warnings.push(`No lifecycle actions are configured for ${options.event}.`);
    const now = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    const record = {
      contract_version: 1,
      work_id: manifest2.work_id,
      run_id: manifest2.run_id,
      provider: config.activity.provider,
      event: options.event,
      status: overallStatus(actions),
      actions,
      warnings,
      manual_fallbacks: manualFallbacks,
      prepared_at: now,
      updated_at: now
    };
    await assertValid("activity-lifecycle-record", record);
    syncManifestEvent(manifest2, record, recordPath2, workspaceRoot20);
    await assertValid("runtime-manifest", manifest2);
    await writeJsonAtomic(recordPath2, record);
    await writeJsonAtomic(manifestPath, manifest2);
    return record;
  });
}
async function recordActivityLifecycleAction(options) {
  const workspaceRoot20 = resolve9(options.workspaceRoot);
  const { manifestPath, recordPath: recordPath2 } = paths(workspaceRoot20, options.runId, options.event);
  return withExclusiveFile(`${recordPath2}.lock`, async () => {
    const record = await readJson(recordPath2);
    await assertValid("activity-lifecycle-record", record);
    const action = record.actions.find((item) => item.id === options.actionId);
    if (!action) throw new Error(`Lifecycle event has no action named ${options.actionId}`);
    const evidence = safeEvidence(options.evidence, "Evidence");
    const externalReference = options.externalReference ? safeEvidence(options.externalReference, "External reference") : null;
    if (action.status === "completed" || action.status === "failed") {
      if (action.status === options.status && action.evidence === evidence && action.external_reference === externalReference) return record;
      throw new Error(`Lifecycle action ${action.id} already has a different terminal result`);
    }
    if (action.status === "skipped") throw new Error(`Lifecycle action ${action.id} was skipped and cannot receive an external result`);
    action.status = options.status;
    action.evidence = evidence;
    action.external_reference = externalReference;
    if (options.status === "failed" && action.policy === "required") {
      const index = record.actions.indexOf(action);
      for (const remaining of record.actions.slice(index + 1)) {
        if (remaining.status === "pending" || remaining.status === "manual") {
          remaining.status = "skipped";
          remaining.evidence = "Not attempted after a required lifecycle action failed.";
        }
      }
      if (!record.manual_fallbacks.includes(action.description)) record.manual_fallbacks.push(action.description);
    } else if (options.status === "failed") {
      record.warnings.push(`${action.id}: ${evidence}`);
    }
    record.status = overallStatus(record.actions);
    record.updated_at = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    const manifest2 = await readJson(manifestPath);
    syncManifestEvent(manifest2, record, recordPath2, workspaceRoot20);
    await assertValid("activity-lifecycle-record", record);
    await assertValid("runtime-manifest", manifest2);
    await writeJsonAtomic(recordPath2, record);
    await writeJsonAtomic(manifestPath, manifest2);
    return record;
  });
}
var import_yaml5;
var init_activity_lifecycle = __esm({
  "scripts/lib/activity-lifecycle.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    import_yaml5 = __toESM(require_dist(), 1);
    init_io();
    init_validation();
  }
});

// scripts/lib/plans.ts
import { createHash as createHash3, randomUUID } from "node:crypto";
import { lstat as lstat5, mkdir as mkdir3, readdir as readdir3, readFile as readFile8, realpath as realpath4, rename as rename2, rm } from "node:fs/promises";
import { basename, join as join7, resolve as resolve10 } from "node:path";
function contractMessages(errors2) {
  return errors2.map((error) => `${error.instancePath || "/"} ${error.message}`);
}
function markdownList(values20, empty) {
  return values20.length > 0 ? values20.map((value2) => `- ${value2}`).join("\n") : `- ${empty}`;
}
function productImpactBody(declaration) {
  const references = markdownList(declaration.references, "None referenced.");
  const proposed = declaration.proposed_change ?? "No product behavior change is proposed.";
  return `- Impact: ${declaration.impact}

Referenced Product Knowledge:

${references}

Proposed change:

${proposed}`;
}
function assertMarkdownCell(value2, field) {
  if (value2.includes("|") || /[\r\n]/.test(value2)) throw new Error(`${field} cannot contain a table delimiter or newline`);
}
function allocateWorkItems(request4) {
  const ids = /* @__PURE__ */ new Map();
  request4.work_items.forEach((item, index) => {
    const suffix = index === 0 ? 1 : index * 10;
    ids.set(item.key, `${request4.work_prefix}-${String(suffix).padStart(3, "0")}`);
  });
  return request4.work_items.map((item) => ({
    work_id: ids.get(item.key),
    title: item.title,
    parent: item.parent ? ids.get(item.parent) ?? null : null,
    depends_on: (item.depends_on ?? []).map((key) => ids.get(key) ?? key),
    area: item.area,
    repository: item.repository,
    scope: item.scope,
    test_scope: item.test_scope,
    test_policy: item.test_policy,
    ...item.test_rationale ? { test_rationale: item.test_rationale } : {},
    verification_commands: item.verification_commands,
    acceptance_criteria: item.acceptance_criteria,
    external_reference: null
  }));
}
function cycleErrors(items) {
  const dependencies = new Map(items.map((item) => [item.work_id, item.depends_on]));
  const errors2 = [];
  const visited = /* @__PURE__ */ new Set();
  const active = /* @__PURE__ */ new Set();
  const visit = (id) => {
    if (active.has(id)) {
      errors2.push(`work dependency cycle includes ${id}`);
      return;
    }
    if (visited.has(id)) return;
    active.add(id);
    for (const dependency of dependencies.get(id) ?? []) visit(dependency);
    active.delete(id);
    visited.add(id);
  };
  for (const id of dependencies.keys()) visit(id);
  return [...new Set(errors2)];
}
function parentCycleErrors(items) {
  return cycleErrors(items.map((item) => ({ work_id: item.work_id, depends_on: item.parent ? [item.parent] : [] }))).map((error) => error.replace("work dependency cycle", "work parent cycle"));
}
function planDraftSemanticErrors(request4, config) {
  const errors2 = [];
  const keys = /* @__PURE__ */ new Set();
  for (const item of request4.work_items) {
    if (keys.has(item.key)) errors2.push(`duplicate work item key: ${item.key}`);
    keys.add(item.key);
    for (const [field, value2] of [["title", item.title], ["area", item.area]]) {
      try {
        assertMarkdownCell(value2, `work item ${item.key} ${field}`);
      } catch (error) {
        errors2.push(error.message);
      }
    }
  }
  for (const item of request4.work_items) {
    if (item.parent && !keys.has(item.parent)) errors2.push(`work item ${item.key} has unknown parent: ${item.parent}`);
    if (item.parent === item.key) errors2.push(`work item ${item.key} cannot be its own parent`);
    for (const dependency of item.depends_on ?? []) {
      if (!keys.has(dependency)) errors2.push(`work item ${item.key} has unknown dependency: ${dependency}`);
      if (dependency === item.key) errors2.push(`work item ${item.key} cannot depend on itself`);
    }
  }
  const productKnowledge = request4.product_knowledge;
  if (productKnowledge) {
    const requiresChange = ["behavior-change", "new-workflow", "retired-workflow"];
    if (requiresChange.includes(productKnowledge.impact) && !productKnowledge.proposed_change?.trim()) {
      errors2.push(`product knowledge impact '${productKnowledge.impact}' requires a proposed_change summary`);
    }
    if (productKnowledge.impact === "none" && productKnowledge.proposed_change) {
      errors2.push("product knowledge impact 'none' must not include a proposed_change");
    }
  }
  const keyedDependencies = request4.work_items.map((item) => ({ work_id: item.key, depends_on: item.depends_on ?? [] }));
  errors2.push(...cycleErrors(keyedDependencies));
  errors2.push(...parentCycleErrors(request4.work_items.map((item) => ({ work_id: item.key, parent: item.parent ?? null }))));
  if (config) {
    for (const repository of request4.affected_repositories) {
      if (!config.repositories[repository]) errors2.push(`affected repository is not registered: ${repository}`);
    }
    for (const item of request4.work_items) {
      if (!config.repositories[item.repository]) errors2.push(`work item ${item.key} repository is not registered: ${item.repository}`);
      if (!request4.affected_repositories.includes(item.repository)) errors2.push(`work item ${item.key} repository is not affected: ${item.repository}`);
    }
  }
  return [...new Set(errors2)];
}
function planWorkBreakdownSemanticErrors(breakdown, config) {
  const errors2 = [];
  const ids = /* @__PURE__ */ new Set();
  for (const item of breakdown.items) {
    if (ids.has(item.work_id)) errors2.push(`duplicate work ID: ${item.work_id}`);
    ids.add(item.work_id);
    if (!item.work_id.startsWith(`${breakdown.work_prefix}-`)) {
      errors2.push(`work ID does not use ${breakdown.work_prefix} prefix: ${item.work_id}`);
    }
  }
  for (const item of breakdown.items) {
    if (item.parent && !ids.has(item.parent)) errors2.push(`${item.work_id} has unknown parent: ${item.parent}`);
    if (item.parent === item.work_id) errors2.push(`${item.work_id} cannot be its own parent`);
    for (const dependency of item.depends_on) {
      if (!ids.has(dependency)) errors2.push(`${item.work_id} has unknown dependency: ${dependency}`);
      if (dependency === item.work_id) errors2.push(`${item.work_id} cannot depend on itself`);
    }
    if (config && !config.repositories[item.repository]) errors2.push(`${item.work_id} repository is not registered: ${item.repository}`);
  }
  errors2.push(...cycleErrors(breakdown.items));
  errors2.push(...parentCycleErrors(breakdown.items));
  return [...new Set(errors2)];
}
function materialDigest(files, names) {
  const hash = createHash3("sha256");
  for (const name of names) hash.update(`${name}\0${files.get(name) ?? ""}\0`);
  return `sha256:${hash.digest("hex")}`;
}
function parsePlanIndex(raw) {
  const match = raw.match(/^---\r?\n([\s\S]*?)\r?\n---(?:\r?\n|$)/);
  if (!match) throw new Error("Plan README must begin with YAML frontmatter");
  return (0, import_yaml6.parse)(match[1]);
}
function parseWorkBreakdown(raw, index, config) {
  const lines = raw.replace(/\r\n/g, "\n").split("\n");
  const currentHeader = lines.indexOf(tableHeader);
  const legacyHeader = lines.indexOf(legacyTableHeader);
  const legacy = currentHeader === -1 && legacyHeader !== -1;
  const header = currentHeader === -1 ? legacyHeader : currentHeader;
  if (header === -1 || lines[header + 1] !== (legacy ? legacyTableSeparator : tableSeparator)) {
    throw new Error("Work breakdown must contain the canonical seven-column table and must not add live status columns");
  }
  const summaries = [];
  for (const line of lines.slice(header + 2)) {
    if (!line.startsWith("|")) break;
    const cells = line.slice(1, -1).split("|").map((cell) => cell.trim());
    if (cells.length !== (legacy ? 6 : 7)) throw new Error(`Invalid work breakdown row: ${line}`);
    const [workId, title, parent, dependencies] = cells;
    const repository = legacy ? null : cells[4];
    const area = cells[legacy ? 4 : 5];
    const external = cells[legacy ? 5 : 6];
    summaries.push({
      work_id: workId,
      title,
      parent: parent === "\u2014" ? null : parent,
      depends_on: dependencies === "\u2014" ? [] : dependencies.split(",").map((value2) => value2.trim()),
      repository,
      area,
      external_reference: external === "\u2014" ? null : external
    });
  }
  const executionMatch = raw.match(/## Execution contracts\r?\n\r?\n```json\r?\n([\s\S]*?)\r?\n```/);
  if (!executionMatch) throw new Error("Work breakdown must contain the canonical execution contracts JSON block");
  const execution = JSON.parse(executionMatch[1]);
  if (![1, 2].includes(execution.contract_version) || !Array.isArray(execution.items)) throw new Error("Invalid work execution contracts block");
  if (!legacy && execution.contract_version !== 2) throw new Error("The canonical seven-column work breakdown requires execution contract version 2");
  if (legacy && execution.contract_version !== 1) throw new Error("The legacy six-column work breakdown requires execution contract version 1");
  const executionById = new Map(execution.items.map((item) => [item.work_id, item]));
  const items = summaries.map((summary2) => {
    const details = executionById.get(summary2.work_id);
    if (!details) throw new Error(`Missing execution contract for ${summary2.work_id}`);
    let repository = details.repository ?? summary2.repository;
    if (summary2.repository && details.repository && summary2.repository !== details.repository) {
      throw new Error(`Repository mismatch for ${summary2.work_id}: table has ${summary2.repository}, execution contract has ${details.repository}`);
    }
    if (!repository && legacy && execution.contract_version === 1) {
      if (!config) throw new Error(`Legacy work item ${summary2.work_id} has no repository; validate it inside a configured workspace or migrate the plan`);
      repository = config.repositories[summary2.area] ? summary2.area : null;
      if (!repository) throw new Error(`Legacy work item ${summary2.work_id} has no repository and area '${summary2.area}' is not an exact registered repository key; add an explicit repository through a material plan revision`);
    }
    if (!repository) throw new Error(`Work item ${summary2.work_id} has no explicit repository`);
    return { ...summary2, ...details, repository };
  });
  for (const workId of executionById.keys()) if (!summaries.some((item) => item.work_id === workId)) throw new Error(`Execution contract references unknown work ID: ${workId}`);
  return { contract_version: 2, plan_id: index.plan_id, work_prefix: index.work_prefix, items };
}
async function regularFile(path2) {
  try {
    const info = await lstat5(path2);
    return info.isFile() && !info.isSymbolicLink();
  } catch {
    return false;
  }
}
async function validatePlanDirectory(planDirectory3, expectedPlanId = basename(planDirectory3)) {
  const directory = resolve10(planDirectory3);
  const errors2 = [];
  let index = null;
  let breakdown = null;
  try {
    const workspaceRoot20 = resolve10(directory, "../../..");
    const config = await readData(join7(workspaceRoot20, "workspace.yaml"));
    const workspaceErrors = contractMessages(await validateContract("workspace", config));
    workspaceErrors.push(...workspaceSemanticErrors(config));
    if (workspaceErrors.length > 0) throw new Error(`Invalid workspace configuration: ${workspaceErrors.join("; ")}`);
    const info = await lstat5(directory);
    if (!info.isDirectory() || info.isSymbolicLink()) throw new Error("Plan path must be a real directory");
    if (!await regularFile(join7(directory, "README.md"))) throw new Error("Plan README must be a real file");
    index = parsePlanIndex(await readFile8(join7(directory, "README.md"), "utf8"));
    const indexErrors = contractMessages(await validateContract("plan-index", index));
    errors2.push(...indexErrors);
    if (indexErrors.length > 0) return { index, work_breakdown: null, errors: [...new Set(errors2)] };
    if (index.plan_id !== expectedPlanId) errors2.push(`plan_id must match directory name: ${expectedPlanId}`);
    if (Date.parse(index.updated_at) < Date.parse(index.created_at)) errors2.push("updated_at cannot be earlier than created_at");
    const sorted = [...index.documents].sort();
    if (JSON.stringify(sorted) !== JSON.stringify(index.documents)) errors2.push("numbered plan documents must be listed in ascending order");
    if (!index.documents.includes(index.work_breakdown)) errors2.push("work_breakdown must be listed in documents");
    const actualNumbered = (await readdir3(directory)).filter((name) => /^[0-9]{4}-.+\.md$/.test(name)).sort();
    for (const document of actualNumbered) {
      if (!index.documents.includes(document)) errors2.push(`numbered plan document is not listed in the index: ${document}`);
    }
    for (const document of index.documents) {
      if (!await regularFile(join7(directory, document))) errors2.push(`plan document is missing or unsafe: ${document}`);
    }
    if (errors2.length === 0) {
      const material = /* @__PURE__ */ new Map();
      for (const document of index.documents) material.set(document, await readFile8(join7(directory, document), "utf8"));
      const digest = materialDigest(material, index.documents);
      if (index.material_digest !== digest) errors2.push("material_digest does not match the numbered plan documents");
      if (index.status === "approved" && index.approved_digest !== digest) errors2.push("approved_digest does not match the approved plan material");
      breakdown = parseWorkBreakdown(material.get(index.work_breakdown), index, config);
      errors2.push(...contractMessages(await validateContract("plan-work-breakdown", breakdown)));
      errors2.push(...planWorkBreakdownSemanticErrors(breakdown, config));
    }
  } catch (error) {
    errors2.push(error.message);
  }
  return { index, work_breakdown: breakdown, errors: [...new Set(errors2)] };
}
async function actualMaterialDigest(directory, index) {
  const material = /* @__PURE__ */ new Map();
  for (const document of index.documents) material.set(document, await readFile8(join7(directory, document), "utf8"));
  return materialDigest(material, index.documents);
}
async function setPlanState(planDirectory3, transition2, now = /* @__PURE__ */ new Date()) {
  const directory = resolve10(planDirectory3);
  const validation = await validatePlanDirectory(directory);
  const allowedStaleDigestErrors = /* @__PURE__ */ new Set([
    "material_digest does not match the numbered plan documents",
    "approved_digest does not match the approved plan material"
  ]);
  const blocking = validation.errors.filter((error) => !allowedStaleDigestErrors.has(error));
  if (!validation.index || blocking.length > 0) throw new Error(`Plan state transition validation failed:
- ${blocking.join("\n- ")}`);
  const index = validation.index;
  const digest = await actualMaterialDigest(directory, index);
  if (transition2.kind === "approve") {
    if (index.status !== "draft") throw new Error("Only a draft plan can be approved");
    if (!transition2.approved_by.trim()) throw new Error("Approval requires a non-empty approver");
    index.status = "approved";
    index.approved_at = now.toISOString();
    index.approved_by = transition2.approved_by.trim();
    index.material_digest = digest;
    index.approved_digest = digest;
  } else if (transition2.kind === "material-revision") {
    if (index.status !== "approved") throw new Error("Material revision transition requires an approved plan");
    if (!transition2.reason.trim()) throw new Error("Material revision requires a reason");
    index.status = "draft";
    index.plan_version += 1;
    index.approved_at = null;
    index.approved_by = null;
    index.approved_digest = null;
    index.material_digest = digest;
    index.revision_reason = transition2.reason.trim();
  } else {
    if (index.status !== "approved") throw new Error("Non-material repair transition requires an approved plan");
    index.material_digest = digest;
    index.approved_digest = digest;
  }
  index.updated_at = now.toISOString();
  const readmePath = join7(directory, "README.md");
  const raw = await readFile8(readmePath, "utf8");
  const match = raw.match(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/);
  if (!match) throw new Error("Plan README must begin with YAML frontmatter");
  await writeTextAtomic(readmePath, raw.replace(match[0], `---
${(0, import_yaml6.stringify)(index).trimEnd()}
---
`));
  const after = await validatePlanDirectory(directory);
  if (after.errors.length > 0) throw new Error(`Plan state transition produced invalid metadata:
- ${after.errors.join("\n- ")}`);
  return index;
}
function renderDocument(title, sections) {
  return `# ${title}

${sections.map(([heading, body]) => `## ${heading}

${body}`).join("\n\n")}
`;
}
function renderPlan(request4, createdAt) {
  const workItems = allocateWorkItems(request4);
  const breakdown = { contract_version: 2, plan_id: request4.plan_id, work_prefix: request4.work_prefix, items: workItems };
  const files = /* @__PURE__ */ new Map();
  const overviewSections = [
    ["Summary", request4.summary],
    ["Source", `${request4.source.kind}: ${request4.source.reference}`],
    ["Affected repositories", markdownList(request4.affected_repositories, "None identified.")]
  ];
  if (request4.product_knowledge) overviewSections.push(["Product impact", productImpactBody(request4.product_knowledge)]);
  overviewSections.push(
    ["Assumptions", markdownList(request4.assumptions, "None recorded.")],
    ["Open questions", markdownList(request4.open_questions, "None recorded.")]
  );
  files.set("0001-overview.md", renderDocument("Overview", overviewSections));
  files.set("0010-requirements.md", renderDocument("Requirements", [["Requirements and acceptance criteria", markdownList(request4.requirements, "None recorded.")]]));
  files.set("0020-solution.md", renderDocument("Solution", [["Proposed solution", markdownList(request4.solution, "None recorded.")]]));
  files.set("0040-delivery.md", renderDocument("Delivery", [["Delivery order", markdownList(request4.delivery, "None recorded.")]]));
  files.set("0050-verification.md", renderDocument("Verification", [["Verification strategy", markdownList(request4.verification, "None recorded.")]]));
  files.set("0070-risks.md", renderDocument("Risks", [["Risks and mitigations", markdownList(request4.risks, "None recorded.")]]));
  const rows = workItems.map((item) => `| ${item.work_id} | ${item.title} | ${item.parent ?? "\u2014"} | ${item.depends_on.join(", ") || "\u2014"} | ${item.repository} | ${item.area} | \u2014 |`).join("\n");
  const execution = {
    contract_version: 2,
    items: workItems.map(({ work_id, repository, scope, test_scope, test_policy, test_rationale, verification_commands, acceptance_criteria }) => ({
      work_id,
      repository,
      scope,
      test_scope,
      test_policy,
      ...test_rationale ? { test_rationale } : {},
      verification_commands,
      acceptance_criteria
    }))
  };
  files.set("0080-work-breakdown.md", `# Work breakdown

${tableHeader}
${tableSeparator}
${rows}

## Execution contracts

\`\`\`json
${JSON.stringify(execution, null, 2)}
\`\`\`

Live task status does not belong in this plan. Add confirmed external references only after an explicit publication action.
`);
  const index = {
    contract_version: 1,
    plan_id: request4.plan_id,
    title: request4.title,
    status: "draft",
    plan_version: 1,
    approved_at: null,
    approved_by: null,
    revision_reason: "Initial draft",
    source: request4.source,
    work_prefix: request4.work_prefix,
    documents: [...documents],
    work_breakdown: "0080-work-breakdown.md",
    material_digest: materialDigest(files, documents),
    approved_digest: null,
    created_at: createdAt,
    updated_at: createdAt,
    ...request4.product_knowledge ? { product_knowledge: request4.product_knowledge } : {}
  };
  const links = documents.map((document) => `- [${document.replace(/^[0-9]{4}-|\.md$/g, "").replaceAll("-", " ")}](./${document})`).join("\n");
  files.set("README.md", `---
${(0, import_yaml6.stringify)(index).trimEnd()}
---

# ${request4.title}

${request4.summary}

## Plan documents

${links}

## Approval gate

Human approval must explicitly cover scope, solution, delivery order, risks, and acceptance criteria before the metadata status changes to \`approved\`. The machine-readable frontmatter status is authoritative; approval updates metadata without rewriting this prose.
`);
  return { index, breakdown, files };
}
async function createPlanDraft(workspaceRootInput, request4, now = /* @__PURE__ */ new Date()) {
  const workspaceRoot20 = resolve10(workspaceRootInput);
  const contractErrors2 = contractMessages(await validateContract("plan-draft-request", request4));
  const config = await readData(join7(workspaceRoot20, "workspace.yaml"));
  const workspaceErrors = contractMessages(await validateContract("workspace", config));
  const errors2 = [...contractErrors2, ...workspaceErrors];
  if (contractErrors2.length === 0 && workspaceErrors.length === 0) {
    errors2.push(...workspaceSemanticErrors(config), ...planDraftSemanticErrors(request4, config));
  }
  if (errors2.length > 0) throw new Error(`Invalid plan draft request:
- ${errors2.join("\n- ")}`);
  const realWorkspace = await realpath4(workspaceRoot20);
  const contextRoot = assertInside(workspaceRoot20, join7(workspaceRoot20, "context"));
  try {
    const info = await lstat5(contextRoot);
    if (!info.isDirectory() || info.isSymbolicLink()) throw new Error(`Context root must be a real directory: ${contextRoot}`);
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
    await mkdir3(contextRoot, { mode: 493 });
  }
  if (await realpath4(contextRoot) !== join7(realWorkspace, "context")) throw new Error(`Context root must not traverse symbolic links: ${contextRoot}`);
  const plansRoot = assertInside(contextRoot, join7(contextRoot, "plans"));
  try {
    const info = await lstat5(plansRoot);
    if (!info.isDirectory() || info.isSymbolicLink()) throw new Error(`Plan root must be a real directory: ${plansRoot}`);
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
    await mkdir3(plansRoot, { mode: 493 });
  }
  const realPlansRoot = await realpath4(plansRoot);
  assertInside(realWorkspace, realPlansRoot);
  if (realPlansRoot !== join7(realWorkspace, "context", "plans")) {
    throw new Error(`Plan root must not traverse symbolic links: ${plansRoot}`);
  }
  const destination = assertInside(realPlansRoot, join7(realPlansRoot, request4.plan_id));
  try {
    await lstat5(destination);
    throw new Error(`Plan already exists; refusing to overwrite: ${destination}`);
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
  }
  const temporary = join7(realPlansRoot, `.${request4.plan_id}.${randomUUID()}.tmp`);
  const rendered = renderPlan(request4, now.toISOString());
  try {
    await mkdir3(temporary, { mode: 493 });
    for (const [name, contents] of rendered.files) await writeTextExclusive(join7(temporary, name), contents);
    const validation = await validatePlanDirectory(temporary, request4.plan_id);
    if (validation.errors.length > 0) throw new Error(`Generated plan failed validation:
- ${validation.errors.join("\n- ")}`);
    await rename2(temporary, destination);
  } catch (error) {
    await rm(temporary, { recursive: true, force: true });
    throw error;
  }
  return {
    plan_id: request4.plan_id,
    status: "draft",
    plan_version: 1,
    directory: destination,
    index: join7(destination, "README.md"),
    documents: [...documents],
    work_ids: rendered.breakdown.items.map((item) => item.work_id),
    approval_required: true
  };
}
var import_yaml6, documents, tableHeader, tableSeparator, legacyTableHeader, legacyTableSeparator;
var init_plans = __esm({
  "scripts/lib/plans.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    import_yaml6 = __toESM(require_dist(), 1);
    init_io();
    init_validation();
    documents = [
      "0001-overview.md",
      "0010-requirements.md",
      "0020-solution.md",
      "0040-delivery.md",
      "0050-verification.md",
      "0070-risks.md",
      "0080-work-breakdown.md"
    ];
    tableHeader = "| Work ID | Title | Parent | Depends on | Repository | Area | External reference |";
    tableSeparator = "| --- | --- | --- | --- | --- | --- | --- |";
    legacyTableHeader = "| Work ID | Title | Parent | Depends on | Area | External reference |";
    legacyTableSeparator = "| --- | --- | --- | --- | --- | --- |";
  }
});

// scripts/lib/run-task.ts
import { randomBytes as randomBytes2 } from "node:crypto";
import { access as access4, readFile as readFile9 } from "node:fs/promises";
import { join as join8, relative as relative5, resolve as resolve11 } from "node:path";
async function assertValid2(name, value2) {
  const errors2 = await validateContract(name, value2);
  if (errors2.length > 0) {
    throw new Error(`Generated ${name} is invalid: ${errors2.map((error) => `${error.instancePath} ${error.message}`).join("; ")}`);
  }
}
function repositoryExpectation(input) {
  const paths2 = normalizeScope(input.test_scope, `test scope for ${input.name}`);
  if ((input.test_policy === "required" || input.test_policy === "existing-coverage") && paths2.length === 0) {
    throw new Error(`${input.test_policy} test policy requires at least one test scope entry for ${input.name}`);
  }
  return { policy: input.test_policy, paths: paths2, rationale: input.test_rationale?.trim() || defaultTestRationale(input.test_policy) };
}
function normalizeContractFirstRequest(input, workId, runId, createdAt) {
  const request4 = input.request.trim();
  const acceptanceCriteria = input.acceptance_criteria.map((item) => item.trim()).filter(Boolean);
  if (!request4) throw new Error("A direct request is required");
  if (acceptanceCriteria.length === 0) throw new Error("At least one acceptance criterion is required");
  const names = input.repositories.map((repository) => repository.name);
  if (new Set(names).size !== names.length) throw new Error("Repository names must be unique");
  if (!names.includes(input.shared_contract.repository)) throw new Error("Shared contract repository must be included in repositories");
  const contractPaths = normalizeScope(input.shared_contract.paths, "shared contract paths");
  const byName = new Map(input.repositories.map((repository) => [repository.name, repository]));
  for (const repository of input.repositories) {
    if (repository.depends_on.includes(repository.name)) throw new Error(`${repository.name} cannot depend on itself`);
    for (const dependency of repository.depends_on) if (!byName.has(dependency)) throw new Error(`${repository.name} has unknown dependency ${dependency}`);
  }
  const contractOwner = byName.get(input.shared_contract.repository);
  if (contractOwner.depends_on.length > 0) throw new Error("Shared contract repository cannot depend on another repository");
  const visiting = /* @__PURE__ */ new Set();
  const orders = /* @__PURE__ */ new Map();
  const orderOf = (name) => {
    const known = orders.get(name);
    if (known !== void 0) return known;
    if (visiting.has(name)) throw new Error(`Repository dependency cycle includes ${name}`);
    visiting.add(name);
    const repository = byName.get(name);
    const order = repository.depends_on.length === 0 ? 0 : Math.max(...repository.depends_on.map(orderOf)) + 1;
    visiting.delete(name);
    orders.set(name, order);
    return order;
  };
  const dependsOnContract = (name, seen = /* @__PURE__ */ new Set()) => {
    if (name === input.shared_contract.repository) return true;
    if (seen.has(name)) return false;
    seen.add(name);
    return byName.get(name).depends_on.some((dependency) => dependsOnContract(dependency, seen));
  };
  for (const name of names) {
    orderOf(name);
    if (name !== input.shared_contract.repository && !dependsOnContract(name)) throw new Error(`${name} must depend on the shared contract repository`);
  }
  const targets = input.repositories.map((repository) => {
    const implementationScope = normalizeScope(repository.scope, `implementation scope for ${repository.name}`);
    const testExpectation = repositoryExpectation(repository);
    const scope = [.../* @__PURE__ */ new Set([...implementationScope, ...testExpectation.policy === "required" ? testExpectation.paths : []])];
    const repositoryAcceptance = repository.acceptance_criteria.map((criterion) => criterion.trim()).filter(Boolean);
    if (repositoryAcceptance.length === 0) throw new Error(`At least one acceptance criterion is required for ${repository.name}`);
    return {
      name: repository.name,
      dependency_order: orders.get(repository.name),
      depends_on: repository.depends_on,
      scope,
      implementation_scope: implementationScope,
      test_expectation: testExpectation,
      verification_commands: repository.verification_commands.map((command3) => command3.trim()).filter(Boolean),
      acceptance_criteria: repositoryAcceptance
    };
  }).sort((left, right) => left.dependency_order - right.dependency_order || left.name.localeCompare(right.name));
  const ownerScope = targets.find((target) => target.name === input.shared_contract.repository).scope;
  for (const path2 of contractPaths) if (!ownerScope.some((scope) => path2 === scope || path2.startsWith(`${scope}/`))) {
    throw new Error(`Shared contract path is outside ${input.shared_contract.repository} scope: ${path2}`);
  }
  return {
    contract_version: 1,
    work_id: workId,
    run_id: runId,
    source: { kind: "direct-request" },
    requested_outcome: request4,
    scope: [...new Set(targets.flatMap((target) => target.scope))],
    acceptance_criteria: acceptanceCriteria,
    repositories: targets,
    shared_contract: { repository: input.shared_contract.repository, paths: contractPaths },
    plan: { reference: null, approval_state: "not-applicable" },
    activity: { reference: null, claim_status: "not-applicable", duplicate_effort_warning: true },
    assumptions: ["The direct request is authoritative for this planless run.", "Dependent repositories remain locked until their declared dependencies pass independent verification."],
    risks: ["No authoritative claim is available; duplicate effort is possible."],
    verification_commands: [...new Set(targets.flatMap((target) => target.verification_commands ?? []))],
    authorization: { kind: "explicit-user-request", evidence: "The human explicitly invoked run-task for this direct request." },
    created_at: createdAt
  };
}
function normalizePlanRequest(input, item, runId, createdAt) {
  const implementationScope = normalizeScope(item.scope, `implementation scope for ${item.work_id}`);
  const testExpectation = repositoryExpectation({
    name: item.repository,
    depends_on: [],
    scope: item.scope,
    test_scope: item.test_scope,
    test_policy: item.test_policy,
    ...item.test_rationale ? { test_rationale: item.test_rationale } : {},
    verification_commands: item.verification_commands,
    acceptance_criteria: item.acceptance_criteria
  });
  const scope = [.../* @__PURE__ */ new Set([...implementationScope, ...testExpectation.policy === "required" ? testExpectation.paths : []])];
  const target = {
    name: item.repository,
    dependency_order: 0,
    depends_on: [],
    scope,
    implementation_scope: implementationScope,
    test_expectation: testExpectation,
    verification_commands: item.verification_commands,
    acceptance_criteria: item.acceptance_criteria
  };
  return {
    contract_version: 1,
    work_id: input.work_ids[0],
    run_id: runId,
    source: { kind: "plan", reference: input.source.reference },
    requested_outcome: item.title,
    scope,
    acceptance_criteria: item.acceptance_criteria,
    repositories: [target],
    plan: {
      reference: input.source.reference,
      approval_state: "approved",
      plan_version: input.source.plan_version,
      approved_digest: input.source.approved_digest,
      work_ids: input.work_ids
    },
    activity: { reference: null, claim_status: "not-applicable", duplicate_effort_warning: true },
    assumptions: ["The selected work IDs and approved plan material are authoritative for this run."],
    risks: ["No authoritative claim is available; duplicate effort is possible."],
    verification_commands: item.verification_commands,
    authorization: { kind: "confirmed-selection", evidence: `The human selected approved plan work: ${input.work_ids.join(", ")}.` },
    created_at: createdAt
  };
}
function normalizeDirectRequest(input) {
  const request4 = input.request.trim();
  const acceptanceCriteria = input.acceptanceCriteria.map((item) => item.trim()).filter(Boolean);
  const implementationScope = normalizeScope(input.scope, "implementation scope");
  const testScope = normalizeScope(input.testScope ?? [], "test scope");
  const testPolicy = input.testPolicy ?? (testScope.length > 0 ? "required" : "verifier-only");
  if (!request4) throw new Error("A direct request is required");
  if (acceptanceCriteria.length === 0) throw new Error("At least one acceptance criterion is required");
  if (implementationScope.length === 0) throw new Error("At least one implementation scope entry is required");
  if ((testPolicy === "required" || testPolicy === "existing-coverage") && testScope.length === 0) {
    throw new Error(`${testPolicy} test policy requires at least one test scope entry`);
  }
  const testExpectation = {
    policy: testPolicy,
    paths: testScope,
    rationale: input.testRationale?.trim() || defaultTestRationale(testPolicy)
  };
  const scope = [.../* @__PURE__ */ new Set([...implementationScope, ...testPolicy === "required" ? testScope : []])];
  return {
    contract_version: 1,
    work_id: input.workId,
    run_id: input.runId,
    source: { kind: "direct-request" },
    requested_outcome: request4,
    scope,
    implementation_scope: implementationScope,
    test_expectation: testExpectation,
    acceptance_criteria: acceptanceCriteria,
    repositories: [{ name: input.repository, dependency_order: 0 }],
    plan: { reference: null, approval_state: "not-applicable" },
    activity: {
      reference: null,
      claim_status: "not-applicable",
      duplicate_effort_warning: true
    },
    assumptions: ["The direct request is authoritative for this planless run."],
    risks: ["No authoritative claim is available; duplicate effort is possible."],
    verification_commands: input.verificationCommands,
    authorization: {
      kind: "explicit-user-request",
      evidence: "The human explicitly invoked run-task for this direct request."
    },
    created_at: input.createdAt
  };
}
function normalizeScope(values20, label) {
  const normalized = values20.map((item) => item.trim().replace(/\/$/, "")).filter(Boolean);
  for (const path2 of normalized) {
    if (path2.startsWith("/") || path2.includes("\\") || path2.split("/").includes("..")) {
      throw new Error(`${label} entries must be repository-relative paths: ${path2}`);
    }
  }
  return [...new Set(normalized)];
}
function defaultTestRationale(policy) {
  switch (policy) {
    case "required":
      return "The worker must add or update tests in the declared test scope.";
    case "existing-coverage":
      return "Declared existing tests are expected to cover the requested behavior.";
    case "not-required":
      return "No repository test change is required for this task.";
    default:
      return "No test edit scope is authorized; the verifier must supply independent acceptance evidence.";
  }
}
async function preparePlanlessTask(options) {
  const workspaceRoot20 = resolve11(options.workspaceRoot);
  const configPath = join8(workspaceRoot20, "workspace.yaml");
  const config = (0, import_yaml7.parse)(await readFile9(configPath, "utf8"));
  await assertValid2("workspace", config);
  const semanticErrors = workspaceSemanticErrors(config);
  if (semanticErrors.length > 0) throw new Error(`Invalid workspace: ${semanticErrors.join("; ")}`);
  const repositoryConfig = config.repositories[options.repository];
  if (!repositoryConfig) throw new Error(`Unknown repository: ${options.repository}`);
  const repositoryPath = assertInside(workspaceRoot20, join8(workspaceRoot20, repositoryConfig.path));
  await access4(repositoryPath);
  await assertCleanRepository(repositoryPath);
  const baseCommit = await git(repositoryPath, ["rev-parse", repositoryConfig.default_branch]);
  const requiredInstructionPaths = [
    join8(workspaceRoot20, "AGENTS.md"),
    join8(workspaceRoot20, "agents", `${repositoryConfig.agent}.md`),
    join8(workspaceRoot20, "agents", "repository-worker.md"),
    join8(workspaceRoot20, "agents", "verifier.md")
  ];
  await Promise.all(requiredInstructionPaths.map((path2) => access4(path2)));
  const runtimeRoot = assertInside(workspaceRoot20, join8(workspaceRoot20, ".runtime"));
  const now = options.now ?? /* @__PURE__ */ new Date();
  const { workId, runId } = await generateIds(runtimeRoot, options.request, now, options.discriminator);
  const createdAt = now.toISOString();
  const branch = `agent/${workId.toLowerCase()}-${slugify2(options.request)}-${runId.slice(-8)}`;
  const runRoot = assertInside(runtimeRoot, join8(runtimeRoot, "runs", runId));
  const worktree = assertInside(runtimeRoot, join8(runtimeRoot, "worktrees", runId, options.repository));
  const taskBriefPath = join8(runtimeRoot, "tasks", `${runId}.json`);
  const manifestPath = join8(runRoot, "manifest.json");
  const workerInputPath = join8(runRoot, `${options.repository}-worker-input.json`);
  const verifierInputPath = join8(runRoot, `${options.repository}-verifier-input.json`);
  const workerResultPath = join8(runtimeRoot, "results", `${runId}-${options.repository}-worker.json`);
  const verifierResultPath = join8(runtimeRoot, "results", `${runId}-${options.repository}-verifier.json`);
  const taskBrief = normalizeDirectRequest({
    request: options.request,
    repository: options.repository,
    acceptanceCriteria: options.acceptanceCriteria,
    scope: options.scope,
    ...options.testScope ? { testScope: options.testScope } : {},
    ...options.testPolicy ? { testPolicy: options.testPolicy } : {},
    ...options.testRationale ? { testRationale: options.testRationale } : {},
    verificationCommands: options.verificationCommands ?? [],
    workId,
    runId,
    createdAt
  });
  await assertValid2("task-brief", taskBrief);
  if (taskBrief.test_expectation?.policy === "existing-coverage") {
    for (const path2 of taskBrief.test_expectation.paths) {
      try {
        await access4(assertInside(repositoryPath, join8(repositoryPath, path2)));
      } catch {
        throw new Error(`Existing-coverage test path does not exist: ${path2}`);
      }
    }
  }
  await writeJsonAtomic(taskBriefPath, taskBrief);
  const runtimeRepository = {
    name: options.repository,
    base_path: relative5(workspaceRoot20, repositoryPath),
    base_commit: baseCommit,
    branch,
    worktree,
    worker_input: workerInputPath,
    verifier_input: verifierInputPath,
    repair_attempts: 0
  };
  const manifest2 = {
    contract_version: 1,
    work_id: workId,
    run_id: runId,
    source_kind: "direct-request",
    status: "preparing",
    created_at: createdAt,
    updated_at: createdAt,
    task_brief: taskBriefPath,
    repositories: [runtimeRepository],
    evidence: [taskBriefPath, manifestPath, workerInputPath, verifierInputPath],
    warnings: config.activity.provider === "none" ? ["No activity tool is configured; this run cannot guarantee exclusive ownership."] : [],
    execution_events: [],
    lifecycle_events: config.activity.provider === "none" ? [{
      event: "task.starting",
      status: "skipped",
      idempotency_key: `${runId}:task.starting:activity-none`,
      occurred_at: createdAt
    }] : []
  };
  await assertValid2("runtime-manifest", manifest2);
  await writeJsonAtomic(manifestPath, manifest2);
  const instructionPaths = [
    ...requiredInstructionPaths.slice(0, 3)
  ];
  try {
    await access4(join8(repositoryPath, "AGENTS.md"));
    instructionPaths.push(join8(worktree, "AGENTS.md"));
  } catch {
  }
  const workerInput = {
    contract_version: 1,
    role: "repository-worker",
    task_brief: taskBriefPath,
    repository: options.repository,
    worktree,
    branch,
    base_commit: baseCommit,
    ready: true,
    blocked_by: [],
    allowed_scope: taskBrief.scope,
    implementation_scope: taskBrief.implementation_scope,
    test_expectation: taskBrief.test_expectation,
    instruction_paths: instructionPaths,
    result_contract: join8(workspaceRoot20, ".agents", "contracts", "worker-result.schema.json"),
    result_path: workerResultPath
  };
  const verifierInput = {
    contract_version: 1,
    role: "verifier",
    read_only: true,
    task_brief: taskBriefPath,
    repository: options.repository,
    worktree,
    branch,
    base_commit: baseCommit,
    worker_result: workerResultPath,
    acceptance_criteria: taskBrief.acceptance_criteria,
    test_expectation: taskBrief.test_expectation,
    verification_commands: taskBrief.verification_commands,
    instruction_paths: [join8(workspaceRoot20, "AGENTS.md"), join8(workspaceRoot20, "agents", "verifier.md"), ...instructionPaths.slice(3)],
    result_contract: join8(workspaceRoot20, ".agents", "contracts", "verifier-result.schema.json"),
    result_path: verifierResultPath
  };
  await writeJsonAtomic(workerInputPath, workerInput);
  await writeJsonAtomic(verifierInputPath, verifierInput);
  if (config.activity.provider !== "none") {
    const lifecycle = await prepareActivityLifecycle({
      workspaceRoot: workspaceRoot20,
      runId,
      event: "task.starting",
      availableCapabilities: options.availableCapabilities ?? [],
      now
    });
    if (lifecycle.status !== "completed" && lifecycle.status !== "skipped") {
      return {
        workId,
        runId,
        branch,
        worktree,
        taskBrief: taskBriefPath,
        manifest: manifestPath,
        workerInput: workerInputPath,
        verifierInput: verifierInputPath,
        repositories: [{ name: options.repository, branch, worktree, workerInput: workerInputPath, verifierInput: verifierInputPath, ready: true, blockedBy: [] }],
        preparationStatus: lifecycle.status === "failed" ? "blocked" : "awaiting-activity"
      };
    }
  }
  try {
    await ensurePrivateDirectory(join8(runtimeRoot, "worktrees", runId));
    await git(repositoryPath, ["worktree", "add", "-b", branch, worktree, baseCommit]);
    runtimeRepository.status = "prepared";
    manifest2.status = "prepared";
    manifest2.updated_at = (/* @__PURE__ */ new Date()).toISOString();
    await assertValid2("runtime-manifest", manifest2);
    await writeJsonAtomic(manifestPath, manifest2);
  } catch (error) {
    manifest2.status = "blocked";
    manifest2.updated_at = (/* @__PURE__ */ new Date()).toISOString();
    manifest2.evidence.push(`Preparation failed: ${error.message}`);
    await writeJsonAtomic(manifestPath, manifest2);
    throw error;
  }
  return {
    workId,
    runId,
    branch,
    worktree,
    taskBrief: taskBriefPath,
    manifest: manifestPath,
    workerInput: workerInputPath,
    verifierInput: verifierInputPath,
    repositories: [{ name: options.repository, branch, worktree, workerInput: workerInputPath, verifierInput: verifierInputPath, ready: true, blockedBy: [] }],
    preparationStatus: "prepared"
  };
}
async function prepareContractFirstTask(options) {
  const requestErrors = await validateContract("run-task-request", options.request);
  if (requestErrors.length > 0) {
    throw new Error(`Invalid run-task-request: ${requestErrors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
  }
  const workspaceRoot20 = resolve11(options.workspaceRoot);
  const config = (0, import_yaml7.parse)(await readFile9(join8(workspaceRoot20, "workspace.yaml"), "utf8"));
  await assertValid2("workspace", config);
  const semanticErrors = workspaceSemanticErrors(config);
  if (semanticErrors.length > 0) throw new Error(`Invalid workspace: ${semanticErrors.join("; ")}`);
  const repositoryBases = /* @__PURE__ */ new Map();
  for (const target of options.request.repositories) {
    const registered = config.repositories[target.name];
    if (!registered) throw new Error(`Unknown repository: ${target.name}`);
    const path2 = assertInside(workspaceRoot20, join8(workspaceRoot20, registered.path));
    await access4(path2);
    await assertCleanRepository(path2);
    repositoryBases.set(target.name, { path: path2, commit: await git(path2, ["rev-parse", registered.default_branch]) });
    await Promise.all([
      access4(join8(workspaceRoot20, "AGENTS.md")),
      access4(join8(workspaceRoot20, "agents", `${registered.agent}.md`)),
      access4(join8(workspaceRoot20, "agents", "repository-worker.md")),
      access4(join8(workspaceRoot20, "agents", "verifier.md"))
    ]);
  }
  const runtimeRoot = assertInside(workspaceRoot20, join8(workspaceRoot20, ".runtime"));
  const now = options.now ?? /* @__PURE__ */ new Date();
  const { workId, runId } = await generateIds(runtimeRoot, options.request.request, now, options.discriminator);
  const createdAt = now.toISOString();
  const branch = `agent/${workId.toLowerCase()}-${slugify2(options.request.request)}-${runId.slice(-8)}`;
  const runRoot = assertInside(runtimeRoot, join8(runtimeRoot, "runs", runId));
  const taskBriefPath = join8(runtimeRoot, "tasks", `${runId}.json`);
  const manifestPath = join8(runRoot, "manifest.json");
  const taskBrief = normalizeContractFirstRequest(options.request, workId, runId, createdAt);
  await assertValid2("task-brief", taskBrief);
  for (const target of taskBrief.repositories) {
    if (target.test_expectation?.policy !== "existing-coverage") continue;
    const base = repositoryBases.get(target.name);
    for (const path2 of target.test_expectation.paths) {
      try {
        await access4(assertInside(base.path, join8(base.path, path2)));
      } catch {
        throw new Error(`Existing-coverage test path does not exist in ${target.name}: ${path2}`);
      }
    }
  }
  await writeJsonAtomic(taskBriefPath, taskBrief);
  const runtimeRepositories = [];
  const preparedRepositories = [];
  for (const target of taskBrief.repositories) {
    const base = repositoryBases.get(target.name);
    const registered = config.repositories[target.name];
    const worktree = assertInside(runtimeRoot, join8(runtimeRoot, "worktrees", runId, target.name));
    const workerInputPath = join8(runRoot, `${target.name}-worker-input.json`);
    const verifierInputPath = join8(runRoot, `${target.name}-verifier-input.json`);
    const workerResultPath = join8(runtimeRoot, "results", `${runId}-${target.name}-worker.json`);
    const verifierResultPath = join8(runtimeRoot, "results", `${runId}-${target.name}-verifier.json`);
    const blockedBy = target.depends_on ?? [];
    const ready = blockedBy.length === 0;
    const instructionPaths = [
      join8(workspaceRoot20, "AGENTS.md"),
      join8(workspaceRoot20, "agents", `${registered.agent}.md`),
      join8(workspaceRoot20, "agents", "repository-worker.md")
    ];
    try {
      await access4(join8(base.path, "AGENTS.md"));
      instructionPaths.push(join8(worktree, "AGENTS.md"));
    } catch {
    }
    const sharedContract = {
      repository: taskBrief.shared_contract.repository,
      paths: taskBrief.shared_contract.paths,
      worktree: assertInside(runtimeRoot, join8(runtimeRoot, "worktrees", runId, taskBrief.shared_contract.repository)),
      approval: target.name === taskBrief.shared_contract.repository ? "must-pass-independent-verification" : "pending"
    };
    await writeJsonAtomic(workerInputPath, {
      contract_version: 1,
      role: "repository-worker",
      task_brief: taskBriefPath,
      repository: target.name,
      worktree,
      branch,
      base_commit: base.commit,
      ready,
      blocked_by: blockedBy,
      shared_contract: sharedContract,
      allowed_scope: target.scope,
      implementation_scope: target.implementation_scope,
      test_expectation: target.test_expectation,
      instruction_paths: instructionPaths,
      result_contract: join8(workspaceRoot20, ".agents", "contracts", "worker-result.schema.json"),
      result_path: workerResultPath
    });
    await writeJsonAtomic(verifierInputPath, {
      contract_version: 1,
      role: "verifier",
      read_only: true,
      task_brief: taskBriefPath,
      repository: target.name,
      worktree,
      branch,
      base_commit: base.commit,
      worker_result: workerResultPath,
      acceptance_criteria: target.acceptance_criteria,
      test_expectation: target.test_expectation,
      verification_commands: target.verification_commands,
      shared_contract: sharedContract,
      instruction_paths: [join8(workspaceRoot20, "AGENTS.md"), join8(workspaceRoot20, "agents", "verifier.md")],
      result_contract: join8(workspaceRoot20, ".agents", "contracts", "verifier-result.schema.json"),
      result_path: verifierResultPath
    });
    runtimeRepositories.push({
      name: target.name,
      base_path: relative5(workspaceRoot20, base.path),
      base_commit: base.commit,
      branch,
      worktree,
      worker_input: workerInputPath,
      verifier_input: verifierInputPath,
      status: ready ? "prepared" : "waiting",
      depends_on: blockedBy,
      repair_attempts: 0
    });
    preparedRepositories.push({ name: target.name, branch, worktree, workerInput: workerInputPath, verifierInput: verifierInputPath, ready, blockedBy });
  }
  const manifest2 = {
    contract_version: 1,
    work_id: workId,
    run_id: runId,
    source_kind: "direct-request",
    status: "preparing",
    created_at: createdAt,
    updated_at: createdAt,
    task_brief: taskBriefPath,
    repositories: runtimeRepositories,
    evidence: [taskBriefPath, manifestPath, ...preparedRepositories.flatMap((repository) => [repository.workerInput, repository.verifierInput])],
    warnings: config.activity.provider === "none" ? ["No activity tool is configured; this run cannot guarantee exclusive ownership."] : [],
    execution_events: [],
    lifecycle_events: config.activity.provider === "none" ? [{ event: "task.starting", status: "skipped", idempotency_key: `${runId}:task.starting:activity-none`, occurred_at: createdAt }] : []
  };
  await assertValid2("runtime-manifest", manifest2);
  await writeJsonAtomic(manifestPath, manifest2);
  if (config.activity.provider !== "none") {
    const lifecycle = await prepareActivityLifecycle({ workspaceRoot: workspaceRoot20, runId, event: "task.starting", availableCapabilities: options.availableCapabilities ?? [], now });
    if (lifecycle.status !== "completed" && lifecycle.status !== "skipped") {
      const primary2 = preparedRepositories.find((repository) => repository.name === taskBrief.shared_contract.repository);
      return { workId, runId, branch: primary2.branch, worktree: primary2.worktree, taskBrief: taskBriefPath, manifest: manifestPath, workerInput: primary2.workerInput, verifierInput: primary2.verifierInput, repositories: preparedRepositories, preparationStatus: lifecycle.status === "failed" ? "blocked" : "awaiting-activity" };
    }
  }
  try {
    await ensurePrivateDirectory(join8(runtimeRoot, "worktrees", runId));
    for (const repository of runtimeRepositories) {
      const base = repositoryBases.get(repository.name);
      await git(base.path, ["worktree", "add", "-b", repository.branch, repository.worktree, repository.base_commit]);
    }
    manifest2.status = "prepared";
    manifest2.updated_at = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    await assertValid2("runtime-manifest", manifest2);
    await writeJsonAtomic(manifestPath, manifest2);
  } catch (error) {
    manifest2.status = "blocked";
    manifest2.updated_at = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    manifest2.evidence.push(`Preparation failed: ${error.message}`);
    await writeJsonAtomic(manifestPath, manifest2);
    throw error;
  }
  const primary = preparedRepositories.find((repository) => repository.name === taskBrief.shared_contract.repository);
  return { workId, runId, branch: primary.branch, worktree: primary.worktree, taskBrief: taskBriefPath, manifest: manifestPath, workerInput: primary.workerInput, verifierInput: primary.verifierInput, repositories: preparedRepositories, preparationStatus: "prepared" };
}
async function preparePlanTask(options) {
  const requestErrors = await validateContract("run-task-request", options.request);
  if (requestErrors.length > 0) {
    throw new Error(`Invalid run-task-request: ${requestErrors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
  }
  const workspaceRoot20 = resolve11(options.workspaceRoot);
  const config = (0, import_yaml7.parse)(await readFile9(join8(workspaceRoot20, "workspace.yaml"), "utf8"));
  await assertValid2("workspace", config);
  const semanticErrors = workspaceSemanticErrors(config);
  if (semanticErrors.length > 0) throw new Error(`Invalid workspace: ${semanticErrors.join("; ")}`);
  const planDirectory3 = assertInside(join8(workspaceRoot20, "context", "plans"), resolve11(workspaceRoot20, options.request.source.reference));
  const validation = await validatePlanDirectory(planDirectory3);
  if (!validation.index || !validation.work_breakdown || validation.errors.length > 0) {
    throw new Error(`Plan validation failed: ${validation.errors.join("; ") || "plan metadata is unavailable"}`);
  }
  const index = validation.index;
  const breakdown = validation.work_breakdown;
  if (index.status !== "approved") throw new Error(`Plan ${index.plan_id} is draft; explicit approval is required`);
  if (options.request.source.plan_version !== index.plan_version) {
    throw new Error(`Plan version is stale: requested ${options.request.source.plan_version}, current ${index.plan_version}`);
  }
  if (options.request.source.approved_digest !== index.approved_digest) {
    throw new Error("Plan approval digest is stale or does not match the approved plan material");
  }
  const items = new Map(breakdown.items.map((item2) => [item2.work_id, item2]));
  const workId = options.request.work_ids[0];
  const item = items.get(workId);
  if (!item) throw new Error(`Unknown plan work ID: ${workId}`);
  if (item.depends_on.length > 0) {
    throw new Error(`${workId} is dependency-blocked by ${item.depends_on.join(", ")}; plan execution currently requires an independently executable item`);
  }
  if (!config.repositories[item.repository]) throw new Error(`Plan work ${workId} repository is not registered: ${item.repository}`);
  const repositoryBases = /* @__PURE__ */ new Map();
  for (const target of [{ name: item.repository }]) {
    const registered = config.repositories[target.name];
    if (!registered) throw new Error(`Unknown repository: ${target.name}`);
    const path2 = assertInside(workspaceRoot20, join8(workspaceRoot20, registered.path));
    await access4(path2);
    await assertCleanRepository(path2);
    repositoryBases.set(target.name, { path: path2, commit: await git(path2, ["rev-parse", registered.default_branch]) });
    await Promise.all([
      access4(join8(workspaceRoot20, "AGENTS.md")),
      access4(join8(workspaceRoot20, "agents", `${registered.agent}.md`)),
      access4(join8(workspaceRoot20, "agents", "repository-worker.md")),
      access4(join8(workspaceRoot20, "agents", "verifier.md"))
    ]);
  }
  const runtimeRoot = assertInside(workspaceRoot20, join8(workspaceRoot20, ".runtime"));
  const now = options.now ?? /* @__PURE__ */ new Date();
  const runId = generateRunId(item.title, now, options.discriminator ?? randomBytes2(4).toString("hex"));
  const createdAt = now.toISOString();
  const branch = `agent/${workId.toLowerCase()}-${slugify2(item.title)}-${runId.slice(-8)}`;
  const runRoot = assertInside(runtimeRoot, join8(runtimeRoot, "runs", runId));
  const taskBriefPath = join8(runtimeRoot, "tasks", `${runId}.json`);
  const manifestPath = join8(runRoot, "manifest.json");
  const taskBrief = normalizePlanRequest(options.request, item, runId, createdAt);
  await assertValid2("task-brief", taskBrief);
  for (const target of taskBrief.repositories) {
    if (target.test_expectation?.policy !== "existing-coverage") continue;
    const base = repositoryBases.get(target.name);
    for (const path2 of target.test_expectation.paths) {
      try {
        await access4(assertInside(base.path, join8(base.path, path2)));
      } catch {
        throw new Error(`Existing-coverage test path does not exist in ${target.name}: ${path2}`);
      }
    }
  }
  await writeJsonAtomic(taskBriefPath, taskBrief);
  const runtimeRepositories = [];
  const preparedRepositories = [];
  for (const target of taskBrief.repositories) {
    const base = repositoryBases.get(target.name);
    const registered = config.repositories[target.name];
    const worktree = assertInside(runtimeRoot, join8(runtimeRoot, "worktrees", runId, target.name));
    const workerInputPath = join8(runRoot, `${target.name}-worker-input.json`);
    const verifierInputPath = join8(runRoot, `${target.name}-verifier-input.json`);
    const workerResultPath = join8(runtimeRoot, "results", `${runId}-${target.name}-worker.json`);
    const verifierResultPath = join8(runtimeRoot, "results", `${runId}-${target.name}-verifier.json`);
    const blockedBy = target.depends_on ?? [];
    const ready = blockedBy.length === 0;
    const instructionPaths = [join8(workspaceRoot20, "AGENTS.md"), join8(workspaceRoot20, "agents", `${registered.agent}.md`), join8(workspaceRoot20, "agents", "repository-worker.md")];
    try {
      await access4(join8(base.path, "AGENTS.md"));
      instructionPaths.push(join8(worktree, "AGENTS.md"));
    } catch {
    }
    await writeJsonAtomic(workerInputPath, {
      contract_version: 1,
      role: "repository-worker",
      task_brief: taskBriefPath,
      repository: target.name,
      worktree,
      branch,
      base_commit: base.commit,
      ready,
      blocked_by: blockedBy,
      allowed_scope: target.scope,
      implementation_scope: target.implementation_scope,
      test_expectation: target.test_expectation,
      instruction_paths: instructionPaths,
      result_contract: join8(workspaceRoot20, ".agents", "contracts", "worker-result.schema.json"),
      result_path: workerResultPath
    });
    await writeJsonAtomic(verifierInputPath, {
      contract_version: 1,
      role: "verifier",
      read_only: true,
      task_brief: taskBriefPath,
      repository: target.name,
      worktree,
      branch,
      base_commit: base.commit,
      worker_result: workerResultPath,
      acceptance_criteria: target.acceptance_criteria,
      test_expectation: target.test_expectation,
      verification_commands: target.verification_commands,
      instruction_paths: [join8(workspaceRoot20, "AGENTS.md"), join8(workspaceRoot20, "agents", "verifier.md")],
      result_contract: join8(workspaceRoot20, ".agents", "contracts", "verifier-result.schema.json"),
      result_path: verifierResultPath
    });
    runtimeRepositories.push({ name: target.name, base_path: relative5(workspaceRoot20, base.path), base_commit: base.commit, branch, worktree, worker_input: workerInputPath, verifier_input: verifierInputPath, status: ready ? "prepared" : "waiting", depends_on: blockedBy, repair_attempts: 0 });
    preparedRepositories.push({ name: target.name, branch, worktree, workerInput: workerInputPath, verifierInput: verifierInputPath, ready, blockedBy });
  }
  const manifest2 = {
    contract_version: 1,
    work_id: workId,
    run_id: runId,
    source_kind: "plan",
    status: "preparing",
    created_at: createdAt,
    updated_at: createdAt,
    task_brief: taskBriefPath,
    repositories: runtimeRepositories,
    plan_work_items: [{ work_id: workId, repository: item.repository, depends_on: item.depends_on, outcome: "pending" }],
    evidence: [taskBriefPath, manifestPath, ...preparedRepositories.flatMap((repository) => [repository.workerInput, repository.verifierInput])],
    warnings: config.activity.provider === "none" ? ["No activity tool is configured; this run cannot guarantee exclusive ownership."] : [],
    execution_events: [],
    lifecycle_events: config.activity.provider === "none" ? [{ event: "task.starting", status: "skipped", idempotency_key: `${runId}:task.starting:activity-none`, occurred_at: createdAt }] : []
  };
  await assertValid2("runtime-manifest", manifest2);
  await writeJsonAtomic(manifestPath, manifest2);
  if (config.activity.provider !== "none") {
    const lifecycle = await prepareActivityLifecycle({ workspaceRoot: workspaceRoot20, runId, event: "task.starting", availableCapabilities: options.availableCapabilities ?? [], now });
    if (lifecycle.status !== "completed" && lifecycle.status !== "skipped") {
      const primary2 = preparedRepositories[0];
      return { workId, runId, branch, worktree: primary2.worktree, taskBrief: taskBriefPath, manifest: manifestPath, workerInput: primary2.workerInput, verifierInput: primary2.verifierInput, repositories: preparedRepositories, preparationStatus: lifecycle.status === "failed" ? "blocked" : "awaiting-activity" };
    }
  }
  try {
    await ensurePrivateDirectory(join8(runtimeRoot, "worktrees", runId));
    for (const repository of runtimeRepositories) {
      await git(repositoryBases.get(repository.name).path, ["worktree", "add", "-b", repository.branch, repository.worktree, repository.base_commit]);
    }
    manifest2.status = "prepared";
    manifest2.updated_at = now.toISOString();
    await assertValid2("runtime-manifest", manifest2);
    await writeJsonAtomic(manifestPath, manifest2);
  } catch (error) {
    manifest2.status = "blocked";
    manifest2.updated_at = (/* @__PURE__ */ new Date()).toISOString();
    manifest2.evidence.push(`Preparation failed: ${error.message}`);
    await writeJsonAtomic(manifestPath, manifest2);
    throw error;
  }
  const primary = preparedRepositories[0];
  return { workId, runId, branch, worktree: primary.worktree, taskBrief: taskBriefPath, manifest: manifestPath, workerInput: primary.workerInput, verifierInput: primary.verifierInput, repositories: preparedRepositories, preparationStatus: "prepared" };
}
async function resumePlanlessTask(options) {
  const workspaceRoot20 = resolve11(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot20, join8(workspaceRoot20, ".runtime"));
  const manifestPath = assertInside(runtimeRoot, join8(runtimeRoot, "runs", options.runId, "manifest.json"));
  const manifest2 = JSON.parse(await readFile9(manifestPath, "utf8"));
  await assertValid2("runtime-manifest", manifest2);
  if (manifest2.run_id !== options.runId || manifest2.status !== "preparing") throw new Error(`Run ${options.runId} is not awaiting preparation`);
  const lifecycle = manifest2.lifecycle_events.find((event) => event.event === "task.starting");
  if (!lifecycle || lifecycle.status !== "completed" && lifecycle.status !== "skipped") {
    throw new Error(`task.starting lifecycle is ${lifecycle?.status ?? "missing"}; complete required or manual actions before resuming`);
  }
  const config = (0, import_yaml7.parse)(await readFile9(join8(workspaceRoot20, "workspace.yaml"), "utf8"));
  await assertValid2("workspace", config);
  const basePaths = /* @__PURE__ */ new Map();
  for (const repository of manifest2.repositories) {
    const registered = config.repositories[repository.name];
    if (!registered) throw new Error(`Unknown repository: ${repository.name}`);
    const repositoryPath = assertInside(workspaceRoot20, join8(workspaceRoot20, repository.base_path));
    await assertCleanRepository(repositoryPath);
    const currentBase = await git(repositoryPath, ["rev-parse", registered.default_branch]);
    if (currentBase !== repository.base_commit) throw new Error(`Repository base changed during activity preflight for ${repository.name}; prepare a fresh run`);
    basePaths.set(repository.name, repositoryPath);
  }
  try {
    await ensurePrivateDirectory(join8(runtimeRoot, "worktrees", options.runId));
    for (const repository of manifest2.repositories) {
      await git(basePaths.get(repository.name), ["worktree", "add", "-b", repository.branch, repository.worktree, repository.base_commit]);
    }
    manifest2.status = "prepared";
    manifest2.updated_at = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    await assertValid2("runtime-manifest", manifest2);
    await writeJsonAtomic(manifestPath, manifest2);
  } catch (error) {
    manifest2.status = "blocked";
    manifest2.updated_at = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    manifest2.evidence.push(`Preparation failed: ${error.message}`);
    await writeJsonAtomic(manifestPath, manifest2);
    throw error;
  }
  const brief = JSON.parse(await readFile9(manifest2.task_brief, "utf8"));
  const primary = manifest2.repositories.find((repository) => repository.name === brief.shared_contract?.repository) ?? manifest2.repositories[0];
  const repositories = manifest2.repositories.map((repository) => ({
    name: repository.name,
    branch: repository.branch,
    worktree: repository.worktree,
    workerInput: repository.worker_input,
    verifierInput: repository.verifier_input,
    ready: (repository.depends_on?.length ?? 0) === 0,
    blockedBy: repository.depends_on ?? []
  }));
  return {
    workId: manifest2.work_id,
    runId: manifest2.run_id,
    branch: primary.branch,
    worktree: primary.worktree,
    taskBrief: manifest2.task_brief,
    manifest: manifestPath,
    workerInput: primary.worker_input,
    verifierInput: primary.verifier_input,
    repositories,
    preparationStatus: "prepared"
  };
}
var import_yaml7;
var init_run_task = __esm({
  "scripts/lib/run-task.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    import_yaml7 = __toESM(require_dist(), 1);
    init_git();
    init_ids();
    init_io();
    init_validation();
    init_activity_lifecycle();
    init_plans();
  }
});

// scripts/run-task.ts
var run_task_exports = {};
import { dirname as dirname9, resolve as resolve12 } from "node:path";
import { readFile as readFile10 } from "node:fs/promises";
import { parseArgs as parseArgs4 } from "node:util";
import { fileURLToPath as fileURLToPath5 } from "node:url";
var testPolicies, activityCapabilities, workspaceRoot4, values4, prepared;
var init_run_task2 = __esm({
  async "scripts/run-task.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_run_task();
    testPolicies = ["required", "existing-coverage", "verifier-only", "not-required"];
    activityCapabilities = ["read-tasks", "update-status", "create-tasks", "assign-task", "timers"];
    workspaceRoot4 = resolve12(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve12(dirname9(fileURLToPath5(import.meta.url)), ".."));
    ({ values: values4 } = parseArgs4({
      options: {
        request: { type: "string" },
        repository: { type: "string" },
        acceptance: { type: "string", multiple: true, default: [] },
        scope: { type: "string", multiple: true, default: [] },
        "test-scope": { type: "string", multiple: true, default: [] },
        "test-policy": { type: "string" },
        "test-rationale": { type: "string" },
        verify: { type: "string", multiple: true, default: [] },
        available: { type: "string", multiple: true, default: [] },
        "resume-run": { type: "string" },
        "request-file": { type: "string" }
      }
    }));
    if (values4["resume-run"]) {
      console.log(JSON.stringify(await resumePlanlessTask({ workspaceRoot: workspaceRoot4, runId: values4["resume-run"] }), null, 2));
      process.exit(0);
    }
    if (values4.available.some((capability) => !activityCapabilities.includes(capability))) {
      throw new Error(`Unknown activity capability; expected one of: ${activityCapabilities.join(", ")}`);
    }
    if (values4["request-file"]) {
      const request4 = JSON.parse(await readFile10(resolve12(values4["request-file"]), "utf8"));
      const prepared2 = "source" in request4 && request4.source?.kind === "plan" ? await preparePlanTask({ workspaceRoot: workspaceRoot4, request: request4, availableCapabilities: values4.available }) : await prepareContractFirstTask({ workspaceRoot: workspaceRoot4, request: request4, availableCapabilities: values4.available });
      console.log(JSON.stringify(prepared2, null, 2));
      process.exit(0);
    }
    if (!values4.request || !values4.repository || values4["test-policy"] && !testPolicies.includes(values4["test-policy"])) {
      throw new Error("Usage: run-task --request <text> --repository <name> --acceptance <criterion> --scope <path> [...] | run-task --request-file <json>");
    }
    prepared = await preparePlanlessTask({
      workspaceRoot: workspaceRoot4,
      request: values4.request,
      repository: values4.repository,
      acceptanceCriteria: values4.acceptance,
      scope: values4.scope,
      testScope: values4["test-scope"],
      ...values4["test-policy"] ? { testPolicy: values4["test-policy"] } : {},
      ...values4["test-rationale"] ? { testRationale: values4["test-rationale"] } : {},
      verificationCommands: values4.verify,
      availableCapabilities: values4.available
    });
    console.log(JSON.stringify(prepared, null, 2));
    if (prepared.preparationStatus === "prepared") console.warn("Warning: exclusive ownership is guaranteed only when a configured starting action confirmed it.");
    else console.warn(`Worktree not created: activity preflight is ${prepared.preparationStatus}. Complete the recorded actions, then rerun with --resume-run ${prepared.runId}.`);
  }
});

// scripts/lib/record-result.ts
import { chmod as chmod2, readFile as readFile11 } from "node:fs/promises";
import { join as join9, resolve as resolve13 } from "node:path";
function assertIdentifier(value2, label, pattern) {
  if (!pattern.test(value2)) throw new Error(`Invalid ${label}: ${value2}`);
}
async function readJson2(path2) {
  return JSON.parse(await readFile11(path2, "utf8"));
}
async function assertValid3(name, value2) {
  const errors2 = await validateContract(name, value2);
  if (errors2.length > 0) {
    throw new Error(`Invalid ${name}: ${errors2.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
  }
}
function assertEqual(actual, expected, label) {
  if (actual !== expected) throw new Error(`${label} mismatch: expected ${String(expected)}, received ${String(actual)}`);
}
function inAllowedScope(path2, scopes) {
  return scopes.some((scope) => {
    const normalized = scope.replace(/\/$/, "");
    return path2 === normalized || path2.startsWith(`${normalized}/`);
  });
}
function sameMembers(left, right) {
  return [...left].sort().join("\n") === [...right].sort().join("\n");
}
function sameTestExpectation(left, right) {
  return left.policy === right.policy && left.rationale === right.rationale && sameMembers(left.paths, right.paths);
}
function findRepository(manifest2, name) {
  const repository = manifest2.repositories.find((candidate) => candidate.name === name);
  if (!repository) throw new Error(`Run ${manifest2.run_id} has no repository named ${name}`);
  return repository;
}
function assertPlanWorkItemAssociation(manifest2, brief, repository) {
  assertEqual(manifest2.source_kind, brief.source.kind, "manifest source_kind");
  const planLinked = brief.source.kind === "plan";
  if (!planLinked) {
    if (manifest2.plan_work_items !== void 0) throw new Error("Non-plan run must not contain plan work items");
    return;
  }
  if (brief.plan.approval_state !== "approved") throw new Error("Plan-linked task brief must contain approved plan metadata");
  if (!manifest2.plan_work_items || manifest2.plan_work_items.length !== 1) {
    throw new Error("Plan-linked run must contain exactly one plan work item");
  }
  const item = manifest2.plan_work_items[0];
  assertEqual(item.work_id, manifest2.work_id, "plan work item work_id");
  if (brief.plan.work_ids.length !== 1 || brief.plan.work_ids[0] !== item.work_id) {
    throw new Error("Plan work item identity does not match task brief work IDs");
  }
  if (manifest2.repositories.length !== 1 || manifest2.repositories[0].name !== item.repository || brief.repositories.length !== 1 || brief.repositories[0].name !== item.repository || item.repository !== repository) {
    throw new Error("Plan work item repository does not match task brief and recorded repository");
  }
}
function eventKey(runId, repository, stage, attempt) {
  const suffix = attempt === 0 ? "" : `:attempt-${attempt}`;
  return `${runId}:execution:${repository}:${stage}${suffix}`;
}
function appendEvent(manifest2, stage, repository, from, to, occurredAt, inferred, attempt, resultPath) {
  const event = {
    stage,
    repository,
    from_status: from,
    to_status: to,
    inferred,
    attempt,
    idempotency_key: eventKey(manifest2.run_id, repository, stage, attempt),
    occurred_at: occurredAt
  };
  if (resultPath) event.result_path = resultPath;
  manifest2.execution_events ??= [];
  manifest2.execution_events.push(event);
  const runtimeRepository = findRepository(manifest2, repository);
  runtimeRepository.status = to;
  refreshManifestStatus(manifest2);
  manifest2.updated_at = occurredAt;
  if (resultPath && !manifest2.evidence.includes(resultPath)) manifest2.evidence.push(resultPath);
}
function refreshManifestStatus(manifest2) {
  const statuses = manifest2.repositories.map((repository) => repository.status ?? manifest2.status);
  if (statuses.every((status2) => status2 === "passed")) manifest2.status = "passed";
  else if (statuses.includes("failed")) manifest2.status = "failed";
  else if (statuses.includes("blocked")) manifest2.status = "blocked";
  else if (statuses.includes("verifying")) manifest2.status = "verifying";
  else if (statuses.includes("running")) manifest2.status = "running";
  else manifest2.status = "prepared";
}
async function unlockDependents(runtimeRoot, manifest2) {
  for (const candidate of manifest2.repositories) {
    if (candidate.status !== "waiting") continue;
    const dependencies = candidate.depends_on ?? [];
    if (!dependencies.every((name) => findRepository(manifest2, name).status === "passed")) continue;
    for (const [path2, worker] of [[candidate.worker_input, true], [candidate.verifier_input, false]]) {
      const inputPath2 = assertInside(runtimeRoot, path2);
      const input = await readJson2(inputPath2);
      if (worker) {
        input.ready = true;
        input.blocked_by = [];
      }
      if (input.shared_contract && typeof input.shared_contract === "object") {
        input.shared_contract.approval = "verified";
      }
      await writeJsonAtomic(inputPath2, input);
    }
    candidate.status = "prepared";
  }
  refreshManifestStatus(manifest2);
}
async function assertWorktree(repository) {
  await assertCleanRepository(repository.worktree);
  const branch = await git(repository.worktree, ["branch", "--show-current"]);
  assertEqual(branch, repository.branch, "worktree branch");
  return git(repository.worktree, ["rev-parse", "HEAD"]);
}
function assertTaskIdentity(manifest2, brief, repository) {
  assertEqual(brief.work_id, manifest2.work_id, "task brief work_id");
  assertEqual(brief.run_id, manifest2.run_id, "task brief run_id");
  if (!brief.repositories.some((candidate) => candidate.name === repository)) {
    throw new Error(`Task brief does not include repository ${repository}`);
  }
}
async function validateWorkerResult(manifest2, repository, input, testExpectation) {
  const result3 = await readJson2(input.result_path);
  await assertValid3("worker-result", result3);
  assertEqual(result3.work_id, manifest2.work_id, "worker result work_id");
  assertEqual(result3.run_id, manifest2.run_id, "worker result run_id");
  assertEqual(result3.repository, repository.name, "worker result repository");
  assertEqual(result3.branch, repository.branch, "worker result branch");
  assertEqual(resolve13(result3.worktree), resolve13(repository.worktree), "worker result worktree");
  if (result3.status === "completed") {
    if (result3.commits.length === 0) throw new Error("Completed worker result must record at least one commit");
    if (result3.checks.some((check) => check.status === "failed")) throw new Error("Completed worker result cannot contain a failed check");
    const head = await assertWorktree(repository);
    assertEqual(result3.commits.at(-1), head, "worker result final commit");
    const commits = (await git(repository.worktree, ["rev-list", "--reverse", `${repository.base_commit}..${head}`])).split("\n").filter(Boolean);
    if (commits.join("\n") !== result3.commits.join("\n")) throw new Error("worker result commits does not match the ordered base-to-head Git history");
    const changedFiles = (await git(repository.worktree, ["diff", "--name-only", `${repository.base_commit}...${head}`])).split("\n").filter(Boolean);
    if (!sameMembers(changedFiles, result3.changed_files)) throw new Error("worker result changed_files does not match the base-to-head Git diff");
    const outsideScope = changedFiles.filter((path2) => !inAllowedScope(path2, input.allowed_scope));
    if (outsideScope.length > 0) throw new Error(`Worker changed files outside allowed scope: ${outsideScope.join(", ")}`);
    if (testExpectation.policy === "required" && !changedFiles.some((path2) => inAllowedScope(path2, testExpectation.paths))) {
      throw new Error(`Required test policy needs a changed file in test scope: ${testExpectation.paths.join(", ")}`);
    }
    return "verifying";
  }
  return result3.status;
}
async function validateVerifierResult(manifest2, repository, input) {
  const worker = await readJson2(input.worker_result);
  await assertValid3("worker-result", worker);
  const result3 = await readJson2(input.result_path);
  await assertValid3("verifier-result", result3);
  assertEqual(result3.work_id, manifest2.work_id, "verifier result work_id");
  assertEqual(result3.run_id, manifest2.run_id, "verifier result run_id");
  assertEqual(result3.repository, repository.name, "verifier result repository");
  if (!sameMembers(result3.acceptance.map((item) => item.criterion), input.acceptance_criteria)) {
    throw new Error("verifier result acceptance criteria do not match verifier input");
  }
  const acceptanceStatuses = result3.acceptance.map((item) => item.status);
  if (result3.status === "pass" && acceptanceStatuses.some((status2) => status2 !== "passed")) {
    throw new Error("Passing verifier result requires every acceptance criterion to pass");
  }
  if (result3.status === "fail" && !acceptanceStatuses.includes("failed")) {
    throw new Error("Failing verifier result must identify a failed acceptance criterion");
  }
  if (result3.status === "blocked" && !acceptanceStatuses.includes("blocked")) {
    throw new Error("Blocked verifier result must identify a blocked acceptance criterion");
  }
  const head = await assertWorktree(repository);
  assertEqual(worker.commits.at(-1), head, "verified worker commit");
  return result3.status === "pass" ? "passed" : result3.status === "fail" ? "failed" : "blocked";
}
async function recordResult(options) {
  assertIdentifier(options.runId, "run ID", /^[0-9]{8}T[0-9]{6}Z-[a-f0-9]{8}$/);
  assertIdentifier(options.repository, "repository", /^[a-z][a-z0-9-]*$/);
  const workspaceRoot20 = resolve13(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot20, join9(workspaceRoot20, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const manifestPath = assertInside(runtimeRoot, join9(runtimeRoot, "runs", options.runId, "manifest.json"));
  const lockPath = `${manifestPath}.lock`;
  return withExclusiveFile(lockPath, async () => {
    const manifest2 = await readJson2(manifestPath);
    await assertValid3("runtime-manifest", manifest2);
    assertEqual(manifest2.run_id, options.runId, "manifest run_id");
    const repository = findRepository(manifest2, options.repository);
    const attempt = repository.repair_attempts ?? 0;
    assertInside(runtimeRoot, repository.worktree);
    const taskBriefPath = assertInside(runtimeRoot, manifest2.task_brief);
    const workerInputPath = assertInside(runtimeRoot, repository.worker_input);
    const verifierInputPath = assertInside(runtimeRoot, repository.verifier_input);
    const brief = await readJson2(taskBriefPath);
    await assertValid3("task-brief", brief);
    assertPlanWorkItemAssociation(manifest2, brief, options.repository);
    assertTaskIdentity(manifest2, brief, options.repository);
    const target = brief.repositories.find((candidate) => candidate.name === options.repository);
    const targetScope = target.scope ?? brief.scope;
    const implementationScope = target.implementation_scope ?? brief.implementation_scope ?? targetScope;
    const testExpectation = target.test_expectation ?? brief.test_expectation ?? {
      policy: "verifier-only",
      paths: [],
      rationale: "Legacy task brief has no authorized test edit scope; verifier evidence is required."
    };
    const workerInput = await readJson2(workerInputPath);
    const verifierInput = await readJson2(verifierInputPath);
    assertInside(runtimeRoot, workerInput.result_path);
    assertInside(runtimeRoot, verifierInput.worker_result);
    assertInside(runtimeRoot, verifierInput.result_path);
    assertEqual(workerInput.repository, repository.name, "worker input repository");
    assertEqual(resolve13(workerInput.task_brief), resolve13(taskBriefPath), "worker input task_brief");
    assertEqual(resolve13(workerInput.worktree), resolve13(repository.worktree), "worker input worktree");
    assertEqual(workerInput.branch, repository.branch, "worker input branch");
    assertEqual(workerInput.base_commit, repository.base_commit, "worker input base_commit");
    if (!sameMembers(workerInput.allowed_scope, targetScope)) throw new Error("worker input allowed_scope does not match task brief repository scope");
    if ((target.implementation_scope || brief.implementation_scope) && (!workerInput.implementation_scope || !sameMembers(workerInput.implementation_scope, implementationScope))) {
      throw new Error("worker input implementation_scope does not match task brief");
    }
    if ((target.test_expectation || brief.test_expectation) && (!workerInput.test_expectation || !sameTestExpectation(workerInput.test_expectation, testExpectation))) {
      throw new Error("worker input test_expectation does not match task brief");
    }
    assertEqual(verifierInput.repository, repository.name, "verifier input repository");
    assertEqual(resolve13(verifierInput.task_brief), resolve13(taskBriefPath), "verifier input task_brief");
    assertEqual(resolve13(verifierInput.worktree), resolve13(repository.worktree), "verifier input worktree");
    assertEqual(verifierInput.branch, repository.branch, "verifier input branch");
    assertEqual(verifierInput.base_commit, repository.base_commit, "verifier input base_commit");
    assertEqual(resolve13(verifierInput.worker_result), resolve13(workerInput.result_path), "verifier input worker_result");
    if (!sameMembers(verifierInput.acceptance_criteria, target.acceptance_criteria ?? brief.acceptance_criteria)) throw new Error("verifier input acceptance_criteria does not match task brief repository criteria");
    if ((target.test_expectation || brief.test_expectation) && (!verifierInput.test_expectation || !sameTestExpectation(verifierInput.test_expectation, testExpectation))) {
      throw new Error("verifier input test_expectation does not match task brief");
    }
    const existing = manifest2.execution_events?.find((event) => event.idempotency_key === eventKey(options.runId, options.repository, options.stage, attempt));
    const occurredAt = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    const currentStatus = repository.status ?? manifest2.status;
    if (options.stage === "worker-started") {
      if (existing) return manifest2;
      if (currentStatus === "waiting") throw new Error(`worker-started for ${repository.name} is blocked by: ${(repository.depends_on ?? []).join(", ")}`);
      if (currentStatus !== "prepared") throw new Error(`worker-started requires prepared repository status, received ${currentStatus}`);
      const head = await assertWorktree(repository);
      assertEqual(head, repository.base_commit, "worker start HEAD");
      appendEvent(manifest2, options.stage, options.repository, "prepared", "running", occurredAt, false, attempt);
    } else if (options.stage === "worker-result") {
      const target2 = await validateWorkerResult(manifest2, repository, workerInput, testExpectation);
      await chmod2(workerInput.result_path, 384);
      if (existing) return manifest2;
      if (currentStatus === "prepared") {
        appendEvent(manifest2, "worker-started", options.repository, "prepared", "running", occurredAt, true, attempt);
      }
      if ((repository.status ?? manifest2.status) !== "running") throw new Error(`worker-result requires running repository status, received ${repository.status ?? manifest2.status}`);
      appendEvent(manifest2, options.stage, options.repository, "running", target2, occurredAt, false, attempt, workerInput.result_path);
    } else {
      const target2 = await validateVerifierResult(manifest2, repository, verifierInput);
      await chmod2(verifierInput.result_path, 384);
      if (existing) return manifest2;
      if (currentStatus !== "verifying") throw new Error(`verifier-result requires verifying repository status, received ${currentStatus}`);
      appendEvent(manifest2, options.stage, options.repository, "verifying", target2, occurredAt, false, attempt, verifierInput.result_path);
      const item = manifest2.plan_work_items?.find((candidate) => candidate.work_id === manifest2.work_id);
      if (manifest2.plan_work_items && (!item || item.repository !== options.repository)) {
        throw new Error("Plan work item identity does not match the verified manifest work and repository");
      }
      if (item) item.outcome = target2;
      if (target2 === "passed") await unlockDependents(runtimeRoot, manifest2);
    }
    await assertValid3("runtime-manifest", manifest2);
    await writeJsonAtomic(manifestPath, manifest2);
    return manifest2;
  });
}
var init_record_result = __esm({
  "scripts/lib/record-result.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_git();
    init_io();
    init_validation();
  }
});

// scripts/record-result.ts
var record_result_exports = {};
import { dirname as dirname10, resolve as resolve14 } from "node:path";
import { parseArgs as parseArgs5 } from "node:util";
import { fileURLToPath as fileURLToPath6 } from "node:url";
var stages, workspaceRoot5, values5, manifest;
var init_record_result2 = __esm({
  async "scripts/record-result.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_record_result();
    stages = ["worker-started", "worker-result", "verifier-result"];
    workspaceRoot5 = resolve14(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve14(dirname10(fileURLToPath6(import.meta.url)), ".."));
    ({ values: values5 } = parseArgs5({
      options: {
        "run-id": { type: "string" },
        repository: { type: "string" },
        stage: { type: "string" }
      }
    }));
    if (!values5["run-id"] || !values5.repository || !values5.stage || !stages.includes(values5.stage)) {
      throw new Error("Usage: record-result --run-id <id> --repository <name> --stage <worker-started|worker-result|verifier-result>");
    }
    manifest = await recordResult({
      workspaceRoot: workspaceRoot5,
      runId: values5["run-id"],
      repository: values5.repository,
      stage: values5.stage
    });
    console.log(JSON.stringify({ runId: manifest.run_id, status: manifest.status, executionEvents: manifest.execution_events?.length ?? 0 }, null, 2));
  }
});

// scripts/lib/review-lifecycle.ts
import { readFile as readFile12 } from "node:fs/promises";
import { join as join10, resolve as resolve15 } from "node:path";
async function assertValid4(name, value2) {
  const errors2 = await validateContract(name, value2);
  if (errors2.length > 0) throw new Error(`Invalid ${name}: ${errors2.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
}
function findRepository2(manifest2, name) {
  const repository = manifest2.repositories.find((candidate) => candidate.name === name);
  if (!repository) throw new Error(`Run ${manifest2.run_id} has no repository named ${name}`);
  return repository;
}
function addEvidence(manifest2, ...paths2) {
  for (const path2 of paths2) if (!manifest2.evidence.includes(path2)) manifest2.evidence.push(path2);
}
function addExecutionEvent(manifest2, repository, stage, fromStatus, toStatus, attempt, occurredAt, resultPath) {
  const event = {
    stage,
    repository,
    from_status: fromStatus,
    to_status: toStatus,
    inferred: false,
    attempt,
    idempotency_key: `${manifest2.run_id}:execution:${repository}:${stage}:attempt-${attempt}`,
    occurred_at: occurredAt
  };
  if (resultPath) event.result_path = resultPath;
  manifest2.execution_events ??= [];
  manifest2.execution_events.push(event);
  const runtimeRepository = findRepository2(manifest2, repository);
  runtimeRepository.status = toStatus;
  const statuses = manifest2.repositories.map((candidate) => candidate.status ?? manifest2.status);
  if (statuses.every((status2) => status2 === "passed")) manifest2.status = "passed";
  else if (statuses.includes("failed")) manifest2.status = "failed";
  else if (statuses.includes("blocked")) manifest2.status = "blocked";
  else if (statuses.includes("verifying")) manifest2.status = "verifying";
  else if (statuses.includes("running")) manifest2.status = "running";
  else manifest2.status = "prepared";
  manifest2.updated_at = occurredAt;
}
async function loadWorkspace(workspaceRoot20) {
  const config = (0, import_yaml8.parse)(await readFile12(join10(workspaceRoot20, "workspace.yaml"), "utf8"));
  await assertValid4("workspace", config);
  const semanticErrors = workspaceSemanticErrors(config);
  if (semanticErrors.length > 0) throw new Error(`Invalid workspace: ${semanticErrors.join("; ")}`);
  return config;
}
function sameMembers2(left, right) {
  return left.slice().sort().join("\n") === right.slice().sort().join("\n");
}
async function assertCurrentWorker(runtimeRoot, manifest2, repository, workerInput) {
  const worker = await readJsonRegularInside(runtimeRoot, workerInput.result_path, "Worker result");
  await assertValid4("worker-result", worker);
  if (worker.work_id !== manifest2.work_id || worker.run_id !== manifest2.run_id || worker.repository !== repository.name) {
    throw new Error("Worker result identity does not match the active run");
  }
  if (worker.status !== "completed") throw new Error(`Review lifecycle requires a completed worker, received ${worker.status}`);
  if (worker.branch !== repository.branch || resolve15(worker.worktree) !== resolve15(repository.worktree)) {
    throw new Error("Worker result branch or worktree does not match the runtime manifest");
  }
  await assertCleanRepository(repository.worktree);
  const branch = await git(repository.worktree, ["branch", "--show-current"]);
  if (branch !== repository.branch) throw new Error(`Worktree branch mismatch: expected ${repository.branch}, received ${branch}`);
  const head = await git(repository.worktree, ["rev-parse", "HEAD"]);
  if (worker.commits.at(-1) !== head) throw new Error("Current worktree HEAD does not match the recorded worker result");
  const commits = (await git(repository.worktree, ["rev-list", "--reverse", `${repository.base_commit}..${head}`])).split("\n").filter(Boolean);
  const changedFiles = (await git(repository.worktree, ["diff", "--name-only", `${repository.base_commit}...${head}`])).split("\n").filter(Boolean);
  if (commits.join("\n") !== worker.commits.join("\n") || !sameMembers2(changedFiles, worker.changed_files)) {
    throw new Error("Worker result no longer matches the current base-to-head Git history");
  }
  return worker;
}
function assertVerifier(manifest2, repository, brief, verifier) {
  if (verifier.work_id !== manifest2.work_id || verifier.run_id !== manifest2.run_id || verifier.repository !== repository.name) {
    throw new Error("Verifier result identity does not match the active run");
  }
  const acceptanceCriteria = brief.repositories.find((candidate) => candidate.name === repository.name)?.acceptance_criteria ?? brief.acceptance_criteria;
  if (!sameMembers2(verifier.acceptance.map((item) => item.criterion), acceptanceCriteria)) {
    throw new Error("Verifier acceptance criteria do not match the task brief");
  }
  if (verifier.status === "pass" && verifier.acceptance.some((item) => item.status !== "passed")) {
    throw new Error("Passing verifier result no longer has complete passing acceptance evidence");
  }
  if (verifier.status === "fail" && !verifier.acceptance.some((item) => item.status === "failed")) {
    throw new Error("Failing verifier result no longer identifies failed acceptance evidence");
  }
}
function repairFindings(verifier) {
  const findings = verifier.findings.map((finding) => `[${finding.severity}] ${finding.description} Evidence: ${finding.evidence}`);
  for (const acceptance of verifier.acceptance) {
    if (acceptance.status !== "passed") findings.push(`Acceptance ${acceptance.status}: ${acceptance.criterion}. Evidence: ${acceptance.evidence}`);
  }
  return findings;
}
async function prepareRepair(options) {
  const workspaceRoot20 = resolve15(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot20, join10(workspaceRoot20, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const manifestPath = assertInside(runtimeRoot, join10(runtimeRoot, "runs", options.runId, "manifest.json"));
  const lockPath = `${manifestPath}.lock`;
  const config = await loadWorkspace(workspaceRoot20);
  return withExclusiveFile(lockPath, async () => {
    const manifest2 = await readJsonRegularInside(runtimeRoot, manifestPath, "Runtime manifest");
    await assertValid4("runtime-manifest", manifest2);
    if (manifest2.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository2(manifest2, options.repository);
    const attempt = repository.repair_attempts ?? 0;
    const maximumAttempts = config.workflow.maximum_repair_attempts;
    const lastEvent = manifest2.execution_events?.filter((event) => event.repository === repository.name).at(-1);
    const repositoryStatus = repository.status ?? manifest2.status;
    if (repositoryStatus === "running" && lastEvent?.stage === "repair-prepared" && lastEvent.attempt === attempt) {
      return { status: "prepared", attempt, maximum_attempts: maximumAttempts, manifest: manifestPath, worker_input: repository.worker_input, verifier_input: repository.verifier_input };
    }
    if (repositoryStatus === "blocked" && lastEvent?.stage === "repair-exhausted") {
      return { status: "exhausted", attempt, maximum_attempts: maximumAttempts, manifest: manifestPath };
    }
    if (repositoryStatus !== "failed") throw new Error(`Repair preparation requires failed status, received ${repositoryStatus}`);
    const occurredAt = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    if (attempt >= maximumAttempts) {
      addExecutionEvent(manifest2, repository.name, "repair-exhausted", "failed", "blocked", attempt, occurredAt);
      manifest2.warnings.push(`Maximum repair attempts exhausted for ${repository.name}: ${maximumAttempts}`);
      await assertValid4("runtime-manifest", manifest2);
      await writeJsonAtomic(manifestPath, manifest2);
      return { status: "exhausted", attempt, maximum_attempts: maximumAttempts, manifest: manifestPath };
    }
    const taskBriefPath = assertInside(runtimeRoot, manifest2.task_brief);
    const taskBrief = await readJsonRegularInside(runtimeRoot, taskBriefPath, "Task brief");
    await assertValid4("task-brief", taskBrief);
    const taskTarget = taskBrief.repositories.find((candidate) => candidate.name === repository.name);
    if (!taskTarget) throw new Error(`Task brief does not include repository ${repository.name}`);
    const priorWorkerInput = await readJsonRegularInside(runtimeRoot, repository.worker_input, "Worker input");
    const priorVerifierInput = await readJsonRegularInside(runtimeRoot, repository.verifier_input, "Verifier input");
    assertInside(runtimeRoot, priorWorkerInput.result_path);
    assertInside(runtimeRoot, priorVerifierInput.result_path);
    await assertCurrentWorker(runtimeRoot, manifest2, repository, priorWorkerInput);
    const verifier = await readJsonRegularInside(runtimeRoot, priorVerifierInput.result_path, "Verifier result");
    await assertValid4("verifier-result", verifier);
    assertVerifier(manifest2, repository, taskBrief, verifier);
    if (verifier.status !== "fail") throw new Error(`Repair preparation requires a failing verifier result, received ${verifier.status}`);
    const findings = repairFindings(verifier);
    if (findings.length === 0) throw new Error("Failing verifier result contains no actionable findings");
    const nextAttempt = attempt + 1;
    const runRoot = join10(runtimeRoot, "runs", options.runId);
    const workerInputPath = join10(runRoot, `${repository.name}-repair-${nextAttempt}-worker-input.json`);
    const verifierInputPath = join10(runRoot, `${repository.name}-repair-${nextAttempt}-verifier-input.json`);
    const workerResultPath = join10(runtimeRoot, "results", `${options.runId}-${repository.name}-repair-${nextAttempt}-worker.json`);
    const verifierResultPath = join10(runtimeRoot, "results", `${options.runId}-${repository.name}-repair-${nextAttempt}-verifier.json`);
    const repositoryConfig = config.repositories[repository.name];
    if (!repositoryConfig) throw new Error(`Workspace has no repository named ${repository.name}`);
    const instructionPaths = [
      join10(workspaceRoot20, "AGENTS.md"),
      join10(workspaceRoot20, "agents", "repository-worker.md"),
      join10(workspaceRoot20, "agents", `${repositoryConfig.agent}.md`)
    ];
    const workerInput = {
      contract_version: 1,
      role: "repair-worker",
      attempt: nextAttempt,
      task_brief: taskBriefPath,
      repository: repository.name,
      worktree: repository.worktree,
      branch: repository.branch,
      base_commit: repository.base_commit,
      allowed_scope: taskTarget.scope ?? taskBrief.scope,
      implementation_scope: taskTarget.implementation_scope ?? taskBrief.implementation_scope,
      test_expectation: taskTarget.test_expectation ?? taskBrief.test_expectation,
      findings,
      previous_worker_result: priorWorkerInput.result_path,
      previous_verifier_result: priorVerifierInput.result_path,
      instruction_paths: instructionPaths,
      result_contract: join10(workspaceRoot20, ".agents", "contracts", "worker-result.schema.json"),
      result_path: workerResultPath
    };
    const verifierInput = {
      contract_version: 1,
      role: "verifier",
      read_only: true,
      attempt: nextAttempt,
      task_brief: taskBriefPath,
      repository: repository.name,
      worktree: repository.worktree,
      branch: repository.branch,
      base_commit: repository.base_commit,
      worker_result: workerResultPath,
      acceptance_criteria: taskTarget.acceptance_criteria ?? taskBrief.acceptance_criteria,
      test_expectation: taskTarget.test_expectation ?? taskBrief.test_expectation,
      verification_commands: taskTarget.verification_commands ?? taskBrief.verification_commands,
      instruction_paths: [join10(workspaceRoot20, "AGENTS.md"), join10(workspaceRoot20, "agents", "verifier.md")],
      result_contract: join10(workspaceRoot20, ".agents", "contracts", "verifier-result.schema.json"),
      result_path: verifierResultPath
    };
    await writeJsonAtomic(workerInputPath, workerInput);
    await writeJsonAtomic(verifierInputPath, verifierInput);
    repository.worker_input = workerInputPath;
    repository.verifier_input = verifierInputPath;
    repository.repair_attempts = nextAttempt;
    addEvidence(manifest2, workerInputPath, verifierInputPath);
    addExecutionEvent(manifest2, repository.name, "repair-prepared", "failed", "running", nextAttempt, occurredAt);
    await assertValid4("runtime-manifest", manifest2);
    await writeJsonAtomic(manifestPath, manifest2);
    return { status: "prepared", attempt: nextAttempt, maximum_attempts: maximumAttempts, manifest: manifestPath, worker_input: workerInputPath, verifier_input: verifierInputPath };
  });
}
function reviewBody(brief, verifier) {
  const acceptance = brief.acceptance_criteria.map((criterion) => `- [x] ${criterion}`).join("\n");
  const checks = verifier.checks.length > 0 ? verifier.checks.map((check) => `- ${check}`).join("\n") : "- No repository command was configured; verifier evidence is recorded.";
  return `## Summary

${brief.requested_outcome}

## Acceptance

${acceptance}

## Verification

${checks}

Prepared from run \`${brief.run_id}\`. No push or pull request was performed.
`;
}
function shellQuote(value2) {
  return `'${value2.replaceAll("'", `'"'"'`)}'`;
}
function command(description, cwd, argv) {
  return { description, cwd, argv, shell: `cd -- ${shellQuote(cwd)} && ${argv.map(shellQuote).join(" ")}` };
}
async function assertBranchName(repository, branch, label) {
  try {
    await git(repository, ["check-ref-format", "--branch", branch]);
  } catch {
    throw new Error(`${label} is not a valid Git branch name`);
  }
}
async function prepareReview(options) {
  const workspaceRoot20 = resolve15(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot20, join10(workspaceRoot20, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const manifestPath = assertInside(runtimeRoot, join10(runtimeRoot, "runs", options.runId, "manifest.json"));
  const lockPath = `${manifestPath}.lock`;
  const config = await loadWorkspace(workspaceRoot20);
  return withExclusiveFile(lockPath, async () => {
    const manifest2 = await readJsonRegularInside(runtimeRoot, manifestPath, "Runtime manifest");
    await assertValid4("runtime-manifest", manifest2);
    if (manifest2.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository2(manifest2, options.repository);
    if (manifest2.status !== "passed") throw new Error(`Draft review preparation requires passed status, received ${manifest2.status}`);
    const taskBrief = await readJsonRegularInside(runtimeRoot, manifest2.task_brief, "Task brief");
    await assertValid4("task-brief", taskBrief);
    const workerInput = await readJsonRegularInside(runtimeRoot, repository.worker_input, "Worker input");
    const verifierInput = await readJsonRegularInside(runtimeRoot, repository.verifier_input, "Verifier input");
    assertInside(runtimeRoot, workerInput.result_path);
    assertInside(runtimeRoot, verifierInput.result_path);
    const worker = await assertCurrentWorker(runtimeRoot, manifest2, repository, workerInput);
    const verifier = await readJsonRegularInside(runtimeRoot, verifierInput.result_path, "Verifier result");
    await assertValid4("verifier-result", verifier);
    assertVerifier(manifest2, repository, taskBrief, verifier);
    if (verifier.status !== "pass") throw new Error(`Draft review preparation requires a passing verifier result, received ${verifier.status}`);
    const headCommit = await git(repository.worktree, ["rev-parse", "HEAD"]);
    const commits = (await git(repository.worktree, ["rev-list", "--reverse", `${repository.base_commit}..${headCommit}`])).split("\n").filter(Boolean);
    const changedFiles = (await git(repository.worktree, ["diff", "--name-only", `${repository.base_commit}...${headCommit}`])).split("\n").filter(Boolean);
    if (commits.length === 0 || changedFiles.length === 0) throw new Error("Draft review preparation requires committed changes");
    if (repository.review_preparation) {
      const existing = await readJsonRegularInside(runtimeRoot, repository.review_preparation, "Review preparation");
      await assertValid4("review-preparation", existing);
      if (existing.contract_version === 2 && existing.head_commit === headCommit && existing.worker_result === workerInput.result_path && existing.verifier_result === verifierInput.result_path) {
        return existing;
      }
    }
    const remotes = (await git(repository.worktree, ["remote"])).split("\n").filter(Boolean);
    const remote = remotes.includes("origin") ? "origin" : null;
    const baseBranch = config.repositories[repository.name]?.default_branch ?? config.workspace.default_branch;
    await assertBranchName(repository.worktree, repository.branch, "Recorded source branch");
    await assertBranchName(repository.worktree, baseBranch, "Configured default branch");
    const taskTarget = taskBrief.repositories.find((candidate) => candidate.name === repository.name);
    const testArgv = taskTarget?.verification_commands ?? taskBrief.verification_commands;
    const commands = {
      diff: command("Inspect the exact base-to-head diff", repository.worktree, ["git", "diff", "--stat", `${repository.base_commit}...${headCommit}`]),
      commits: command("Inspect the exact commit list", repository.worktree, ["git", "log", "--oneline", `${repository.base_commit}..${headCommit}`]),
      show: command("Inspect the exact verified head commit", repository.worktree, ["git", "show", "--stat", "--oneline", headCommit]),
      tests: testArgv.map((value2) => command(`Run recorded verification: ${value2}`, repository.worktree, ["sh", "-lc", value2])),
      switch_target: command("Switch the base repository to the configured target branch", assertInside(workspaceRoot20, join10(workspaceRoot20, repository.base_path)), ["git", "switch", baseBranch]),
      merge: command("Human-only merge of the exact verified head", assertInside(workspaceRoot20, join10(workspaceRoot20, repository.base_path)), ["git", "merge", "--no-ff", headCommit])
    };
    const confirmArgv = ["node", ".agents/bin/cc.mjs", "confirm-merge", "--run-id", manifest2.run_id, "--repository", repository.name, "--merge-commit", "<full-merge-commit>", "--author", "<author-slug>", "--evidence", "<single-line-human-merge-evidence>"];
    const preparedAt = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    const preparationPath = join10(runtimeRoot, "runs", options.runId, `${repository.name}-draft-pr.json`);
    const preparation = {
      contract_version: 2,
      work_id: manifest2.work_id,
      run_id: manifest2.run_id,
      repository: repository.name,
      status: remote ? "ready-for-publication" : "ready-for-local-review",
      remote,
      base_branch: baseBranch,
      head_branch: repository.branch,
      base_commit: repository.base_commit,
      head_commit: headCommit,
      commits,
      changed_files: changedFiles,
      title: `${manifest2.work_id}: ${taskBrief.requested_outcome}`,
      body: reviewBody(taskBrief, verifier),
      worker_result: workerInput.result_path,
      verifier_result: verifierInput.result_path,
      blockers: [],
      prepared_at: preparedAt,
      commands,
      merge_handoff: { status: "merge-confirmation-required", confirmation_argv: confirmArgv, confirmation_shell: confirmArgv.map(shellQuote).join(" ") }
    };
    if (!sameMembers2(worker.changed_files, changedFiles)) {
      throw new Error("Current Git diff does not match the recorded worker result");
    }
    await assertValid4("review-preparation", preparation);
    await writeJsonAtomic(preparationPath, preparation);
    repository.review_preparation = preparationPath;
    repository.review_state = preparation.status;
    addEvidence(manifest2, preparationPath);
    const eventKey2 = `${manifest2.run_id}:execution:${repository.name}:review-prepared:attempt-${repository.repair_attempts ?? 0}`;
    if (!manifest2.execution_events?.some((event) => event.idempotency_key === eventKey2)) {
      addExecutionEvent(manifest2, repository.name, "review-prepared", "passed", "passed", repository.repair_attempts ?? 0, preparedAt, preparationPath);
    }
    await assertValid4("runtime-manifest", manifest2);
    await writeJsonAtomic(manifestPath, manifest2);
    return preparation;
  });
}
function safeReviewEvidence(value2, label) {
  const trimmed = value2.trim();
  if (!trimmed || /[\r\n]/.test(trimmed)) throw new Error(`${label} must be a non-empty single line`);
  if (/-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----|https?:\/\/[^\s/@:]+:[^\s/@]+@/i.test(trimmed)) throw new Error(`${label} appears to contain credentials`);
  return trimmed;
}
async function recordReviewPublication(options) {
  const workspaceRoot20 = resolve15(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot20, join10(workspaceRoot20, ".runtime"));
  const manifestPath = assertInside(runtimeRoot, join10(runtimeRoot, "runs", options.runId, "manifest.json"));
  return withExclusiveFile(`${manifestPath}.lock`, async () => {
    const manifest2 = await readJsonRegularInside(runtimeRoot, manifestPath, "Runtime manifest");
    await assertValid4("runtime-manifest", manifest2);
    if (manifest2.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository2(manifest2, options.repository);
    if (!repository.review_preparation) throw new Error("Prepare the draft pull-request handoff before recording publication");
    const preparation = await readJsonRegularInside(runtimeRoot, repository.review_preparation, "Review preparation");
    await assertValid4("review-preparation", preparation);
    if (preparation.work_id !== manifest2.work_id || preparation.run_id !== manifest2.run_id || preparation.repository !== repository.name) throw new Error("Review preparation identity does not match the active run");
    if (preparation.contract_version !== 2 || preparation.status !== "ready-for-publication" || preparation.remote !== "origin") {
      throw new Error(`Remote publication requires a ready-for-publication handoff with origin; current state is ${preparation.status}`);
    }
    if (!options.authorized) throw new Error("Remote publication recording requires explicit authorization confirmation");
    await assertCleanRepository(repository.worktree);
    const headCommit = await git(repository.worktree, ["rev-parse", "HEAD"]);
    if (headCommit !== preparation.head_commit) throw new Error("Worktree HEAD changed after review preparation");
    const evidence = safeReviewEvidence(options.evidence, "Publication evidence");
    const pullRequest = options.status === "published" ? safeReviewEvidence(options.pullRequest ?? "", "Pull-request reference") : null;
    if (options.status === "failed" && options.pullRequest) throw new Error("Failed publication cannot record a pull-request reference");
    const recordPath2 = join10(runtimeRoot, "runs", options.runId, `${repository.name}-review-publication.json`);
    const record = {
      contract_version: options.status === "published" ? 2 : 1,
      work_id: manifest2.work_id,
      run_id: manifest2.run_id,
      repository: repository.name,
      status: options.status,
      tool: options.tool,
      pull_request: pullRequest,
      evidence,
      head_commit: headCommit,
      idempotency_key: `${manifest2.run_id}:review-publication:${repository.name}`,
      recorded_at: (options.now ?? /* @__PURE__ */ new Date()).toISOString(),
      ...options.status === "published" ? { review_state: "published-for-review" } : {}
    };
    await assertValid4("review-publication-record", record);
    if (repository.review_publication) {
      const existing = await readJsonRegularInside(runtimeRoot, repository.review_publication, "Review publication record");
      await assertValid4("review-publication-record", existing);
      const comparable = (value2) => JSON.stringify({ ...value2, recorded_at: null });
      if (comparable(existing) !== comparable(record)) throw new Error("Review publication was already recorded with different confirmed evidence");
      return existing;
    }
    await writeJsonAtomic(recordPath2, record);
    repository.review_publication = recordPath2;
    repository.review_state = options.status === "published" ? "published-for-review" : "ready-for-publication";
    addEvidence(manifest2, recordPath2);
    manifest2.updated_at = record.recorded_at;
    await assertValid4("runtime-manifest", manifest2);
    await writeJsonAtomic(manifestPath, manifest2);
    return record;
  });
}
async function isAncestor(repository, ancestor, descendant) {
  try {
    await git(repository, ["merge-base", "--is-ancestor", ancestor, descendant]);
    return true;
  } catch {
    return false;
  }
}
async function confirmMerge(options) {
  const workspaceRoot20 = resolve15(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot20, join10(workspaceRoot20, ".runtime"));
  const manifestPath = assertInside(runtimeRoot, join10(runtimeRoot, "runs", options.runId, "manifest.json"));
  if (!/^[a-f0-9]{40,64}$/.test(options.mergeCommit)) throw new Error("Merge commit must be a full lowercase Git object ID");
  const author = safeReviewEvidence(options.author, "Author");
  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(author)) throw new Error("Author must be a lowercase slug");
  const evidence = safeReviewEvidence(options.evidence, "Merge evidence");
  const config = await loadWorkspace(workspaceRoot20);
  return withExclusiveFile(`${manifestPath}.lock`, async () => {
    const manifest2 = await readJsonRegularInside(runtimeRoot, manifestPath, "Runtime manifest");
    await assertValid4("runtime-manifest", manifest2);
    if (manifest2.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository2(manifest2, options.repository);
    if (!repository.review_preparation) throw new Error("Prepare review before confirming a merge");
    const preparation = await readJsonRegularInside(runtimeRoot, repository.review_preparation, "Review preparation");
    await assertValid4("review-preparation", preparation);
    if (preparation.work_id !== manifest2.work_id || preparation.run_id !== manifest2.run_id || preparation.repository !== repository.name) throw new Error("Review preparation identity does not match the active run");
    if (preparation.contract_version !== 2 || preparation.head_commit !== await git(repository.worktree, ["rev-parse", "HEAD"])) throw new Error("Review preparation does not match the current verified head");
    const baseBranch = config.repositories[repository.name]?.default_branch ?? config.workspace.default_branch;
    if (preparation.base_branch !== baseBranch) throw new Error("Review preparation target differs from the configured default branch");
    const baseRepository = assertInside(workspaceRoot20, join10(workspaceRoot20, repository.base_path));
    await git(baseRepository, ["cat-file", "-e", `${options.mergeCommit}^{commit}`]);
    const targetRefs = [`refs/heads/${baseBranch}`, `refs/remotes/origin/${baseBranch}`];
    let targetRef = null;
    let targetCommit = null;
    for (const ref of targetRefs) {
      try {
        const commit = await git(baseRepository, ["rev-parse", "--verify", `${ref}^{commit}`]);
        if (await isAncestor(baseRepository, options.mergeCommit, commit)) {
          targetRef = ref;
          targetCommit = commit;
          break;
        }
      } catch {
      }
    }
    if (!targetRef || !targetCommit) throw new Error(`Reported merge commit is not reachable from the configured default target ${baseBranch}`);
    if (!await isAncestor(baseRepository, preparation.head_commit, options.mergeCommit)) throw new Error("Verified review head is not reachable from the reported merge commit");
    if (!await isAncestor(baseRepository, preparation.base_commit, options.mergeCommit)) throw new Error("Recorded base is not reachable from the reported merge commit");
    const finishArgv = ["node", ".agents/bin/cc.mjs", "finish-work", "--run-id", manifest2.run_id, "--repository", repository.name, "--outcome", "merged", "--author", author, "--merge-commit", options.mergeCommit];
    const publication = repository.review_publication ? await readJsonRegularInside(runtimeRoot, repository.review_publication, "Review publication record") : null;
    if (publication?.status === "published" && publication.pull_request) finishArgv.push("--pull-request", publication.pull_request);
    const record = {
      contract_version: 1,
      work_id: manifest2.work_id,
      run_id: manifest2.run_id,
      repository: repository.name,
      status: "closeout-ready",
      base_branch: baseBranch,
      target_ref: targetRef,
      target_commit: targetCommit,
      head_commit: preparation.head_commit,
      merge_commit: options.mergeCommit,
      evidence,
      finish_work_argv: finishArgv,
      finish_work_shell: finishArgv.map(shellQuote).join(" "),
      idempotency_key: `${manifest2.run_id}:merge-confirmation:${repository.name}`,
      confirmed_at: (options.now ?? /* @__PURE__ */ new Date()).toISOString()
    };
    await assertValid4("merge-confirmation-record", record);
    const recordPath2 = join10(runtimeRoot, "runs", options.runId, `${repository.name}-merge-confirmation.json`);
    if (repository.merge_confirmation) {
      const existing = await readJsonRegularInside(runtimeRoot, repository.merge_confirmation, "Merge confirmation record");
      await assertValid4("merge-confirmation-record", existing);
      const comparable = (value2) => JSON.stringify({ ...value2, confirmed_at: null });
      if (comparable(existing) !== comparable(record)) throw new Error("Merge was already confirmed with different evidence");
      return existing;
    }
    await writeJsonAtomic(recordPath2, record);
    repository.merge_confirmation = recordPath2;
    repository.review_state = "closeout-ready";
    addEvidence(manifest2, recordPath2);
    manifest2.updated_at = record.confirmed_at;
    await assertValid4("runtime-manifest", manifest2);
    await writeJsonAtomic(manifestPath, manifest2);
    return record;
  });
}
var import_yaml8;
var init_review_lifecycle = __esm({
  "scripts/lib/review-lifecycle.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    import_yaml8 = __toESM(require_dist(), 1);
    init_git();
    init_io();
    init_validation();
  }
});

// scripts/prepare-repair.ts
var prepare_repair_exports = {};
import { dirname as dirname11, resolve as resolve16 } from "node:path";
import { parseArgs as parseArgs6 } from "node:util";
import { fileURLToPath as fileURLToPath7 } from "node:url";
var values6, workspaceRoot6;
var init_prepare_repair = __esm({
  async "scripts/prepare-repair.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_review_lifecycle();
    ({ values: values6 } = parseArgs6({
      options: {
        "run-id": { type: "string" },
        repository: { type: "string" }
      }
    }));
    if (!values6["run-id"] || !values6.repository) throw new Error("Usage: prepare-repair --run-id <id> --repository <name>");
    workspaceRoot6 = resolve16(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve16(dirname11(fileURLToPath7(import.meta.url)), ".."));
    console.log(JSON.stringify(await prepareRepair({ workspaceRoot: workspaceRoot6, runId: values6["run-id"], repository: values6.repository }), null, 2));
  }
});

// scripts/prepare-review.ts
var prepare_review_exports = {};
import { dirname as dirname12, resolve as resolve17 } from "node:path";
import { parseArgs as parseArgs7 } from "node:util";
import { fileURLToPath as fileURLToPath8 } from "node:url";
var values7, workspaceRoot7;
var init_prepare_review = __esm({
  async "scripts/prepare-review.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_review_lifecycle();
    ({ values: values7 } = parseArgs7({
      options: {
        "run-id": { type: "string" },
        repository: { type: "string" }
      }
    }));
    if (!values7["run-id"] || !values7.repository) throw new Error("Usage: prepare-review --run-id <id> --repository <name>");
    workspaceRoot7 = resolve17(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve17(dirname12(fileURLToPath8(import.meta.url)), ".."));
    console.log(JSON.stringify(await prepareReview({ workspaceRoot: workspaceRoot7, runId: values7["run-id"], repository: values7.repository }), null, 2));
  }
});

// scripts/record-review-publication.ts
var record_review_publication_exports = {};
import { dirname as dirname13, resolve as resolve18 } from "node:path";
import { parseArgs as parseArgs8 } from "node:util";
import { fileURLToPath as fileURLToPath9 } from "node:url";
var workspaceRoot8, values8;
var init_record_review_publication = __esm({
  async "scripts/record-review-publication.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_review_lifecycle();
    workspaceRoot8 = resolve18(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve18(dirname13(fileURLToPath9(import.meta.url)), ".."));
    ({ values: values8 } = parseArgs8({ options: {
      "run-id": { type: "string" },
      repository: { type: "string" },
      status: { type: "string" },
      tool: { type: "string" },
      "pull-request": { type: "string" },
      evidence: { type: "string" },
      authorized: { type: "boolean" }
    } }));
    if (!values8["run-id"] || !values8.repository || !values8.evidence || !["published", "failed"].includes(values8.status ?? "") || !["gh", "glab", "manual"].includes(values8.tool ?? "")) {
      throw new Error("Usage: cc record-review-publication --run-id <id> --repository <name> --status <published|failed> --tool <gh|glab|manual> [--pull-request <ref>] --evidence <text> --authorized");
    }
    console.log(JSON.stringify(await recordReviewPublication({
      workspaceRoot: workspaceRoot8,
      runId: values8["run-id"],
      repository: values8.repository,
      status: values8.status,
      tool: values8.tool,
      ...values8["pull-request"] ? { pullRequest: values8["pull-request"] } : {},
      evidence: values8.evidence,
      ...values8.authorized ? { authorized: true } : {}
    }), null, 2));
  }
});

// scripts/confirm-merge.ts
var confirm_merge_exports = {};
import { dirname as dirname14, resolve as resolve19 } from "node:path";
import { parseArgs as parseArgs9 } from "node:util";
import { fileURLToPath as fileURLToPath10 } from "node:url";
var values9, workspaceRoot9;
var init_confirm_merge = __esm({
  async "scripts/confirm-merge.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_review_lifecycle();
    ({ values: values9 } = parseArgs9({ options: {
      "run-id": { type: "string" },
      repository: { type: "string" },
      "merge-commit": { type: "string" },
      evidence: { type: "string" },
      author: { type: "string" }
    } }));
    if (!values9["run-id"] || !values9.repository || !values9["merge-commit"] || !values9.evidence || !values9.author) {
      throw new Error("Usage: cc confirm-merge --run-id <id> --repository <name> --merge-commit <full-sha> --author <slug> --evidence <single-line-evidence>");
    }
    workspaceRoot9 = resolve19(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve19(dirname14(fileURLToPath10(import.meta.url)), ".."));
    console.log(JSON.stringify(await confirmMerge({
      workspaceRoot: workspaceRoot9,
      runId: values9["run-id"],
      repository: values9.repository,
      mergeCommit: values9["merge-commit"],
      evidence: values9.evidence,
      author: values9.author
    }), null, 2));
  }
});

// scripts/lib/finish-work.ts
import { access as access5, lstat as lstat6, mkdir as mkdir4, readFile as readFile13, readdir as readdir4, realpath as realpath5 } from "node:fs/promises";
import { basename as basename2, join as join11, relative as relative6, resolve as resolve20 } from "node:path";
async function assertValid5(name, value2) {
  const errors2 = await validateContract(name, value2);
  if (errors2.length > 0) throw new Error(`Invalid ${name}: ${errors2.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
}
function findRepository3(manifest2, name) {
  const repository = manifest2.repositories.find((candidate) => candidate.name === name);
  if (!repository) throw new Error(`Run ${manifest2.run_id} has no repository named ${name}`);
  return repository;
}
function compactTimestamp(date) {
  return date.toISOString().replace(/[-:]/g, "").replace(/\.\d{3}Z$/, "Z");
}
function safeToken(value2, label) {
  const normalized = value2.toLowerCase().trim().replace(/[^a-z0-9]+/g, "-").replace(/^-|-$/g, "");
  if (!normalized || !/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(normalized)) throw new Error(`${label} must contain letters or numbers`);
  return normalized;
}
function taskSlug(brief) {
  return safeToken(brief.requested_outcome, "Task outcome").slice(0, 48).replace(/-$/, "") || "work";
}
function list(items, empty) {
  return items.length > 0 ? items.map((item) => `- ${item}`).join("\n") : `- ${empty}`;
}
function shellQuote2(value2) {
  return `'${value2.replaceAll("'", `'\\''`)}'`;
}
function cleanupRerun(record) {
  const argv = [
    "node",
    ".agents/bin/cc.mjs",
    "finish-work",
    "--run-id",
    record.run_id,
    "--repository",
    record.repository,
    "--outcome",
    record.outcome,
    "--author",
    record.author
  ];
  if (record.reason) argv.push("--reason", record.reason);
  for (const pullRequest of record.pull_requests) argv.push("--pull-request", pullRequest);
  if (record.merge_commit) argv.push("--merge-commit", record.merge_commit);
  argv.push("--cleanup");
  return argv.map(shellQuote2).join(" ");
}
function contributionDocument(manifest2, repository, brief, record) {
  const outcome = record.outcome === "merged" ? "Merged after human review." : `Deliberately abandoned by the human.${record.reason ? ` ${record.reason}` : ""}`;
  const changed = record.changed_files.length > 0 ? ` Changed files: ${record.changed_files.join(", ")}.` : " No product files changed.";
  const planReference = brief.plan.reference ?? "none";
  return `# ${manifest2.work_id}: ${brief.requested_outcome}

- Run: \`${manifest2.run_id}\`
- Task source: ${manifest2.source_kind}
- Plan: ${planReference === "none" ? "none" : `\`${planReference}\``}
- Author: \`${record.author}\`

## Outcome

${outcome}

## Affected repositories

- \`${repository.name}\` on branch \`${repository.branch}\`.${changed}

## Pull requests and commits

${list(record.pull_requests.map((item) => `Pull request: ${item}`), "No pull-request reference was recorded.")}
${list(record.commits.map((item) => `Commit: \`${item}\``), `No commits beyond base \`${repository.base_commit}\`.`)}
- Recorded head: \`${record.head_commit}\`

## Verification

${list(record.verification, "No verifier evidence was available.")}

## Decisions and deviations

- ${record.outcome === "merged" ? "No closeout deviation was recorded." : "The run was deliberately abandoned instead of merged."}

## Remaining risks and follow-up

- ${record.reason ?? "No closeout-specific follow-up was recorded."}

## Candidate durable learnings

- Review this contribution during the next context synchronization; no canonical-context change is asserted automatically.
`;
}
function contributionDocumentErrors(path2, content, runId) {
  const errors2 = [];
  if (!/^\d{8}T\d{6}Z-[a-z0-9]+(?:-[a-z0-9]+)*-[a-z0-9]+(?:-[a-z0-9]+)*\.md$/.test(basename2(path2))) {
    errors2.push("contribution filename must be <UTC timestamp>-<author>-<slug>.md");
  }
  for (const heading of contributionHeadings) if (!content.includes(`${heading}
`)) errors2.push(`contribution is missing ${heading}`);
  if (runId && !content.includes(`- Run: \`${runId}\``)) errors2.push("contribution does not reference the expected run");
  if (/-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----|https?:\/\/[^\s/@:]+:[^\s/@]+@/i.test(content)) {
    errors2.push("contribution appears to contain a credential or private key");
  }
  return errors2;
}
async function loadWorkspace2(workspaceRoot20) {
  const config = (0, import_yaml9.parse)(await readFile13(join11(workspaceRoot20, "workspace.yaml"), "utf8"));
  await assertValid5("workspace", config);
  const errors2 = workspaceSemanticErrors(config);
  if (errors2.length > 0) throw new Error(`Invalid workspace: ${errors2.join("; ")}`);
  return config;
}
async function optionalVerifier(runtimeRoot, manifest2, repository) {
  try {
    const input = await readJsonRegularInside(runtimeRoot, repository.verifier_input, "Verifier input");
    const resultPath = assertInside(runtimeRoot, input.result_path);
    const result3 = await readJsonRegularInside(runtimeRoot, resultPath, "Verifier result");
    await assertValid5("verifier-result", result3);
    if (result3.work_id !== manifest2.work_id || result3.run_id !== manifest2.run_id || result3.repository !== repository.name) {
      throw new Error("Verifier result identity does not match the closeout run");
    }
    return result3;
  } catch (error) {
    if (error.code === "ENOENT") return null;
    throw error;
  }
}
async function assertCurrentWorker2(runtimeRoot, manifest2, repository, headCommit, commits, changedFiles) {
  const input = await readJsonRegularInside(runtimeRoot, repository.worker_input, "Worker input");
  const worker = await readJsonRegularInside(runtimeRoot, input.result_path, "Worker result");
  await assertValid5("worker-result", worker);
  if (worker.work_id !== manifest2.work_id || worker.run_id !== manifest2.run_id || worker.repository !== repository.name) {
    throw new Error("Worker result identity does not match the closeout run");
  }
  if (worker.status !== "completed" || worker.branch !== repository.branch || resolve20(worker.worktree) !== resolve20(repository.worktree)) {
    throw new Error("Closeout requires the completed worker recorded for this branch and worktree");
  }
  if (worker.commits.at(-1) !== headCommit || worker.commits.join("\n") !== commits.join("\n")) {
    throw new Error("Worktree commits changed after the recorded worker result");
  }
  if (worker.changed_files.slice().sort().join("\n") !== changedFiles.slice().sort().join("\n")) {
    throw new Error("Worktree changed-file set differs from the recorded worker result");
  }
}
async function findExistingContribution(root, runId, repository) {
  try {
    for (const entry of await readdir4(root, { withFileTypes: true })) {
      if (!entry.isFile() || !entry.name.endsWith(".md")) continue;
      const path2 = join11(root, entry.name);
      const content = await readFile13(path2, "utf8");
      if (content.includes(`- Run: \`${runId}\``) && content.includes(`- \`${repository}\` on branch`)) return path2;
    }
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
  }
  return null;
}
async function ensureContributionRoot(workspaceRoot20, path2) {
  await mkdir4(path2, { recursive: true, mode: 493 });
  const info = await lstat6(path2);
  if (!info.isDirectory() || info.isSymbolicLink()) throw new Error(`Contribution path must be a real directory: ${path2}`);
  assertInside(await realpath5(workspaceRoot20), await realpath5(path2));
}
function addExecutionEvent2(manifest2, repository, stage, from, to, occurredAt, resultPath) {
  const key = `${manifest2.run_id}:execution:${repository.name}:${stage}`;
  if (manifest2.execution_events?.some((event) => event.idempotency_key === key)) return;
  manifest2.execution_events ??= [];
  manifest2.execution_events.push({
    stage,
    repository: repository.name,
    from_status: from,
    to_status: to,
    inferred: false,
    attempt: repository.repair_attempts ?? 0,
    result_path: resultPath,
    idempotency_key: key,
    occurred_at: occurredAt
  });
  repository.status = to;
  const statuses = manifest2.repositories.map((candidate) => candidate.status ?? manifest2.status);
  if (statuses.every((status2) => status2 === "closed")) manifest2.status = "closed";
  else if (statuses.some((status2) => status2 === "closing" || status2 === "closed")) manifest2.status = "closing";
  else if (statuses.every((status2) => status2 === "passed")) manifest2.status = "passed";
  else if (statuses.includes("failed")) manifest2.status = "failed";
  else if (statuses.includes("blocked")) manifest2.status = "blocked";
  manifest2.updated_at = occurredAt;
}
function addLifecycleEvent(manifest2, outcome, occurredAt) {
  const event = outcome === "merged" ? "task.completed" : "task.cancelled";
  const key = `${manifest2.run_id}:lifecycle:${event}:activity-none`;
  if (manifest2.lifecycle_events.some((item) => item.idempotency_key === key)) return;
  manifest2.lifecycle_events.push({ event, status: "skipped", idempotency_key: key, occurred_at: occurredAt });
}
function assertCloseoutLifecycleReady(manifest2, config, outcome) {
  if (config.activity.provider === "none") return;
  const event = outcome === "merged" ? "task.completed" : "task.cancelled";
  const lifecycle = manifest2.lifecycle_events.find((item) => item.event === event);
  if (!lifecycle) {
    throw new Error(`Prepare configured activity hooks before closeout: node .agents/bin/cc.mjs prepare-lifecycle --run-id ${manifest2.run_id} --event ${event}`);
  }
  if (lifecycle.status !== "completed" && lifecycle.status !== "skipped") {
    throw new Error(`Configured activity hook ${event} is ${lifecycle.status}; complete required or manual actions before closeout`);
  }
}
async function isAncestor2(repository, ancestor, descendant) {
  try {
    await git(repository, ["merge-base", "--is-ancestor", ancestor, descendant]);
    return true;
  } catch {
    return false;
  }
}
async function assertVerifiedMergeConfirmation(workspaceRoot20, runtimeRoot, config, manifest2, repository, requestedMergeCommit) {
  if (!repository.merge_confirmation) throw new Error(`Merged closeout is not ready: record and verify the human merge with node .agents/bin/cc.mjs confirm-merge --run-id ${manifest2.run_id} --repository ${repository.name} --merge-commit <full-sha> --author <slug> --evidence <single-line-evidence>`);
  const record = await readJsonRegularInside(runtimeRoot, repository.merge_confirmation, "Merge confirmation record");
  await assertValid5("merge-confirmation-record", record);
  if (record.work_id !== manifest2.work_id || record.run_id !== manifest2.run_id || record.repository !== repository.name || record.head_commit !== await git(repository.worktree, ["rev-parse", "HEAD"])) {
    throw new Error("Merge confirmation identity or verified head does not match the active run");
  }
  if (requestedMergeCommit && requestedMergeCommit !== record.merge_commit) throw new Error("Requested merge commit differs from the verified merge confirmation");
  const defaultBranch = config.repositories[repository.name]?.default_branch ?? config.workspace.default_branch;
  if (record.base_branch !== defaultBranch || ![`refs/heads/${defaultBranch}`, `refs/remotes/origin/${defaultBranch}`].includes(record.target_ref)) {
    throw new Error("Merge confirmation does not target the configured default branch");
  }
  const baseRepository = assertInside(workspaceRoot20, join11(workspaceRoot20, repository.base_path));
  const currentTarget = await git(baseRepository, ["rev-parse", "--verify", `${record.target_ref}^{commit}`]);
  if (currentTarget !== record.target_commit || !await isAncestor2(baseRepository, record.merge_commit, currentTarget) || !await isAncestor2(baseRepository, record.head_commit, record.merge_commit)) {
    throw new Error("Merge confirmation no longer proves the exact reviewed head is reachable from the recorded default target");
  }
  return record;
}
async function verifiedDefaultRefs(repository, branch) {
  const refs = [];
  for (const ref of [`refs/heads/${branch}`, `refs/remotes/origin/${branch}`]) {
    try {
      await git(repository, ["rev-parse", "--verify", ref]);
      refs.push(ref);
    } catch {
    }
  }
  return refs;
}
async function cleanupBlockers(workspaceRoot20, config, repository, record) {
  const blockers2 = [];
  const contributionPath = assertInside(workspaceRoot20, join11(workspaceRoot20, record.contribution));
  try {
    await git(workspaceRoot20, ["ls-files", "--error-unmatch", "--", record.contribution]);
    if (await git(workspaceRoot20, ["status", "--porcelain=v1", "--", record.contribution])) {
      blockers2.push("Contribution has uncommitted wrapper changes; commit it through the configured wrapper workflow before cleanup.");
    }
  } catch {
    blockers2.push("Contribution is not durably tracked by wrapper Git; commit it before cleanup.");
  }
  try {
    await access5(contributionPath);
  } catch {
    blockers2.push("Contribution file is missing; runtime cleanup would discard the only closeout record.");
  }
  const baseRepository = assertInside(workspaceRoot20, join11(workspaceRoot20, repository.base_path));
  try {
    if (await git(baseRepository, ["status", "--porcelain=v1", "--untracked-files=normal"])) blockers2.push("Base repository is dirty.");
    if (await git(repository.worktree, ["status", "--porcelain=v1", "--untracked-files=normal"])) blockers2.push("Run worktree has uncommitted changes.");
    if (await git(repository.worktree, ["branch", "--show-current"]) !== repository.branch) blockers2.push("Run worktree is on an unexpected branch.");
    if (await git(repository.worktree, ["rev-parse", "HEAD"]) !== record.head_commit) blockers2.push("Run worktree HEAD changed after closeout preparation.");
  } catch (error) {
    blockers2.push(`Run worktree is unavailable: ${error.message}`);
    return blockers2;
  }
  const defaultBranch = config.repositories[repository.name]?.default_branch ?? config.workspace.default_branch;
  const defaultRefs = await verifiedDefaultRefs(baseRepository, defaultBranch);
  const headOnDefault = (await Promise.all(defaultRefs.map((ref) => isAncestor2(baseRepository, record.head_commit, ref)))).some(Boolean);
  if (record.outcome === "merged") {
    let mergeEvidence = headOnDefault;
    if (!mergeEvidence && record.merge_commit) {
      const mergeOnDefault = (await Promise.all(defaultRefs.map((ref) => isAncestor2(baseRepository, record.merge_commit, ref)))).some(Boolean);
      mergeEvidence = mergeOnDefault && await isAncestor2(baseRepository, record.base_commit, record.merge_commit);
    }
    if (!mergeEvidence) blockers2.push("Merged outcome is not reachable from the configured default branch; fetch the merge or provide a verified merge commit.");
  } else if (record.head_commit !== record.base_commit && !headOnDefault) {
    const remoteRefs = await git(baseRepository, ["for-each-ref", "--format=%(refname)", "--contains", record.head_commit, "refs/remotes/"]);
    if (!remoteRefs) blockers2.push("Abandoned branch contains commits that are neither merged nor preserved by a remote ref.");
  }
  return blockers2;
}
async function closePreparedRun(workspaceRoot20, manifestPath, manifest2, repository, recordPath2, record, config, occurredAt) {
  const detected = await cleanupBlockers(workspaceRoot20, config, repository, record);
  if (detected.length > 0) {
    const checklist = detected.map((blocker, index) => `${index + 1}. ${blocker}`);
    checklist.push(`${checklist.length + 1}. After resolving the blockers, rerun exactly: ${cleanupRerun(record)}`);
    const blocked = { ...record, status: "blocked", cleanup: { ...record.cleanup, requested: true }, blockers: checklist, updated_at: occurredAt };
    await assertValid5("closeout-record", blocked);
    await writeJsonAtomic(recordPath2, blocked);
    return blocked;
  }
  const baseRepository = assertInside(workspaceRoot20, join11(workspaceRoot20, repository.base_path));
  await git(baseRepository, ["worktree", "remove", repository.worktree]);
  const closed = {
    ...record,
    status: "closed",
    cleanup: { requested: true, worktree_removed: true, branch_preserved: true, runtime_evidence_preserved: true },
    blockers: [],
    updated_at: occurredAt
  };
  await assertValid5("closeout-record", closed);
  await writeJsonAtomic(recordPath2, closed);
  addExecutionEvent2(manifest2, repository, "closeout-cleaned", "closing", "closed", occurredAt, recordPath2);
  await assertValid5("runtime-manifest", manifest2);
  await writeJsonAtomic(manifestPath, manifest2);
  return closed;
}
async function finishWork(options) {
  const workspaceRoot20 = resolve20(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot20, join11(workspaceRoot20, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const manifestPath = assertInside(runtimeRoot, join11(runtimeRoot, "runs", options.runId, "manifest.json"));
  const lockPath = `${manifestPath}.lock`;
  const author = safeToken(options.author, "Author");
  const invocationTime = options.now ?? /* @__PURE__ */ new Date();
  if (options.mergeCommit && !/^[a-f0-9]{40,64}$/.test(options.mergeCommit)) throw new Error("--merge-commit must be a full lowercase Git object ID");
  if (options.pullRequests?.some((reference2) => !reference2.trim() || /[\r\n]/.test(reference2))) throw new Error("Pull-request references must be non-empty single lines");
  const config = await loadWorkspace2(workspaceRoot20);
  const wrapperTopLevel = await git(workspaceRoot20, ["rev-parse", "--show-toplevel"]);
  if (await realpath5(wrapperTopLevel) !== await realpath5(workspaceRoot20)) throw new Error("Workspace root must be the wrapper Git root before closeout");
  return withExclusiveFile(lockPath, async () => {
    const manifest2 = await readJsonRegularInside(runtimeRoot, manifestPath, "Runtime manifest");
    await assertValid5("runtime-manifest", manifest2);
    if (manifest2.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    assertCloseoutLifecycleReady(manifest2, config, options.outcome);
    const repository = findRepository3(manifest2, options.repository);
    const recordPath2 = assertInside(runtimeRoot, join11(runtimeRoot, "runs", options.runId, `${repository.name}-closeout.json`));
    if (repository.closeout_record) {
      const existing = await readJsonRegularInside(runtimeRoot, repository.closeout_record, "Closeout record");
      await assertValid5("closeout-record", existing);
      if (existing.outcome !== options.outcome || existing.author !== author) throw new Error("Closeout was already prepared with different human intent");
      if (existing.outcome === "merged") await assertVerifiedMergeConfirmation(workspaceRoot20, runtimeRoot, config, manifest2, repository, options.mergeCommit);
      if (existing.status === "closed" || !options.cleanup) return existing;
      return closePreparedRun(workspaceRoot20, manifestPath, manifest2, repository, recordPath2, existing, config, invocationTime.toISOString());
    }
    const repositoryStatus = repository.status ?? manifest2.status;
    if (!["passed", "failed", "blocked", "cancelled"].includes(repositoryStatus)) throw new Error(`Closeout preparation requires a terminal repository outcome, received ${repositoryStatus}`);
    if (options.outcome === "merged" && repositoryStatus !== "passed") throw new Error(`Merged closeout requires a passed repository, received ${repositoryStatus}`);
    if (options.outcome === "abandoned" && !options.reason?.trim()) throw new Error("Deliberate abandonment requires --reason");
    const mergeConfirmation = options.outcome === "merged" ? await assertVerifiedMergeConfirmation(workspaceRoot20, runtimeRoot, config, manifest2, repository, options.mergeCommit) : null;
    const brief = await readJsonRegularInside(runtimeRoot, manifest2.task_brief, "Task brief");
    await assertValid5("task-brief", brief);
    const headCommit = await git(repository.worktree, ["rev-parse", "HEAD"]);
    if (await git(repository.worktree, ["branch", "--show-current"]) !== repository.branch) throw new Error("Run worktree is on an unexpected branch");
    const commits = (await git(repository.worktree, ["rev-list", "--reverse", `${repository.base_commit}..${headCommit}`])).split("\n").filter(Boolean);
    const changedFiles = (await git(repository.worktree, ["diff", "--name-only", `${repository.base_commit}...${headCommit}`])).split("\n").filter(Boolean);
    if (options.outcome === "merged") await assertCurrentWorker2(runtimeRoot, manifest2, repository, headCommit, commits, changedFiles);
    const verifier = await optionalVerifier(runtimeRoot, manifest2, repository);
    if (options.outcome === "merged" && verifier?.status !== "pass") throw new Error("Merged closeout requires the recorded passing verifier result");
    const verification = verifier ? [verifier.summary, ...verifier.checks, ...verifier.acceptance.map((item) => `${item.criterion}: ${item.status} \u2014 ${item.evidence}`)] : [];
    const preparedAt = invocationTime.toISOString();
    const contributionsRoot = assertInside(workspaceRoot20, join11(workspaceRoot20, "contributions", "general"));
    await ensureContributionRoot(workspaceRoot20, contributionsRoot);
    const existingContribution = await findExistingContribution(contributionsRoot, options.runId, repository.name);
    const contributionPath = existingContribution ?? join11(contributionsRoot, `${compactTimestamp(invocationTime)}-${author}-${taskSlug(brief)}-${repository.name}.md`);
    const contributionRelative = relative6(workspaceRoot20, contributionPath).replaceAll("\\", "/");
    const record = {
      contract_version: 1,
      work_id: manifest2.work_id,
      run_id: manifest2.run_id,
      repository: repository.name,
      outcome: options.outcome,
      status: "prepared",
      author,
      reason: options.reason?.trim() || null,
      contribution: contributionRelative,
      branch: repository.branch,
      base_commit: repository.base_commit,
      head_commit: headCommit,
      merge_commit: mergeConfirmation?.merge_commit ?? options.mergeCommit ?? null,
      pull_requests: [...new Set(options.pullRequests ?? [])],
      commits,
      changed_files: changedFiles,
      verification,
      cleanup: { requested: Boolean(options.cleanup), worktree_removed: false, branch_preserved: true, runtime_evidence_preserved: true },
      blockers: [],
      prepared_at: preparedAt,
      updated_at: preparedAt
    };
    const document = contributionDocument(manifest2, repository, brief, record);
    const documentErrors = contributionDocumentErrors(contributionPath, document, options.runId);
    if (documentErrors.length > 0) throw new Error(`Invalid contribution: ${documentErrors.join("; ")}`);
    if (existingContribution && await readFile13(existingContribution, "utf8") !== document) {
      throw new Error("An append-only contribution already exists for this run with different closeout content");
    }
    if (!existingContribution) await writeTextExclusive(contributionPath, document);
    await assertValid5("closeout-record", record);
    await writeJsonAtomic(recordPath2, record);
    repository.closeout_record = recordPath2;
    repository.contribution = contributionRelative;
    if (!manifest2.evidence.includes(recordPath2)) manifest2.evidence.push(recordPath2);
    const fromStatus = repositoryStatus;
    addExecutionEvent2(manifest2, repository, "closeout-prepared", fromStatus, "closing", preparedAt, recordPath2);
    if (config.activity.provider === "none") addLifecycleEvent(manifest2, options.outcome, preparedAt);
    await assertValid5("runtime-manifest", manifest2);
    await writeJsonAtomic(manifestPath, manifest2);
    if (!options.cleanup) return record;
    return closePreparedRun(workspaceRoot20, manifestPath, manifest2, repository, recordPath2, record, config, preparedAt);
  });
}
var import_yaml9, contributionHeadings;
var init_finish_work = __esm({
  "scripts/lib/finish-work.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    import_yaml9 = __toESM(require_dist(), 1);
    init_git();
    init_io();
    init_validation();
    contributionHeadings = [
      "## Outcome",
      "## Affected repositories",
      "## Pull requests and commits",
      "## Verification",
      "## Decisions and deviations",
      "## Remaining risks and follow-up",
      "## Candidate durable learnings"
    ];
  }
});

// scripts/finish-work.ts
var finish_work_exports = {};
import { dirname as dirname15, resolve as resolve21 } from "node:path";
import { parseArgs as parseArgs10 } from "node:util";
import { fileURLToPath as fileURLToPath11 } from "node:url";
var values10, workspaceRoot10, result;
var init_finish_work2 = __esm({
  async "scripts/finish-work.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_finish_work();
    ({ values: values10 } = parseArgs10({
      options: {
        "run-id": { type: "string" },
        repository: { type: "string" },
        outcome: { type: "string" },
        author: { type: "string" },
        reason: { type: "string" },
        "merge-commit": { type: "string" },
        "pull-request": { type: "string", multiple: true },
        cleanup: { type: "boolean", default: false }
      }
    }));
    if (!values10["run-id"] || !values10.repository || !values10.outcome || !values10.author) {
      throw new Error("Required: --run-id <id> --repository <name> --outcome <merged|abandoned> --author <slug>");
    }
    if (values10.outcome !== "merged" && values10.outcome !== "abandoned") throw new Error("--outcome must be merged or abandoned");
    workspaceRoot10 = resolve21(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve21(dirname15(fileURLToPath11(import.meta.url)), ".."));
    result = await finishWork({
      workspaceRoot: workspaceRoot10,
      runId: values10["run-id"],
      repository: values10.repository,
      outcome: values10.outcome,
      author: values10.author,
      cleanup: values10.cleanup,
      ...values10.reason ? { reason: values10.reason } : {},
      ...values10["merge-commit"] ? { mergeCommit: values10["merge-commit"] } : {},
      ...values10["pull-request"] ? { pullRequests: values10["pull-request"] } : {}
    });
    console.log(JSON.stringify(result, null, 2));
  }
});

// scripts/create-plan.ts
var create_plan_exports = {};
import { readFile as readFile14 } from "node:fs/promises";
import { dirname as dirname16, resolve as resolve22 } from "node:path";
import { parseArgs as parseArgs11 } from "node:util";
import { fileURLToPath as fileURLToPath12 } from "node:url";
var values11, workspaceRoot11, inputPath, request2;
var init_create_plan = __esm({
  async "scripts/create-plan.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_plans();
    ({ values: values11 } = parseArgs11({
      options: { input: { type: "string" } }
    }));
    if (!values11.input) throw new Error("Usage: cc create-plan --input <plan-draft-request.json>");
    workspaceRoot11 = resolve22(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve22(dirname16(fileURLToPath12(import.meta.url)), ".."));
    inputPath = resolve22(process.cwd(), values11.input);
    request2 = JSON.parse(await readFile14(inputPath, "utf8"));
    console.log(JSON.stringify(await createPlanDraft(workspaceRoot11, request2), null, 2));
  }
});

// scripts/validate-plan.ts
var validate_plan_exports = {};
import { resolve as resolve23 } from "node:path";
import { parseArgs as parseArgs12 } from "node:util";
var positionals2, planDirectory, result2;
var init_validate_plan = __esm({
  async "scripts/validate-plan.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_plans();
    ({ positionals: positionals2 } = parseArgs12({ allowPositionals: true }));
    if (!positionals2[0]) throw new Error("Usage: cc validate-plan context/plans/<plan-id>");
    planDirectory = resolve23(process.cwd(), positionals2[0]);
    result2 = await validatePlanDirectory(planDirectory);
    if (result2.errors.length > 0) {
      console.error(`Invalid plan ${planDirectory}:`);
      for (const error of result2.errors) console.error(`- ${error}`);
      process.exitCode = 1;
    } else {
      console.log(`Valid ${result2.index?.status} plan ${result2.index?.plan_id} version ${result2.index?.plan_version}`);
    }
  }
});

// scripts/set-plan-state.ts
var set_plan_state_exports = {};
import { dirname as dirname17, join as join12, resolve as resolve24 } from "node:path";
import { parseArgs as parseArgs13 } from "node:util";
import { fileURLToPath as fileURLToPath13 } from "node:url";
var values12, requested, transition, workspaceRoot12, planDirectory2;
var init_set_plan_state = __esm({
  async "scripts/set-plan-state.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_io();
    init_plans();
    ({ values: values12 } = parseArgs13({
      options: {
        plan: { type: "string" },
        "approve-by": { type: "string" },
        "material-revision": { type: "string" },
        "non-material-repair": { type: "boolean", default: false }
      }
    }));
    if (!values12.plan) throw new Error("Usage: cc set-plan-state --plan context/plans/<plan-id> (--approve-by <name> | --material-revision <reason> | --non-material-repair)");
    requested = [Boolean(values12["approve-by"]), Boolean(values12["material-revision"]), values12["non-material-repair"]].filter(Boolean).length;
    if (requested !== 1) throw new Error("Choose exactly one plan state transition");
    if (values12["approve-by"]) transition = { kind: "approve", approved_by: values12["approve-by"] };
    else if (values12["material-revision"]) transition = { kind: "material-revision", reason: values12["material-revision"] };
    else transition = { kind: "non-material-repair" };
    workspaceRoot12 = resolve24(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve24(dirname17(fileURLToPath13(import.meta.url)), ".."));
    planDirectory2 = assertInside(join12(workspaceRoot12, "context", "plans"), resolve24(process.cwd(), values12.plan));
    console.log(JSON.stringify(await setPlanState(planDirectory2, transition), null, 2));
  }
});

// scripts/lib/whats-next.ts
import { lstat as lstat7, readdir as readdir5, readFile as readFile15, realpath as realpath6 } from "node:fs/promises";
import { join as join13, relative as relative7, resolve as resolve25, sep as sep2 } from "node:path";
function contractMessages2(errors2) {
  return errors2.map((error) => `${error.instancePath || "/"} ${error.message}`);
}
async function isDirectory(path2) {
  try {
    const info = await lstat7(path2);
    return info.isDirectory() && !info.isSymbolicLink();
  } catch {
    return false;
  }
}
function reference(workspaceRoot20, path2) {
  const raw = relative7(workspaceRoot20, path2);
  return raw && !raw.startsWith(`..${sep2}`) && raw !== ".." ? raw.replaceAll("\\", "/") : path2;
}
function inside(root, path2) {
  const resolvedRoot = resolve25(root);
  const resolvedPath = resolve25(path2);
  return resolvedPath === resolvedRoot || resolvedPath.startsWith(`${resolvedRoot}${sep2}`);
}
function acceptanceIsSufficient(raw) {
  return /^- (?!None recorded\.$).+/m.test(raw);
}
function rank(candidate) {
  if (candidate.kind === "reconciliation") return 0;
  if (candidate.urgent) return 1;
  if (candidate.state === "in-progress" || candidate.state === "closeout") return 2;
  if (["review", "verification-failure", "ci-failure"].includes(candidate.kind) || candidate.state === "failed" || candidate.state === "review") return 3;
  if (candidate.kind !== "plan-work-item") return 4;
  return 5;
}
function blockers(candidate, currentUser) {
  const values20 = [];
  if (candidate.kind === "reconciliation") values20.push("read-only sources report contradictory work states");
  if (candidate.owner && candidate.owner !== currentUser) values20.push(`owned by another active contributor: ${candidate.owner}`);
  for (const dependency of candidate.dependencies) {
    if (dependency.state !== "completed") values20.push(`dependency ${dependency.reference} is ${dependency.state}`);
  }
  if (candidate.plan_approval_state !== "approved" && candidate.plan_approval_state !== "not-applicable") {
    values20.push(`governing plan is ${candidate.plan_approval_state}`);
  }
  if (!candidate.scope_sufficient) values20.push("scope is insufficient");
  if (!candidate.acceptance_sufficient) values20.push("acceptance criteria are insufficient");
  if (candidate.repositories.length === 0) values20.push("no affected repository is resolved");
  if (!candidate.access_available) values20.push("required repository access is unavailable");
  if (candidate.contract_blocked) values20.push("an unresolved contract decision blocks implementation");
  return [...new Set(values20)];
}
function actionKind(candidate, candidateBlockers) {
  if (candidate.kind === "reconciliation") return "reconcile";
  if (candidateBlockers.length > 0) return "enable";
  if (candidate.state === "review") return "review";
  if (candidate.state === "closeout") return "closeout";
  return "execute";
}
function actionFor(candidate, candidateBlockers, config) {
  const action = actionKind(candidate, candidateBlockers);
  const stateSources = candidate.state_sources?.length ? candidate.state_sources.map((item) => `state ${item.state}: ${item.source_reference}`) : [`state ${candidate.state}: ${candidate.source_reference}`];
  const evidence = [
    ...stateSources,
    `plan approval: ${candidate.plan_approval_state}`,
    `plan approval source: ${candidate.plan_reference ?? candidate.source_reference}`,
    candidate.dependencies.length === 0 ? "dependencies: none" : `dependencies: ${candidate.dependencies.map((item) => `${item.reference}=${item.state}`).join(", ")}`,
    `dependency source: ${candidate.plan_reference ?? candidate.source_reference}`,
    `scope sufficient: ${candidate.scope_sufficient}`,
    `scope source: ${candidate.source_reference}`,
    `acceptance sufficient: ${candidate.acceptance_sufficient}`,
    `acceptance source: ${candidate.source_reference}`,
    `repository access available: ${candidate.access_available}`,
    "repository access source: workspace.yaml#repositories"
  ];
  const sequence = candidate.repositories.flatMap((repository) => {
    const agent = config.repositories[repository]?.agent ?? repository;
    return [`repository worker (${agent})`, "independent verifier"];
  });
  const title = action === "execute" ? candidate.title : actionTitle(candidate, candidateBlockers, action);
  return {
    action,
    candidate_id: candidate.candidate_id,
    title,
    why: action === "execute" ? `${rankingReason(candidate)}; all readiness checks passed.` : action === "review" ? "Verified implementation evidence is ready for human review or merge preparation; implementation must not be duplicated." : action === "closeout" ? "Implementation has advanced beyond execution and the remaining work is human-gated closeout or cleanup." : action === "reconcile" ? "Configured activity and local outcome evidence disagree; reconcile the cited sources without mutating them or starting duplicate implementation." : `No candidate is currently executable; this is the smallest visible action that addresses the first blocker for ${candidate.title}.`,
    readiness_evidence: evidence,
    source_references: [.../* @__PURE__ */ new Set([candidate.source_reference, ...candidate.plan_reference ? [candidate.plan_reference] : [], ...(candidate.state_sources ?? []).map((item) => item.source_reference)])],
    repositories: candidate.repositories,
    agent_sequence: action === "execute" ? [...new Set(sequence)] : [],
    blockers: candidateBlockers,
    risks: candidate.risks
  };
}
function rankingReason(candidate) {
  if (candidate.urgent) return "It is explicitly urgent";
  if (candidate.state === "in-progress") return "It is actionable work already in progress";
  if (["review", "verification-failure", "ci-failure"].includes(candidate.kind) || candidate.state === "failed") return "It addresses review, verification, or CI feedback";
  if (candidate.kind !== "plan-work-item") return "It is the highest-priority approved ready source candidate";
  return "It is the next dependency-ready item in an approved plan";
}
function actionTitle(candidate, candidateBlockers, action) {
  if (action === "review") return `Review or prepare merge for ${candidate.title}`;
  if (action === "closeout") return `Complete closeout or cleanup for ${candidate.title}`;
  if (action === "reconcile") return `Reconcile contradictory state for ${candidate.title}`;
  const first = candidateBlockers[0] ?? "readiness is not established";
  if (first.startsWith("governing plan is draft")) return `Approve the governing plan for ${candidate.title}`;
  if (first.startsWith("dependency ")) return `Resolve or confirm ${first.replace(" is ", " as ")}`;
  if (first.startsWith("owned by another")) return `Confirm ownership before starting ${candidate.title}`;
  if (first.includes("contract")) return `Resolve the blocking contract for ${candidate.title}`;
  if (first.includes("repository")) return `Register or restore repository access for ${candidate.title}`;
  return `Clarify ${first} for ${candidate.title}`;
}
async function repositoryAccess(workspaceRoot20, config, repositories) {
  if (repositories.length === 0) return false;
  for (const name of repositories) {
    const repository = config.repositories[name];
    if (!repository || !await isDirectory(resolve25(workspaceRoot20, repository.path))) return false;
  }
  return true;
}
function runtimeState(manifest2, workId) {
  const item = manifest2.plan_work_items?.find((candidate) => candidate.work_id === workId);
  const reviewStates = manifest2.repositories.map((repository) => repository.review_state);
  if (manifest2.status === "cancelled") return "cancelled";
  if (manifest2.status === "closing") return "closeout";
  if (manifest2.status === "closed") return "completed";
  if (reviewStates.includes("closeout-ready")) return "closeout";
  if (reviewStates.some((state) => state === "ready-for-local-review" || state === "ready-for-publication" || state === "published-for-review" || state === "merge-confirmation-required")) return "review";
  if (manifest2.status === "passed" || item?.outcome === "passed") return "review";
  if (manifest2.status === "failed" || manifest2.status === "blocked" || item?.outcome === "failed" || item?.outcome === "blocked") return "failed";
  return "in-progress";
}
async function validatedJson(name, path2) {
  const info = await lstat7(path2);
  if (!info.isFile() || info.isSymbolicLink()) throw new Error(`${name} is not a regular file`);
  const value2 = JSON.parse(await readFile15(path2, "utf8"));
  const errors2 = contractMessages2(await validateContract(name, value2));
  if (errors2.length > 0) throw new Error(errors2.join("; "));
  return value2;
}
async function discoverRuntimeObservations(workspaceRoot20) {
  const observations = /* @__PURE__ */ new Map();
  const runs = /* @__PURE__ */ new Map();
  const warnings = [];
  const runsRoot = join13(workspaceRoot20, ".runtime", "runs");
  try {
    const rootInfo = await lstat7(runsRoot);
    if (rootInfo.isSymbolicLink()) return { observations, runs, warnings: ["Ignored symlinked runtime runs directory: .runtime/runs"] };
    if (!rootInfo.isDirectory()) return { observations, runs, warnings };
  } catch {
    return { observations, runs, warnings };
  }
  const entries = (await readdir5(runsRoot, { withFileTypes: true })).sort((a, b) => a.name.localeCompare(b.name));
  for (const entry of entries) {
    if (entry.isSymbolicLink()) {
      warnings.push(`Ignored symlinked runtime evidence ${entry.name}`);
      continue;
    }
    if (!entry.isDirectory()) continue;
    const manifestPath = join13(runsRoot, entry.name, "manifest.json");
    try {
      const info = await lstat7(manifestPath);
      if (!info.isFile() || info.isSymbolicLink()) throw new Error("manifest is not a regular file");
    } catch (error) {
      if (error.code === "ENOENT") continue;
      warnings.push(`Ignored malformed runtime evidence ${entry.name}: ${error.message}`);
      continue;
    }
    try {
      const manifest2 = await validatedJson("runtime-manifest", manifestPath);
      if (manifest2.source_kind !== "plan" || !manifest2.plan_work_items) continue;
      if (!inside(workspaceRoot20, manifest2.task_brief)) throw new Error("task brief escapes the workspace");
      if (!inside(await realpath6(workspaceRoot20), await realpath6(manifest2.task_brief))) throw new Error("task brief resolves outside the workspace");
      const brief = await validatedJson("task-brief", manifest2.task_brief);
      if (brief.source.kind !== "plan" || brief.plan.approval_state !== "approved") throw new Error("task brief is not an approved plan task");
      if (manifest2.run_id !== brief.run_id || manifest2.work_id !== brief.work_id || entry.name !== manifest2.run_id) throw new Error("run or work identity differs between runtime path, manifest, and task brief");
      if (manifest2.plan_work_items.length !== 1 || manifest2.repositories.length !== 1 || brief.repositories.length !== 1 || brief.plan.work_ids.length !== 1) throw new Error("plan runtime must contain one work item and one repository");
      const item = manifest2.plan_work_items[0];
      const repository = manifest2.repositories[0];
      if (item.work_id !== manifest2.work_id || brief.plan.work_ids[0] !== item.work_id || item.repository !== repository.name || brief.repositories[0].name !== item.repository) throw new Error("work item or repository identity differs between manifest and task brief");
      if (brief.source.reference !== brief.plan.reference) throw new Error("task brief plan references differ");
      const planReference = brief.plan.reference.replace(/\/README\.md$/, "").replace(/\/$/, "");
      const planDirectory3 = resolve25(workspaceRoot20, planReference);
      const plansRoot = resolve25(workspaceRoot20, "context", "plans");
      if (!inside(plansRoot, planDirectory3) || planDirectory3 === plansRoot) throw new Error("task brief plan reference is outside context/plans");
      const plansInfo = await lstat7(plansRoot);
      if (!plansInfo.isDirectory() || plansInfo.isSymbolicLink() || !inside(await realpath6(plansRoot), await realpath6(planDirectory3))) throw new Error("task brief plan reference resolves outside a real context/plans directory");
      const plan = await validatePlanDirectory(planDirectory3);
      if (!plan.index || !plan.work_breakdown || plan.errors.length > 0 || plan.index.status !== "approved") throw new Error(`current plan is not valid and approved: ${plan.errors.join("; ")}`);
      if (plan.index.plan_version !== brief.plan.plan_version || plan.index.approved_digest !== brief.plan.approved_digest) throw new Error("task brief approval version or digest is stale");
      const currentItem = plan.work_breakdown.items.find((candidate) => candidate.work_id === item.work_id);
      if (!currentItem || currentItem.repository !== item.repository) throw new Error("runtime work item does not match the current approved plan");
      let state = runtimeState(manifest2, item.work_id);
      let closeout;
      let closeoutReference;
      if (repository.closeout_record) {
        if (!inside(workspaceRoot20, repository.closeout_record)) throw new Error("closeout record escapes the workspace");
        if (!inside(await realpath6(workspaceRoot20), await realpath6(repository.closeout_record))) throw new Error("closeout record resolves outside the workspace");
        closeout = await validatedJson("closeout-record", repository.closeout_record);
        if (closeout.run_id !== manifest2.run_id || closeout.work_id !== item.work_id || closeout.repository !== item.repository) throw new Error("closeout identity does not match manifest");
        state = closeout.status === "closed" ? closeout.outcome === "merged" ? "completed" : "cancelled" : "closeout";
        closeoutReference = reference(workspaceRoot20, repository.closeout_record);
      }
      const manifestReference = reference(workspaceRoot20, manifestPath);
      const briefReference = reference(workspaceRoot20, manifest2.task_brief);
      const sources = [manifestReference, briefReference, ...closeoutReference ? [closeoutReference] : []];
      const values20 = observations.get(item.work_id) ?? [];
      for (const source_reference of sources) values20.push({ state, source_reference, precedence: state === "completed" || state === "cancelled" ? 30 : 20, plan_reference: planReference, run_id: manifest2.run_id });
      observations.set(item.work_id, values20);
      runs.set(manifest2.run_id, { manifest: manifest2, brief, manifest_reference: manifestReference, brief_reference: briefReference, ...closeout ? { closeout } : {}, ...closeoutReference ? { closeout_reference: closeoutReference } : {} });
    } catch (error) {
      warnings.push(`Ignored malformed runtime evidence ${entry.name}: ${error.message}`);
    }
  }
  return { observations, runs, warnings };
}
async function discoverDurableContributions(workspaceRoot20, runs) {
  const observations = /* @__PURE__ */ new Map();
  const warnings = [];
  const root = join13(workspaceRoot20, "contributions");
  try {
    const rootInfo = await lstat7(root);
    if (rootInfo.isSymbolicLink()) return { observations, warnings: ["Ignored symlinked contributions directory: contributions"] };
    if (!rootInfo.isDirectory()) return { observations, warnings };
  } catch {
    return { observations, warnings };
  }
  const groups = (await readdir5(root, { withFileTypes: true })).sort((a, b) => a.name.localeCompare(b.name));
  for (const group of groups) {
    if (group.isSymbolicLink()) {
      warnings.push(`Ignored symlinked contribution group ${group.name}`);
      continue;
    }
    if (!group.isDirectory()) continue;
    const directory = join13(root, group.name);
    const files = (await readdir5(directory, { withFileTypes: true })).sort((a, b) => a.name.localeCompare(b.name));
    for (const file of files) {
      if (file.isSymbolicLink()) {
        warnings.push(`Ignored symlinked contribution ${group.name}/${file.name}`);
        continue;
      }
      if (!file.isFile() || !file.name.endsWith(".md")) continue;
      const path2 = join13(directory, file.name);
      const source_reference = reference(workspaceRoot20, path2);
      try {
        await git(workspaceRoot20, ["ls-files", "--error-unmatch", "--", source_reference]);
      } catch {
        continue;
      }
      const content = await readFile15(path2, "utf8");
      const work = content.match(/^# ([A-Z][A-Z0-9]{1,15}-\d{3,}):/m)?.[1];
      const run = content.match(/^- Run: `([^`]+)`$/m)?.[1];
      const merged = content.includes("Merged after human review.");
      const abandoned = content.includes("Deliberately abandoned by the human.");
      const documentErrors = contributionDocumentErrors(path2, content, run);
      if (!work || !run || merged === abandoned || documentErrors.length > 0) {
        warnings.push(`Ignored unrecognized durable contribution ${source_reference}`);
        continue;
      }
      const validatedRun = runs.get(run);
      const runtimeRepository = validatedRun?.manifest.repositories[0];
      const closeout = validatedRun?.closeout;
      if (!validatedRun || validatedRun.manifest.work_id !== work || !closeout || runtimeRepository?.contribution !== source_reference || closeout.contribution !== source_reference || closeout.outcome !== (merged ? "merged" : "abandoned")) {
        warnings.push(`Ignored unassociated durable contribution ${source_reference}: cited run and closeout relationship are not validated`);
        continue;
      }
      const values20 = observations.get(work) ?? [];
      values20.push({ state: merged ? "completed" : "cancelled", source_reference, precedence: 40, plan_reference: validatedRun.brief.plan.reference.replace(/\/README\.md$/, "").replace(/\/$/, ""), run_id: run });
      observations.set(work, values20);
    }
  }
  return { observations, warnings };
}
function project(observations) {
  if (observations.length === 0) return null;
  const ordered2 = observations.slice().sort((a, b) => b.precedence - a.precedence || a.source_reference.localeCompare(b.source_reference));
  return { state: ordered2[0].state, observations: ordered2, contradiction: new Set(ordered2.map((item) => item.state)).size > 1 };
}
async function discoverPlanCandidates(workspaceRoot20, config, activityFacts, localObservations) {
  const plansRoot = join13(workspaceRoot20, "context", "plans");
  if (!await isDirectory(plansRoot)) return { candidates: [], warnings: [], matchedFacts: /* @__PURE__ */ new Set() };
  const entries = (await readdir5(plansRoot, { withFileTypes: true })).filter((entry) => entry.isDirectory() && !entry.isSymbolicLink()).sort((a, b) => a.name.localeCompare(b.name));
  const candidates = [];
  const warnings = [];
  const matchedFacts = /* @__PURE__ */ new Set();
  const validations = /* @__PURE__ */ new Map();
  const workIdCounts = /* @__PURE__ */ new Map();
  for (const entry of entries) {
    const planDirectory3 = join13(plansRoot, entry.name);
    const validation = await validatePlanDirectory(planDirectory3);
    if (!validation.index || !validation.work_breakdown || validation.errors.length > 0) {
      warnings.push(`Skipped invalid plan ${entry.name}: ${validation.errors.join("; ") || "missing parsed plan material"}`);
      continue;
    }
    validations.set(entry.name, validation);
    for (const item of validation.work_breakdown.items) workIdCounts.set(item.work_id, (workIdCounts.get(item.work_id) ?? 0) + 1);
  }
  for (const entry of entries) {
    const planDirectory3 = join13(plansRoot, entry.name);
    const validation = validations.get(entry.name);
    if (!validation?.index || !validation.work_breakdown) continue;
    const planDirectoryReference = relative7(workspaceRoot20, planDirectory3).replaceAll("\\", "/");
    const projectedByWork = /* @__PURE__ */ new Map();
    for (const item of validation.work_breakdown.items) {
      const relevantLocal = (localObservations.get(item.work_id) ?? []).filter((observation) => {
        if (observation.plan_reference) return observation.plan_reference === planDirectoryReference;
        if (observation.precedence === 40 && workIdCounts.get(item.work_id) > 1) return false;
        return true;
      });
      const activity2 = activityFacts.get(item.work_id) ?? [];
      const value2 = project([...relevantLocal, ...activity2.map((fact) => ({ state: fact.state, source_reference: fact.source_reference, precedence: 10 }))]);
      if (value2) projectedByWork.set(item.work_id, value2);
      if (workIdCounts.get(item.work_id) > 1 && (localObservations.get(item.work_id) ?? []).some((observation) => !observation.plan_reference && observation.precedence === 40)) {
        warnings.push(`Ignored ambiguous durable contribution for ${item.work_id}: multiple plans use that work ID`);
      }
    }
    const requirementPath = join13(planDirectory3, "0010-requirements.md");
    const acceptanceSufficient = acceptanceIsSufficient(await readFile15(requirementPath, "utf8"));
    for (const item of validation.work_breakdown.items) {
      const facts = activityFacts.get(item.work_id) ?? [];
      for (const fact2 of facts) matchedFacts.add(fact2.candidate_id);
      const fact = facts[0];
      const projection = projectedByWork.get(item.work_id);
      const repositories = config.repositories[item.repository] ? [item.repository] : [];
      const activityRepositoryMismatch = Boolean(fact?.repositories.length) && (fact.repositories.length !== 1 || fact.repositories[0] !== item.repository);
      const dependencies = item.depends_on.map((dependency) => ({
        reference: dependency,
        state: projectedByWork.get(dependency)?.state === "completed" ? "completed" : projectedByWork.has(dependency) ? "pending" : "unknown"
      }));
      const planReference = relative7(workspaceRoot20, join13(planDirectory3, "README.md"));
      const state = projection?.state ?? "ready";
      const contradiction = Boolean(projection?.contradiction) || Boolean(projection) && validation.index.status !== "approved" || activityRepositoryMismatch;
      const stateSources = projection?.observations.map(({ state: observed, source_reference }) => ({ state: observed, source_reference })) ?? [{ state: "ready", source_reference: `${relative7(workspaceRoot20, join13(planDirectory3, validation.index.work_breakdown))}#${item.work_id}` }];
      candidates.push({
        contract_version: 1,
        candidate_id: `plan:${validation.index.plan_id}:${item.work_id}`,
        kind: contradiction ? "reconciliation" : state === "review" ? "review" : "plan-work-item",
        work_id: item.work_id,
        title: item.title,
        state,
        urgent: fact?.urgent ?? false,
        priority: fact?.priority ?? 0,
        owner: fact?.owner ?? null,
        plan_reference: planReference,
        plan_approval_state: validation.index.status,
        dependencies,
        scope_sufficient: fact?.scope_sufficient ?? (item.title.trim().length > 0 && repositories.length > 0),
        acceptance_sufficient: fact?.acceptance_sufficient ?? acceptanceSufficient,
        repositories,
        access_available: (fact?.access_available ?? true) && await repositoryAccess(workspaceRoot20, config, repositories),
        contract_blocked: fact?.contract_blocked ?? false,
        source_reference: `${relative7(workspaceRoot20, join13(planDirectory3, validation.index.work_breakdown))}#${item.work_id}`,
        state_sources: stateSources,
        risks: [.../* @__PURE__ */ new Set([...fact?.risks ?? [], ...activityRepositoryMismatch ? [`Activity repository evidence does not match approved plan repository ${item.repository}.`] : [], ...contradiction ? ["Starting implementation before reconciliation could duplicate or overwrite completed work."] : []])]
      });
    }
  }
  return { candidates, warnings, matchedFacts };
}
async function recommendWhatsNext(workspaceRootInput, activity2 = null, now = /* @__PURE__ */ new Date()) {
  const workspaceRoot20 = resolve25(workspaceRootInput);
  const config = await readData(join13(workspaceRoot20, "workspace.yaml"));
  const workspaceErrors = contractMessages2(await validateContract("workspace", config));
  workspaceErrors.push(...workspaceSemanticErrors(config));
  if (workspaceErrors.length > 0) throw new Error(`Invalid workspace configuration:
- ${workspaceErrors.join("\n- ")}`);
  if (activity2) {
    const activityErrors = contractMessages2(await validateContract("fake-activity-source", activity2));
    if (activityErrors.length > 0) throw new Error(`Invalid fake activity source:
- ${activityErrors.join("\n- ")}`);
  }
  const currentUser = activity2?.current_user ?? "local-user";
  const facts = /* @__PURE__ */ new Map();
  for (const candidate of activity2?.candidates ?? []) if (candidate.work_id) facts.set(candidate.work_id, [...facts.get(candidate.work_id) ?? [], candidate]);
  const runtime = await discoverRuntimeObservations(workspaceRoot20);
  const durable = await discoverDurableContributions(workspaceRoot20, runtime.runs);
  const localObservations = new Map(runtime.observations);
  for (const [workId, observations] of durable.observations) {
    const terminalByRun = new Map(observations.map((observation) => [observation.run_id, observation.state]));
    const runtimeObservations = (localObservations.get(workId) ?? []).map((observation) => {
      const terminal = terminalByRun.get(observation.run_id);
      return terminal ? { ...observation, state: terminal } : observation;
    });
    localObservations.set(workId, [...runtimeObservations, ...observations]);
  }
  const discovered = await discoverPlanCandidates(workspaceRoot20, config, facts, localObservations);
  const external = (activity2?.candidates ?? []).filter((candidate) => !discovered.matchedFacts.has(candidate.candidate_id));
  const hydratedExternal = [];
  for (const candidate of external) hydratedExternal.push({ ...candidate, state_sources: candidate.state_sources ?? [{ state: candidate.state, source_reference: candidate.source_reference }], access_available: candidate.access_available && await repositoryAccess(workspaceRoot20, config, candidate.repositories) });
  const candidates = [...discovered.candidates, ...hydratedExternal];
  const duplicateIds = candidates.filter((candidate, index) => candidates.findIndex((value2) => value2.candidate_id === candidate.candidate_id) !== index);
  if (duplicateIds.length > 0) throw new Error(`Duplicate candidate ID: ${duplicateIds[0].candidate_id}`);
  const excluded = candidates.filter((candidate) => candidate.state === "completed" || candidate.state === "cancelled");
  const assessed = candidates.filter((candidate) => candidate.state !== "completed" && candidate.state !== "cancelled").map((candidate) => ({ candidate, blockers: blockers(candidate, currentUser), rank: rank(candidate) }));
  const ordered2 = assessed.slice().sort((left, right) => left.rank - right.rank || right.candidate.priority - left.candidate.priority || left.candidate.candidate_id.localeCompare(right.candidate.candidate_id));
  const executable = ordered2.filter((item) => item.blockers.length === 0);
  const blocked = ordered2.filter((item) => item.blockers.length > 0);
  const reconciliation = blocked.filter((item) => item.candidate.kind === "reconciliation");
  let recommendation;
  let alternatives;
  if (reconciliation.length > 0) {
    recommendation = actionFor(reconciliation[0].candidate, reconciliation[0].blockers, config);
    alternatives = [...reconciliation.slice(1), ...executable].slice(0, 2).map((item) => actionFor(item.candidate, item.blockers, config));
  } else if (executable.length > 0) {
    recommendation = actionFor(executable[0].candidate, [], config);
    alternatives = executable.slice(1, 3).map((item) => actionFor(item.candidate, [], config));
  } else if (blocked.length > 0) {
    recommendation = actionFor(blocked[0].candidate, blocked[0].blockers, config);
    alternatives = blocked.slice(1, 3).map((item) => actionFor(item.candidate, item.blockers, config));
  } else {
    recommendation = { action: "enable", candidate_id: null, title: "Create or approve a scoped work source", why: "No executable or blocked candidate was found in the configured read-only sources.", readiness_evidence: ["approved plan candidates none: context/plans", "activity candidates none: workspace.yaml#activity", "active runtime candidates none: .runtime/runs", "durable outcome candidates none: contributions"], source_references: ["context/plans", "workspace.yaml#activity", ".runtime/runs", "contributions"], repositories: [], agent_sequence: [], blockers: ["no available candidate provides sufficient scope and acceptance criteria"], risks: [] };
    alternatives = [];
  }
  const result3 = { contract_version: 1, generated_at: now.toISOString(), recommendation, alternatives, considered: { total: candidates.length, executable: executable.length, blocked: blocked.length, excluded: excluded.length }, warnings: [.../* @__PURE__ */ new Set([...runtime.warnings, ...durable.warnings, ...discovered.warnings])], no_state_changed: true };
  const resultErrors = contractMessages2(await validateContract("whats-next-result", result3));
  if (resultErrors.length > 0) throw new Error(`Generated invalid whats-next result:
- ${resultErrors.join("\n- ")}`);
  return result3;
}
var init_whats_next = __esm({
  "scripts/lib/whats-next.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_finish_work();
    init_git();
    init_plans();
    init_validation();
  }
});

// scripts/whats-next.ts
var whats_next_exports = {};
import { readFile as readFile16 } from "node:fs/promises";
import { dirname as dirname18, resolve as resolve26 } from "node:path";
import { parseArgs as parseArgs14 } from "node:util";
import { fileURLToPath as fileURLToPath14 } from "node:url";
var values13, workspaceRoot13, activity;
var init_whats_next2 = __esm({
  async "scripts/whats-next.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_whats_next();
    ({ values: values13 } = parseArgs14({ options: { "activity-fixture": { type: "string" } } }));
    workspaceRoot13 = resolve26(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve26(dirname18(fileURLToPath14(import.meta.url)), ".."));
    activity = values13["activity-fixture"] ? JSON.parse(await readFile16(resolve26(process.cwd(), values13["activity-fixture"]), "utf8")) : null;
    console.log(JSON.stringify(await recommendWhatsNext(workspaceRoot13, activity), null, 2));
  }
});

// scripts/prepare-lifecycle.ts
var prepare_lifecycle_exports = {};
import { dirname as dirname19, resolve as resolve27 } from "node:path";
import { parseArgs as parseArgs15 } from "node:util";
import { fileURLToPath as fileURLToPath15 } from "node:url";
var events, capabilities, values14, workspaceRoot14;
var init_prepare_lifecycle = __esm({
  async "scripts/prepare-lifecycle.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_activity_lifecycle();
    events = ["task.starting", "task.review-ready", "task.completed", "task.blocked", "task.cancelled"];
    capabilities = ["read-tasks", "update-status", "create-tasks", "assign-task", "timers"];
    ({ values: values14 } = parseArgs15({ options: {
      "run-id": { type: "string" },
      event: { type: "string" },
      available: { type: "string", multiple: true, default: [] }
    } }));
    if (!values14["run-id"] || !events.includes(values14.event) || values14.available.some((item) => !capabilities.includes(item))) {
      throw new Error("Usage: prepare-lifecycle --run-id <id> --event <semantic-event> [--available <capability>]");
    }
    workspaceRoot14 = resolve27(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve27(dirname19(fileURLToPath15(import.meta.url)), ".."));
    console.log(JSON.stringify(await prepareActivityLifecycle({
      workspaceRoot: workspaceRoot14,
      runId: values14["run-id"],
      event: values14.event,
      availableCapabilities: values14.available
    }), null, 2));
  }
});

// scripts/record-lifecycle-action.ts
var record_lifecycle_action_exports = {};
import { dirname as dirname20, resolve as resolve28 } from "node:path";
import { parseArgs as parseArgs16 } from "node:util";
import { fileURLToPath as fileURLToPath16 } from "node:url";
var events2, values15, workspaceRoot15;
var init_record_lifecycle_action = __esm({
  async "scripts/record-lifecycle-action.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_activity_lifecycle();
    events2 = ["task.starting", "task.review-ready", "task.completed", "task.blocked", "task.cancelled"];
    ({ values: values15 } = parseArgs16({ options: {
      "run-id": { type: "string" },
      event: { type: "string" },
      action: { type: "string" },
      status: { type: "string" },
      evidence: { type: "string" },
      reference: { type: "string" }
    } }));
    if (!values15["run-id"] || !events2.includes(values15.event) || !values15.action || !["completed", "failed"].includes(values15.status ?? "") || !values15.evidence) {
      throw new Error("Usage: record-lifecycle-action --run-id <id> --event <event> --action <id> --status <completed|failed> --evidence <text> [--reference <ref>]");
    }
    workspaceRoot15 = resolve28(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve28(dirname20(fileURLToPath16(import.meta.url)), ".."));
    console.log(JSON.stringify(await recordActivityLifecycleAction({
      workspaceRoot: workspaceRoot15,
      runId: values15["run-id"],
      event: values15.event,
      actionId: values15.action,
      status: values15.status,
      evidence: values15.evidence,
      ...values15.reference ? { externalReference: values15.reference } : {}
    }), null, 2));
  }
});

// scripts/lib/plan-publication.ts
import { readFile as readFile17 } from "node:fs/promises";
import { join as join14, resolve as resolve29 } from "node:path";
async function assertValid6(name, value2) {
  const errors2 = await validateContract(name, value2);
  if (errors2.length) throw new Error(`Invalid ${name}: ${errors2.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
}
function safeLine(value2, field) {
  const clean = value2.trim();
  if (!clean || /[|\r\n]/.test(clean)) throw new Error(`${field} must be a non-empty single line without table delimiters`);
  if (/https?:\/\/[^\s/@:]+:[^\s/@]+@|-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----/i.test(clean)) throw new Error(`${field} appears to contain a credential`);
  return clean;
}
function ordered(items) {
  const remaining = [...items];
  const emitted = /* @__PURE__ */ new Set();
  const result3 = [];
  while (remaining.length) {
    const index = remaining.findIndex((item2) => [...item2.depends_on, ...item2.parent ? [item2.parent] : []].every((id) => emitted.has(id)));
    if (index < 0) throw new Error("Plan publication order cannot resolve dependencies and parents");
    const [item] = remaining.splice(index, 1);
    result3.push(item);
    emitted.add(item.work_id);
  }
  return result3;
}
function status(items) {
  const done = items.filter((item) => item.status === "created" || item.status === "existing").length;
  const failed = items.filter((item) => item.status === "failed").length;
  if (done === items.length) return "completed";
  if (failed && done) return "partial";
  if (failed) return "failed";
  return done ? "in-progress" : "proposed";
}
function recordPath(workspaceRoot20, planId) {
  return assertInside(workspaceRoot20, join14(workspaceRoot20, ".runtime", "publications", `${planId}.json`));
}
async function preparePlanPublication(options) {
  const workspaceRoot20 = resolve29(options.workspaceRoot);
  const config = (0, import_yaml10.parse)(await readFile17(join14(workspaceRoot20, "workspace.yaml"), "utf8"));
  await assertValid6("workspace", config);
  const semantic = workspaceSemanticErrors(config);
  if (semantic.length) throw new Error(`Invalid workspace: ${semantic.join("; ")}`);
  await assertValid6("plan-publication-discovery", options.discovery);
  if (config.activity.provider === "none" || config.activity.provider !== options.discovery.provider) throw new Error("Publication discovery provider must match the configured non-none activity provider");
  if (![...config.activity.required_capabilities, ...config.activity.optional_capabilities].includes("create-tasks")) throw new Error("Configured activity provider does not declare create-tasks capability");
  const planDirectory3 = assertInside(workspaceRoot20, join14(workspaceRoot20, "context", "plans", options.planId));
  const plan = await validatePlanDirectory(planDirectory3);
  if (plan.errors.length || !plan.index || !plan.work_breakdown) throw new Error(`Plan is invalid: ${plan.errors.join("; ")}`);
  if (plan.index.status !== "approved" || !plan.index.approved_digest) throw new Error("Only an approved plan can be published");
  const index = plan.index;
  const breakdown = plan.work_breakdown;
  const discovered = /* @__PURE__ */ new Map();
  for (const mapping of options.discovery.mappings) {
    if (discovered.has(mapping.work_id)) throw new Error(`Duplicate discovered mapping: ${mapping.work_id}`);
    discovered.set(mapping.work_id, { reference: safeLine(mapping.external_reference, "External reference"), evidence: safeLine(mapping.evidence, "Evidence") });
  }
  const path2 = recordPath(workspaceRoot20, options.planId);
  await ensurePrivateDirectory(join14(workspaceRoot20, ".runtime", "publications"));
  return withExclusiveFile(`${path2}.lock`, async () => {
    try {
      const existing = JSON.parse(await readFile17(path2, "utf8"));
      await assertValid6("plan-publication-record", existing);
      if (existing.plan_version !== index.plan_version || existing.provider !== options.discovery.provider || existing.destination !== options.discovery.destination) throw new Error("Existing publication record conflicts with this request");
      return existing;
    } catch (error) {
      if (error.code !== "ENOENT") throw error;
    }
    const items = ordered(breakdown.items).map((item) => {
      const known = item.external_reference ? { reference: item.external_reference, evidence: "Confirmed mapping already stored in the approved plan." } : discovered.get(item.work_id);
      if (item.external_reference && discovered.get(item.work_id)?.reference !== void 0 && discovered.get(item.work_id).reference !== item.external_reference) throw new Error(`Conflicting external mapping for ${item.work_id}`);
      return { work_id: item.work_id, title: item.title, parent: item.parent, depends_on: item.depends_on, area: item.area, repository: item.repository, action: known ? "skip-existing" : "create", status: known ? "existing" : "proposed", external_reference: known?.reference ?? null, evidence: known?.evidence ?? null, idempotency_key: `${index.plan_id}:v${index.plan_version}:${item.work_id}` };
    });
    for (const workId of discovered.keys()) if (!items.some((item) => item.work_id === workId)) throw new Error(`Discovered mapping references unknown work ID: ${workId}`);
    const now = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    const record = { contract_version: 2, plan_id: index.plan_id, plan_version: index.plan_version, approved_digest: index.approved_digest, provider: options.discovery.provider, destination: safeLine(options.discovery.destination, "Destination"), status: status(items), items, warnings: [], prepared_at: now, updated_at: now };
    await assertValid6("plan-publication-record", record);
    await writeJsonAtomic(path2, record);
    return record;
  });
}
async function writeMapping(planDirectory3, breakdownName, workId, reference2, now) {
  const path2 = join14(planDirectory3, breakdownName);
  const raw = await readFile17(path2, "utf8");
  let found = false;
  const updated = raw.split("\n").map((line) => {
    if (!line.startsWith(`| ${workId} |`)) return line;
    const cells = line.slice(1, -1).split("|").map((cell) => cell.trim());
    if (![6, 7].includes(cells.length)) throw new Error(`Invalid work-breakdown row for ${workId}`);
    const referenceCell = cells.length - 1;
    if (cells[referenceCell] !== "\u2014" && cells[referenceCell] !== reference2) throw new Error(`Plan already maps ${workId} to a different external reference`);
    cells[referenceCell] = reference2;
    found = true;
    return `| ${cells.join(" | ")} |`;
  }).join("\n");
  if (!found) throw new Error(`Plan has no work item ${workId}`);
  await writeTextAtomic(path2, updated);
  return (await setPlanState(planDirectory3, { kind: "non-material-repair" }, now)).approved_digest;
}
async function recordPlanPublication(options) {
  const workspaceRoot20 = resolve29(options.workspaceRoot);
  const path2 = recordPath(workspaceRoot20, options.planId);
  return withExclusiveFile(`${path2}.lock`, async () => {
    const record = JSON.parse(await readFile17(path2, "utf8"));
    await assertValid6("plan-publication-record", record);
    const item = record.items.find((candidate) => candidate.work_id === options.workId);
    if (!item) throw new Error(`Publication has no work item ${options.workId}`);
    const evidence = safeLine(options.evidence, "Evidence");
    const reference2 = options.externalReference ? safeLine(options.externalReference, "External reference") : null;
    if (options.status === "created" && !reference2) throw new Error("Created publication result requires a confirmed external reference");
    if (item.status === "created" || item.status === "existing") {
      if (item.external_reference === reference2 && item.evidence === evidence) return record;
      throw new Error(`${item.work_id} already has a different confirmed mapping`);
    }
    if (options.status === "failed") {
      item.status = "failed";
      item.evidence = evidence;
      item.external_reference = null;
    } else {
      const planDirectory3 = assertInside(workspaceRoot20, join14(workspaceRoot20, "context", "plans", record.plan_id));
      const validation = await validatePlanDirectory(planDirectory3);
      if (validation.errors.length || !validation.index || validation.index.status !== "approved" || validation.index.plan_version !== record.plan_version) throw new Error("Approved plan changed during publication");
      record.approved_digest = await writeMapping(planDirectory3, validation.index.work_breakdown, item.work_id, reference2, options.now ?? /* @__PURE__ */ new Date());
      item.status = "created";
      item.external_reference = reference2;
      item.evidence = evidence;
    }
    record.status = status(record.items);
    record.updated_at = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    await assertValid6("plan-publication-record", record);
    await writeJsonAtomic(path2, record);
    return record;
  });
}
var import_yaml10;
var init_plan_publication = __esm({
  "scripts/lib/plan-publication.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    import_yaml10 = __toESM(require_dist(), 1);
    init_io();
    init_plans();
    init_validation();
  }
});

// scripts/prepare-plan-publication.ts
var prepare_plan_publication_exports = {};
import { readFile as readFile18 } from "node:fs/promises";
import { dirname as dirname21, resolve as resolve30 } from "node:path";
import { parseArgs as parseArgs17 } from "node:util";
import { fileURLToPath as fileURLToPath17 } from "node:url";
var values16, workspaceRoot16, discovery;
var init_prepare_plan_publication = __esm({
  async "scripts/prepare-plan-publication.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_plan_publication();
    ({ values: values16 } = parseArgs17({ options: { plan: { type: "string" }, discovery: { type: "string" } } }));
    if (!values16.plan || !values16.discovery) throw new Error("Usage: prepare-plan-publication --plan <plan-id> --discovery <json>");
    workspaceRoot16 = resolve30(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve30(dirname21(fileURLToPath17(import.meta.url)), ".."));
    discovery = JSON.parse(await readFile18(resolve30(process.cwd(), values16.discovery), "utf8"));
    console.log(JSON.stringify(await preparePlanPublication({ workspaceRoot: workspaceRoot16, planId: values16.plan, discovery }), null, 2));
  }
});

// scripts/record-plan-publication.ts
var record_plan_publication_exports = {};
import { dirname as dirname22, resolve as resolve31 } from "node:path";
import { parseArgs as parseArgs18 } from "node:util";
import { fileURLToPath as fileURLToPath18 } from "node:url";
var values17, workspaceRoot17;
var init_record_plan_publication = __esm({
  async "scripts/record-plan-publication.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_plan_publication();
    ({ values: values17 } = parseArgs18({ options: { plan: { type: "string" }, work: { type: "string" }, status: { type: "string" }, evidence: { type: "string" }, reference: { type: "string" } } }));
    if (!values17.plan || !values17.work || !["created", "failed"].includes(values17.status ?? "") || !values17.evidence) throw new Error("Usage: record-plan-publication --plan <id> --work <id> --status <created|failed> --evidence <text> [--reference <external-ref>]");
    workspaceRoot17 = resolve31(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve31(dirname22(fileURLToPath18(import.meta.url)), ".."));
    console.log(JSON.stringify(await recordPlanPublication({ workspaceRoot: workspaceRoot17, planId: values17.plan, workId: values17.work, status: values17.status, evidence: values17.evidence, ...values17.reference ? { externalReference: values17.reference } : {} }), null, 2));
  }
});

// scripts/lib/context-sync.ts
import { createHash as createHash4 } from "node:crypto";
import { lstat as lstat8, readFile as readFile19, realpath as realpath7 } from "node:fs/promises";
import { join as join15, resolve as resolve32 } from "node:path";
async function assertValid7(name, value2) {
  const errors2 = await validateContract(name, value2);
  if (errors2.length > 0) throw new Error(`Invalid ${name}: ${errors2.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
}
async function readJson3(path2) {
  return JSON.parse(await readFile19(path2, "utf8"));
}
function compactTimestamp2(date) {
  return date.toISOString().replace(/[-:]/g, "").replace(/\.\d{3}Z$/, "Z");
}
async function loadWorkspace3(workspaceRoot20) {
  const config = (0, import_yaml11.parse)(await readFile19(join15(workspaceRoot20, "workspace.yaml"), "utf8"));
  await assertValid7("workspace", config);
  const errors2 = workspaceSemanticErrors(config);
  if (errors2.length > 0) throw new Error(`Invalid workspace: ${errors2.join("; ")}`);
  return config;
}
async function validateRequestSemantics(workspaceRoot20, config, request4) {
  const listed = new Set(request4.contributions);
  for (const contribution of request4.contributions) {
    const path2 = assertInside(workspaceRoot20, join15(workspaceRoot20, contribution));
    const info = await lstat8(path2);
    if (!info.isFile() || info.isSymbolicLink()) throw new Error(`Contribution must be a regular file: ${contribution}`);
    assertInside(await realpath7(workspaceRoot20), await realpath7(path2));
    const errors2 = contributionDocumentErrors(path2, await readFile19(path2, "utf8"));
    if (errors2.length > 0) throw new Error(`Invalid contribution ${contribution}: ${errors2.join("; ")}`);
  }
  const durableTargets = [];
  for (const proposal of request4.proposals) {
    if (!listed.has(proposal.source_contribution)) throw new Error(`Proposal source is not listed: ${proposal.source_contribution}`);
    if (proposal.classification === "durable-wrapper") {
      if (!proposal.target || !proposal.proposed_change?.trim() || proposal.target_repository) throw new Error("durable-wrapper proposals require target and proposed_change only");
      if ((await readFile19(join15(workspaceRoot20, proposal.target), "utf8")).includes(proposal.source_contribution)) {
        throw new Error(`Contribution is already cited by ${proposal.target}: ${proposal.source_contribution}`);
      }
      durableTargets.push(proposal.target);
    } else if (proposal.classification === "repository-local") {
      if (!proposal.target_repository || !config.repositories[proposal.target_repository] || proposal.target || proposal.proposed_change) {
        throw new Error("repository-local proposals require one registered target_repository and no wrapper target/change");
      }
    } else if (proposal.target || proposal.target_repository || proposal.proposed_change) {
      throw new Error(`${proposal.classification} proposals cannot mutate wrapper or repository context`);
    }
  }
  if (durableTargets.length === 0) throw new Error("Context synchronization has no durable wrapper proposal; report classifications without creating a worktree");
  return [...new Set(durableTargets)].sort();
}
async function prepareContextSync(options) {
  const workspaceRoot20 = resolve32(options.workspaceRoot);
  await assertValid7("context-sync-request", options.request);
  const config = await loadWorkspace3(workspaceRoot20);
  const allowedPaths = await validateRequestSemantics(workspaceRoot20, config, options.request);
  if (await realpath7(await git(workspaceRoot20, ["rev-parse", "--show-toplevel"])) !== await realpath7(workspaceRoot20)) {
    throw new Error("Workspace root must be the wrapper Git root");
  }
  await assertCleanRepository(workspaceRoot20);
  const baseBranch = await git(workspaceRoot20, ["branch", "--show-current"]);
  if (!baseBranch) throw new Error("Context synchronization requires an attached wrapper branch");
  const baseCommit = await git(workspaceRoot20, ["rev-parse", "HEAD"]);
  const now = options.now ?? /* @__PURE__ */ new Date();
  const digest = createHash4("sha256").update(JSON.stringify(options.request)).update(now.toISOString()).digest("hex").slice(0, 8);
  const syncId = `${compactTimestamp2(now)}-${digest}`;
  const runtimeRoot = assertInside(workspaceRoot20, join15(workspaceRoot20, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const syncRoot = assertInside(runtimeRoot, join15(runtimeRoot, "context-sync", syncId));
  await ensurePrivateDirectory(syncRoot);
  const requestPath = join15(syncRoot, "request.json");
  const recordPath2 = join15(syncRoot, "record.json");
  const worktree = assertInside(runtimeRoot, join15(runtimeRoot, "worktrees", "context-sync", syncId, "wrapper"));
  const branch = `agent/context-sync-${syncId.toLowerCase()}`;
  await writeJsonAtomic(requestPath, options.request);
  const preparedAt = now.toISOString();
  const record = {
    contract_version: 1,
    sync_id: syncId,
    status: "prepared",
    wrapper_mode: config.workflow.wrapper_change_policy,
    base_branch: baseBranch,
    base_commit: baseCommit,
    branch,
    worktree,
    request: requestPath,
    allowed_wrapper_paths: allowedPaths,
    repository_follow_ups: options.request.proposals.filter((proposal) => proposal.classification === "repository-local").map((proposal) => ({ repository: proposal.target_repository, summary: proposal.summary, source_contribution: proposal.source_contribution })),
    future_tasks: options.request.proposals.filter((proposal) => proposal.classification === "future-task").map((proposal) => ({ summary: proposal.summary, source_contribution: proposal.source_contribution })),
    retained_one_offs: options.request.proposals.filter((proposal) => proposal.classification === "one-off").map((proposal) => proposal.summary),
    changed_files: [],
    commits: [],
    remote: null,
    blockers: [],
    prepared_at: preparedAt,
    updated_at: preparedAt
  };
  try {
    await ensurePrivateDirectory(join15(runtimeRoot, "worktrees"));
    await ensurePrivateDirectory(join15(runtimeRoot, "worktrees", "context-sync"));
    await ensurePrivateDirectory(join15(runtimeRoot, "worktrees", "context-sync", syncId));
    await git(workspaceRoot20, ["worktree", "add", "-b", branch, worktree, baseCommit]);
  } catch (error) {
    record.status = "blocked";
    record.blockers = [`Wrapper worktree preparation failed: ${error.message}`];
    await writeJsonAtomic(recordPath2, record);
    throw error;
  }
  await assertValid7("context-sync-record", record);
  await writeJsonAtomic(recordPath2, record);
  return record;
}
async function prepareContextReview(options) {
  const workspaceRoot20 = resolve32(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot20, join15(workspaceRoot20, ".runtime"));
  const recordPath2 = assertInside(runtimeRoot, join15(runtimeRoot, "context-sync", options.syncId, "record.json"));
  return withExclusiveFile(`${recordPath2}.lock`, async () => {
    const record = await readJson3(recordPath2);
    await assertValid7("context-sync-record", record);
    if (record.sync_id !== options.syncId) throw new Error("Context sync ID mismatch");
    const worktree = assertInside(runtimeRoot, record.worktree);
    const config = await loadWorkspace3(workspaceRoot20);
    const request4 = await readJson3(assertInside(runtimeRoot, record.request));
    await assertValid7("context-sync-request", request4);
    const recomputedPaths = await validateRequestSemantics(workspaceRoot20, config, request4);
    if (recomputedPaths.join("\n") !== record.allowed_wrapper_paths.slice().sort().join("\n")) throw new Error("Context sync allowed paths do not match the validated request");
    const registrations = await git(workspaceRoot20, ["worktree", "list", "--porcelain"]);
    if (!registrations.split("\n").includes(`worktree ${await realpath7(worktree)}`)) throw new Error("Context sync worktree is not registered by the wrapper repository");
    await assertCleanRepository(worktree);
    if (await git(worktree, ["branch", "--show-current"]) !== record.branch) throw new Error("Wrapper worktree branch changed");
    const head = await git(worktree, ["rev-parse", "HEAD"]);
    await git(worktree, ["merge-base", "--is-ancestor", record.base_commit, head]);
    const commits = (await git(record.worktree, ["rev-list", "--reverse", `${record.base_commit}..${head}`])).split("\n").filter(Boolean);
    const changedFiles = (await git(record.worktree, ["diff", "--name-only", `${record.base_commit}...${head}`])).split("\n").filter(Boolean);
    const outsideScope = changedFiles.filter((path2) => !record.allowed_wrapper_paths.includes(path2));
    if (outsideScope.length > 0) throw new Error(`Context sync changed files outside approved wrapper scope: ${outsideScope.join(", ")}`);
    if (commits.length === 0 || changedFiles.length === 0) throw new Error("Context review requires committed canonical-context changes");
    for (const proposal of request4.proposals.filter((candidate) => candidate.classification === "durable-wrapper")) {
      const content = await readFile19(join15(record.worktree, proposal.target), "utf8");
      if (!content.includes(proposal.source_contribution)) throw new Error(`Canonical update must cite source contribution: ${proposal.source_contribution}`);
    }
    const remotes = (await git(record.worktree, ["remote"])).split("\n").filter(Boolean);
    const remote = remotes.includes("origin") ? "origin" : null;
    const blockers2 = record.wrapper_mode === "pull-request" && !remote ? ["Wrapper has no origin remote; configure one before pushing or opening the required review."] : [];
    record.status = blockers2.length > 0 ? "blocked" : "review-ready";
    record.commits = commits;
    record.changed_files = changedFiles;
    record.remote = remote;
    record.blockers = blockers2;
    record.updated_at = (options.now ?? /* @__PURE__ */ new Date()).toISOString();
    await assertValid7("context-sync-record", record);
    await writeJsonAtomic(recordPath2, record);
    return record;
  });
}
var import_yaml11;
var init_context_sync = __esm({
  "scripts/lib/context-sync.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    import_yaml11 = __toESM(require_dist(), 1);
    init_finish_work();
    init_git();
    init_io();
    init_validation();
  }
});

// scripts/sync-context.ts
var sync_context_exports = {};
import { readFile as readFile20 } from "node:fs/promises";
import { dirname as dirname23, resolve as resolve33 } from "node:path";
import { parseArgs as parseArgs19 } from "node:util";
import { fileURLToPath as fileURLToPath19 } from "node:url";
var workspaceRoot18, values18, request3;
var init_sync_context = __esm({
  async "scripts/sync-context.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_context_sync();
    workspaceRoot18 = resolve33(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve33(dirname23(fileURLToPath19(import.meta.url)), ".."));
    ({ values: values18 } = parseArgs19({ options: { request: { type: "string" } } }));
    if (!values18.request) throw new Error("Usage: sync-context --request <context-sync-request.json>");
    request3 = JSON.parse(await readFile20(resolve33(values18.request), "utf8"));
    console.log(JSON.stringify(await prepareContextSync({ workspaceRoot: workspaceRoot18, request: request3 }), null, 2));
  }
});

// scripts/prepare-context-review.ts
var prepare_context_review_exports = {};
import { dirname as dirname24, resolve as resolve34 } from "node:path";
import { parseArgs as parseArgs20 } from "node:util";
import { fileURLToPath as fileURLToPath20 } from "node:url";
var workspaceRoot19, values19;
var init_prepare_context_review = __esm({
  async "scripts/prepare-context-review.ts"() {
    "use strict";
    init_define_CC_TEMPLATE_INVENTORY();
    init_context_sync();
    workspaceRoot19 = resolve34(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve34(dirname24(fileURLToPath20(import.meta.url)), ".."));
    ({ values: values19 } = parseArgs20({ options: { "sync-id": { type: "string" } } }));
    if (!values19["sync-id"]) throw new Error("Usage: prepare-context-review --sync-id <id>");
    console.log(JSON.stringify(await prepareContextReview({ workspaceRoot: workspaceRoot19, syncId: values19["sync-id"] }), null, 2));
  }
});

// scripts/cc.ts
init_define_CC_TEMPLATE_INVENTORY();
var command2 = process.argv[2];
if (!command2) throw new Error("Usage: cc <command> [arguments]");
process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT = process.cwd();
process.argv.splice(2, 1);
switch (command2) {
  case "validate":
    await init_validate().then(() => validate_exports);
    break;
  case "initialize-workspace":
    await init_initialize_workspace2().then(() => initialize_workspace_exports);
    break;
  case "configure-workspace":
    await init_configure_workspace2().then(() => configure_workspace_exports);
    break;
  case "run-task":
    await init_run_task2().then(() => run_task_exports);
    break;
  case "record-result":
    await init_record_result2().then(() => record_result_exports);
    break;
  case "prepare-repair":
    await init_prepare_repair().then(() => prepare_repair_exports);
    break;
  case "prepare-review":
    await init_prepare_review().then(() => prepare_review_exports);
    break;
  case "record-review-publication":
    await init_record_review_publication().then(() => record_review_publication_exports);
    break;
  case "confirm-merge":
    await init_confirm_merge().then(() => confirm_merge_exports);
    break;
  case "finish-work":
    await init_finish_work2().then(() => finish_work_exports);
    break;
  case "create-plan":
    await init_create_plan().then(() => create_plan_exports);
    break;
  case "validate-plan":
    await init_validate_plan().then(() => validate_plan_exports);
    break;
  case "set-plan-state":
    await init_set_plan_state().then(() => set_plan_state_exports);
    break;
  case "whats-next":
    await init_whats_next2().then(() => whats_next_exports);
    break;
  case "prepare-lifecycle":
    await init_prepare_lifecycle().then(() => prepare_lifecycle_exports);
    break;
  case "record-lifecycle-action":
    await init_record_lifecycle_action().then(() => record_lifecycle_action_exports);
    break;
  case "prepare-plan-publication":
    await init_prepare_plan_publication().then(() => prepare_plan_publication_exports);
    break;
  case "record-plan-publication":
    await init_record_plan_publication().then(() => record_plan_publication_exports);
    break;
  case "sync-context":
    await init_sync_context().then(() => sync_context_exports);
    break;
  case "prepare-context-review":
    await init_prepare_context_review().then(() => prepare_context_review_exports);
    break;
  default:
    throw new Error(`Unknown Context Circuit command: ${command2}`);
}
