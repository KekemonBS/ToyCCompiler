module kekemon.org/comp/codegen

go 1.15

require(
    kekemon.org/comp/semant             v0.0.0
    kekemon.org/comp/parser             v0.0.0
)

replace(
    kekemon.org/comp/semant             => ./../Semantic_analyzer
    kekemon.org/comp/parser             => ./../Parser
)
