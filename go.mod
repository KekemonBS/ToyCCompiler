module kekemon.org/comp

go 1.15

require(
    kekemon.org/comp/parser             v0.0.0
    kekemon.org/comp/lexer              v0.0.0
    kekemon.org/comp/semant             v0.0.0
    kekemon.org/comp/codegen            v0.0.0
)

replace(
    kekemon.org/comp/parser             => ./Parser
    kekemon.org/comp/lexer              => ./Lexer
    kekemon.org/comp/semant             => ./Semantic_analyzer
    kekemon.org/comp/codegen            => ./Code_generator
)
