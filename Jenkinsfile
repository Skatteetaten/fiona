#!/usr/bin/env groovy

def goleveranse
def config = [
    artifactId               : 'fiona',
    groupId                  : 'no.skatteetaten.aurora',
    deliveryBundleClassifier : "Doozerleveransepakke",
    scriptVersion            : 'v7',
    pipelineScript           : 'https://git.aurora.skead.no/scm/ao/aurora-pipeline-scripts.git',
    iqOrganizationName       : 'Team AOS',
    openShiftBaseImage       : 'turbo',
    openShiftBaseImageVersion: 'latest',
    goVersion                : '1.14',
    artifactPath             : 'bin/',
    sonarQube                : false,
    credentialsId            : "github",
    versionStrategy          : [
        [branch: 'master', versionHint: '1']
    ]
]

fileLoader.withGit(config.pipelineScript, config.scriptVersion) {
  goleveranse = fileLoader.load('templates/goleveranse')
}

goleveranse.run(config.scriptVersion, config)
