package pro.qucy.ktype.models

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Dependency(
    @SerialName("id") val id: String
)

@Serializable
data class SampleWithDependency(
    @SerialName("dependency") val dependency: Dependency
)

@Serializable
data class AnAliasedDependency(
    @SerialName("server") val server: Long
)

@Serializable
data class SampleAliasedDependency(
    @SerialName("aliased_dependency") val aliasedDependency: AnAliasedDependency
)

@Serializable
data class SomeWrappedSample(
    @SerialName("title") val title: String
)

